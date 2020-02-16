package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anothrnick/machinable/auth"
	"github.com/anothrnick/machinable/dsi/interfaces"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// Resources is the constant value for the URL parameter
var Resources = "api"

// JSONKey is the constant value for the URL parameter
var JSONKey = "json"

// BEARER is the key for the bearer authorization token
var BEARER = "bearer"

// APIKEY is the key for the apikey authorization token
var APIKEY = "apikey"

// StoreConfig holds the middleware-relevant config for a collection/resource
type StoreConfig struct {
	Create        bool
	Read          bool
	Update        bool
	Delete        bool
	ParallelRead  bool
	ParallelWrite bool
	Headers       map[string]string
}

// VerbRequiresAuthn returns if the provided HTTP Verb requires authentication for this endpoint
func (s *StoreConfig) VerbRequiresAuthn(verb string) (bool, error) {
	switch verb {
	case "POST":
		return s.Create, nil
	case "GET":
		return s.Read, nil
	case "PUT":
		return s.Update, nil
	case "DELETE":
		return s.Delete, nil
	default:
		return false, errors.New("invalid verb")
	}
}

// ProjectUserRegistrationMiddleware verifies this project has User Registration enabled
func ProjectUserRegistrationMiddleware(store interfaces.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get project slug
		// get project from context, inserted into context from subdomain
		project := c.GetString("project")

		// load the project and check the authn policy
		prj, err := store.GetProjectBySlug(project)
		if err != nil {
			respondWithError(http.StatusNotFound, "project not found", c)
			return
		}

		if !prj.UserRegistration {
			respondWithError(http.StatusNotFound, "path not found", c)
			return
		}

		// inject projectId
		c.Set("projectId", prj.ID)

		c.Next()
	}
}

// ProjectAuthzBuildFiltersMiddleware builds the necessary filters based on the requester's role, permissions, as well as
// the collection/resource's access policies. This middleware requires that the requester `role` has been injected into
// the context.
func ProjectAuthzBuildFiltersMiddleware(store interfaces.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get project from context, inserted into context from subdomain
		verb := c.Request.Method
		filters := map[string]interface{}{}

		// get store config
		storei, exists := c.Get("storeConfig")
		if !exists {
			respondWithError(http.StatusBadRequest, "malformed request - invalid store", c)
			return
		}
		storeConfig := storei.(StoreConfig)

		// check verb authentication policy
		requiresAuthn, err := storeConfig.VerbRequiresAuthn(verb)
		if err != nil {
			respondWithError(http.StatusNotImplemented, "unexpected HTTP verb when checking for authentication", c)
			return
		}
		// if this verb does not require authn, or the user is creating an object (no need for creator filter), let it on by!
		if !requiresAuthn || verb == "POST" {
			c.Set("filters", filters)
			c.Next()
			return
		}

		rRole := c.GetString("authRole")
		rID := c.GetString("authID")

		// based on the requester's role and resource access policies, build filters
		if rRole == auth.RoleUser {
			if verb == "GET" && storeConfig.ParallelRead == false {
				filters["_metadata.creator"] = rID
			} else if (verb == "PUT" || verb == "DELETE") && storeConfig.ParallelWrite == false {
				filters["_metadata.creator"] = rID
			}

			c.Set("filters", filters)
			c.Next()
			return
		} else if rRole == auth.RoleAdmin {
			// `admin` role:
			//    no filter needed
			c.Set("filters", filters)
			c.Next()
			return
		}

		// unknown role, cancel request
		respondWithError(http.StatusForbidden, "unknown role", c)
		return
	}
}

// ProjectUserAuthzMiddleware authenticates the JWT and verifies the requesting user has access to this project. This middleware
// requires that the `project` has been injected into the context.
func ProjectUserAuthzMiddleware(store interfaces.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get request method
		verb := c.Request.Method
		// get project slug
		// get project from context, inserted into context from subdomain
		projectSlug := c.GetString("project")
		if projectSlug == "" {
			respondWithError(http.StatusUnauthorized, "invalid project", c)
			return
		}

		// load the project and check the authn policy
		project, err := store.GetProjectDetailBySlug(projectSlug)
		if err != nil {
			respondWithError(http.StatusNotFound, "project not found", c)
			return
		}

		// TODO: query hooks as part of project detail view
		hooks, err := store.ListHooks(project.ID)
		if err != nil {
			respondWithError(http.StatusNotFound, "error loading project details", c)
			return
		}
		project.Hooks = hooks

		c.Set("projectObject", project)
		c.Set("projectId", project.ID)
		c.Set("accountRequestLimit", project.Requests)
		c.Set("accountId", project.UserID)

		// load resource access policies
		params := strings.Split(c.Request.URL.Path, "/")

		if len(params) < 3 {
			respondWithError(http.StatusBadRequest, "malformed request - invalid params", c)
			return
		}

		storeType := params[1]
		storeConfig := StoreConfig{}

		// check store type, load store and get config for access policies
		if storeType == Resources {
			resourceName := params[2]
			// TODO: Perhaps move this to a view with the project so we only make one DB query?
			def, err := store.GetDefinitionByPathName(project.ID, resourceName)
			if err != nil {
				respondWithError(http.StatusNotFound, "error retrieving resource - does not exist", c)
				return
			}
			storeConfig.Create = def.Create
			storeConfig.Read = def.Read
			storeConfig.Update = def.Update
			storeConfig.Delete = def.Delete
			storeConfig.ParallelRead = def.ParallelRead
			storeConfig.ParallelWrite = def.ParallelWrite
		} else if storeType == JSONKey {
			rootKeyStr := params[2]
			rootKey, err := store.GetRootKey(project.ID, rootKeyStr)
			if err != nil {
				respondWithError(http.StatusNotFound, "error retrieving root key - does not exist", c)
				return
			}
			storeConfig.Create = rootKey.Create
			storeConfig.Read = rootKey.Read
			storeConfig.Update = rootKey.Update
			storeConfig.Delete = rootKey.Delete
		} else {
			respondWithError(http.StatusNotFound, "not found", c)
			return
		}

		// put store in context
		c.Set("storeConfig", storeConfig)

		// validate Authorization header
		if values, _ := c.Request.Header["Authorization"]; len(values) > 0 {

			vals := strings.Split(values[0], " ")

			authType := strings.ToLower(vals[0])

			if authType == BEARER {
				tokenString := vals[1]
				token, err := jwt.Parse(tokenString, auth.TokenLookup)

				if err == nil {

					// token is valid, get claims and perform authorization
					claims := token.Claims.(jwt.MapClaims)

					projects, ok := claims["projects"].(map[string]interface{})
					if !ok {
						respondWithError(http.StatusUnauthorized, "improperly formatted access token", c)
						return
					}

					_, ok = projects[projectSlug]
					if !ok {
						// project user does not have access to this project
						respondWithError(http.StatusNotFound, "project not found", c)
						return
					}

					user, ok := claims["user"].(map[string]interface{})
					if !ok {
						respondWithError(http.StatusUnauthorized, "improperly formatted access token", c)
						return
					}

					userType, ok := user["type"].(string)
					if !ok || userType != "project" {
						respondWithError(http.StatusUnauthorized, "invalid access token", c)
						return
					}

					userIsActive, ok := user["active"].(bool)
					if !ok || !userIsActive {
						respondWithError(http.StatusUnauthorized, "user is not active, please confirm your account", c)
						return
					}

					// check user permissions
					perms := map[string]bool{}
					if user["read"].(bool) {
						perms["GET"] = true
					}
					if user["write"].(bool) {
						perms["POST"] = true
						perms["DELETE"] = true
						perms["PUT"] = true
						// perms["PATCH"] = true
					}

					if _, ok := perms[verb]; !ok {
						respondWithError(http.StatusUnauthorized, fmt.Sprintf("user does not have permission to '%s'", verb), c)
						return
					}

					// inject claims into context
					c.Set("authType", "user")
					c.Set("authString", user["name"])
					c.Set("authID", user["id"])
					c.Set("authRole", user["role"])

					c.Next()
					return
				}
			} else if authType == APIKEY {
				// authenticate api key
				if len(vals) < 2 {
					respondWithError(http.StatusNotFound, "invalid key", c)
					return
				}
				apiKey := vals[1]
				hashedKey := auth.SHA1(apiKey)

				key, err := store.GetAPIKeyByKey(project.ID, hashedKey)
				if err != nil {
					respondWithError(http.StatusNotFound, "invalid key", c)
					return
				}
				// check user permissions
				perms := map[string]bool{}
				if key.Read {
					perms["GET"] = true
				}
				if key.Write {
					perms["POST"] = true
					perms["DELETE"] = true
					perms["PUT"] = true
					// perms["PATCH"] = true
				}

				if _, ok := perms[verb]; !ok {
					respondWithError(http.StatusUnauthorized, fmt.Sprintf("user does not have permission to '%s'", verb), c)
					return
				}

				// inject claims into context
				c.Set("authType", "apikey")
				c.Set("authString", key.Description)
				c.Set("authID", key.ID)
				c.Set("authRole", key.Role)

				c.Next()
				return
			}

			respondWithError(http.StatusUnauthorized, "invalid access token", c)
			return
		}

		// if no Authorization header is present, check the authn policy
		requiresAuthn, err := storeConfig.VerbRequiresAuthn(verb)
		if err != nil {
			respondWithError(http.StatusNotImplemented, "unexpected HTTP verb when checking for authentication", c)
			return
		}
		if !requiresAuthn {
			c.Set("authType", "anonymous")
			c.Set("authString", "anonymous")
			c.Set("authID", "anonymous")
			c.Set("authRole", "anonymous")

			// project does not require authentication, carry on
			c.Next()
			return
		}

		respondWithError(http.StatusUnauthorized, "access token required", c)
		return
	}
}

// RequestRateLimit checks the account rate limit and returns 429 if over app tier limit
func RequestRateLimit(store interfaces.Datastore, cache redis.UniversalClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountLimit := c.GetInt("accountRequestLimit")
		accountID := c.GetString("accountId")
		hour := time.Now().Hour()
		currentKey := fmt.Sprintf("requestCount:%s:%d", accountID, hour)

		// get the request count key for the current window
		val, err := cache.Get(currentKey).Int()

		if err == redis.Nil {
			// {currentKey} does not exist
		} else if err != nil {
			log.Println("could not read from cache ", err.Error())
			// continue handler chain as to not disrupt user experience
			c.Next()
		}

		if val > accountLimit {
			respondWithError(http.StatusTooManyRequests, "request count exceeded account rate limit", c)
			return
		}

		// increment and set request count in redis
		val++
		// expire key after 1 hour
		err = cache.Set(currentKey, val, time.Hour*1).Err()
		if err != nil {
			log.Println("could not write to cache ", err.Error())
			// continue handler chain as to not disrupt user experience
			c.Next()
		}

		// currently under rate limit, continue handler chain
		c.Next()
	}
}
