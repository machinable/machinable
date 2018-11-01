package middleware

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/nsjostrom/machinable/auth"
	"bitbucket.org/nsjostrom/machinable/management/database"
	"bitbucket.org/nsjostrom/machinable/management/models"
	pdb "bitbucket.org/nsjostrom/machinable/projects/database"
	pmodels "bitbucket.org/nsjostrom/machinable/projects/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

var BEARER = "bearer"
var APIKEY = "apikey"

// ProjectUserAuthzMiddleware authenticates the JWT and verifies the requesting user has access to this project
func ProjectUserAuthzMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get project slug
		// get project from context, inserted into context from subdomain
		project := c.GetString("project")
		if project == "" {
			respondWithError(http.StatusUnauthorized, "invalid project", c)
			return
		}

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
					// c.Set("projects", projects)
					// c.Set("user_id", user["id"])
					// c.Set("username", user["name"])
					c.Set("authType", "user")
					c.Set("authString", user["name"])

					c.Next()
					return
				}
			} else if authType == APIKEY {
				// authenticate api key
				apiKey := vals[1]
				hashedKey := auth.SHA1(apiKey)
				collection := database.Connect().Collection(pdb.KeyDocs(project))

				// Find api key
				documentResult := collection.FindOne(
					nil,
					bson.NewDocument(
						bson.EC.String("key_hash", hashedKey),
					),
					nil,
				)

				if documentResult == nil {
					respondWithError(http.StatusInternalServerError, "invalid key", c)
					return
				}

				// Decode result into document
				key := &pmodels.ProjectAPIKey{}
				err := documentResult.Decode(key)
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
				// c.Set("projects", projects)
				// c.Set("user_id", user["id"])
				// c.Set("username", user["name"])
				c.Set("authType", "api key")
				c.Set("authString", key.Description)

				c.Next()
				return
			}

			respondWithError(http.StatusUnauthorized, "invalid access token", c)
			return
		}
		// if no Authorization header is present, load the project and check the authn policy

		// get the project collection
		col := database.Connect().Collection(database.Projects)

		// look up the user
		documentResult := col.FindOne(
			nil,
			bson.NewDocument(
				bson.EC.String("slug", project),
			),
			nil,
		)

		prj := &models.Project{}
		// decode user document
		err := documentResult.Decode(prj)
		if err != nil {
			respondWithError(http.StatusNotFound, "project not found", c)
			return
		}

		if !prj.Authn {
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
func ProjectLoggingMiddleware() gin.HandlerFunc {
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

		if authString == "" {
			authString = "anonymous"
		} else {
			authType = authType + ":"
		}

		plog := &pmodels.Log{
			Event:      fmt.Sprintf("%s %s", verb, path),
			StatusCode: statusCode,
			Created:    time.Now(),
			Initiator:  fmt.Sprintf("%s %s", authType, authString),
		}

		// Get the logs collection
		rc := pdb.Collection(pdb.LogDocs(projectSlug))
		_, err := rc.InsertOne(
			context.Background(),
			plog,
		)

		if err != nil {
			log.Println("an error occured trying to save the log")
			log.Println(err.Error())
		}
	}
}
