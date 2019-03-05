package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anothrnick/machinable/auth"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Collection is the constant value for the URL parameter
var Collections = "collections"

// Resources is the constant value for the URL parameter
var Resources = "api"

// BEARER is the key for the bearer authorization token
var BEARER = "bearer"

// APIKEY is the key for the apikey authorization token
var APIKEY = "apikey"

// ProjectAuthzBuildFiltersMiddleware builds the necessary filters based on the requester's role, permissions, as well as
// the collection/resource's access policies. This middleware requires that the requester `role` has been injected into
// the context.
func ProjectAuthzBuildFiltersMiddleware(store interfaces.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get project from context, inserted into context from subdomain
		project := c.GetString("project")
		_, projectAuthn := c.Get("projectAuthn")
		verb := c.Request.Method
		filters := map[string]interface{}{}
		// if the projectAuthn key exists, this project does not require authn or authz
		// if the requester is doing a POST, just continue
		if !projectAuthn || verb == "POST" {
			c.Set("filters", filters)
			c.Next()
			return
		}

		rRole := c.GetString("authRole")
		rID := c.GetString("authID")

		// based on the requester's role and collection/resource access policies, build filters
		if rRole == auth.RoleUser {
			// `user` role:
			//	  > Load collection/resource access policies
			params := strings.Split(c.Request.URL.Path, "/")

			if len(params) < 3 {
				respondWithError(http.StatusBadRequest, "malformed request - invalid params", c)
				return
			}

			storeType := params[1]
			collectionName := params[2]
			parallelRead, parallelWrite := false, false

			if storeType == Collections {
				col, err := store.GetCollection(project, collectionName)
				if err != nil {
					respondWithError(http.StatusInternalServerError, "error retrieving collection", c)
					return
				}
				parallelRead = col.ParallelRead
				parallelWrite = col.ParallelWrite
			} else if storeType == Resources {
				def, err := store.GetDefinitionByPathName(project, collectionName)
				if err != nil {
					respondWithError(http.StatusInternalServerError, "error retrieving resource", c)
					return
				}
				parallelRead = def.ParallelRead
				parallelWrite = def.ParallelWrite
				fmt.Println("resource filters not supported")
			} else {
				respondWithError(http.StatusBadRequest, "malformed request - unknown path", c)
				return
			}

			if verb == "GET" && parallelRead == false {
				filters["_metadata.creator"] = rID
			} else if (verb == "PUT" || verb == "DELETE") && parallelWrite == false {
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
		// get project slug
		// get project from context, inserted into context from subdomain
		project := c.GetString("project")

		// load the project and check the authn policy
		prj, err := store.GetProjectBySlug(project)
		if err != nil {
			respondWithError(http.StatusNotFound, "project not found", c)
			return
		}

		c.Set("projectAuthn", prj.Authn)

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

					_, ok = projects[project]
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

					if _, ok := perms[c.Request.Method]; !ok {
						respondWithError(http.StatusUnauthorized, fmt.Sprintf("user does not have permission to '%s'", c.Request.Method), c)
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

				key, err := store.GetAPIKeyByKey(project, hashedKey)
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

				if _, ok := perms[c.Request.Method]; !ok {
					respondWithError(http.StatusUnauthorized, fmt.Sprintf("user does not have permission to '%s'", c.Request.Method), c)
					return
				}

				// inject claims into context
				c.Set("authType", "apikey")
				c.Set("authString", key.Description)
				c.Set("authID", key.ID.Hex())
				c.Set("authRole", key.Role)

				c.Next()
				return
			}

			respondWithError(http.StatusUnauthorized, "invalid access token", c)
			return
		}

		// if no Authorization header is present, load the project and check the authn policy
		if !prj.Authn {
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
		// inject custom writer
		lw := &logWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = lw

		// continue handler chain
		c.Next()

		// get status code
		statusCode := c.Writer.Status()
		verb := c.Request.Method
		path := c.Request.URL.Path

		projectSlug := c.GetString("project")
		authType := c.GetString("authType")
		authString := c.GetString("authString")
		authID := c.GetString("authID")

		if authString == "" {
			authString = "unknown"
		}

		plog := &models.Log{
			Event:         fmt.Sprintf("%s %s", verb, path),
			StatusCode:    statusCode,
			Created:       time.Now().Unix(),
			Initiator:     authString,
			InitiatorType: authType,
			InitiatorID:   authID,
		}

		// Get the logs collection
		err := store.AddProjectLog(projectSlug, plog)

		if err != nil {
			log.Println("an error occured trying to save the log")
			log.Println(err.Error())
		}
	}
}
