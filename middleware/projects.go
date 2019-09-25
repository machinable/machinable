package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/anothrnick/machinable/auth"
	"github.com/anothrnick/machinable/dsi/interfaces"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Resources is the constant value for the URL parameter
var Resources = "api"

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
		store := storei.(StoreConfig)

		// check verb authentication policy
		requiresAuthn, err := store.VerbRequiresAuthn(verb)
		if err != nil {
			respondWithError(http.StatusNotImplemented, "unexpected HTTP verb when checking for authentication", c)
			return
		}
		if !requiresAuthn || verb == "POST" {
			c.Set("filters", filters)
			c.Next()
			return
		}

		rRole := c.GetString("authRole")
		rID := c.GetString("authID")

		// based on the requester's role and collection/resource access policies, build filters
		if rRole == auth.RoleUser {
			if verb == "GET" && store.ParallelRead == false {
				filters["_metadata.creator"] = rID
			} else if (verb == "PUT" || verb == "DELETE") && store.ParallelWrite == false {
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
		// else if rRole == auth.RoleAnon {
		// 	// `anon` role:
		// 	//    no filter needed, trust that the previous middleware checked the project policy

		// 	c.Next()
		// 	return
		// }

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
		project, err := store.GetProjectBySlug(projectSlug)
		if err != nil {
			respondWithError(http.StatusNotFound, "project not found", c)
			return
		}

		c.Set("projectId", project.ID)

		//	  > Load collection/resource access policies
		params := strings.Split(c.Request.URL.Path, "/")

		if len(params) < 3 {
			respondWithError(http.StatusBadRequest, "malformed request - invalid params", c)
			return
		}

		storeType := params[1]
		collectionName := params[2]

		storeConfig := StoreConfig{}

		// check store type, load store and get config for access policies
		if storeType == Resources {
			def, err := store.GetDefinitionByPathName(project.ID, collectionName)
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

type logWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w logWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ProjectLoggingMiddleware logs the request
func ProjectLoggingMiddleware(store interfaces.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// continue handler chain
		c.Next()
	}
}
