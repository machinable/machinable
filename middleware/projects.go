package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/nsjostrom/machinable/auth"
	"bitbucket.org/nsjostrom/machinable/management/database"
	"bitbucket.org/nsjostrom/machinable/management/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

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

			tokenString := strings.Split(values[0], " ")[1]
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
				c.Set("projects", projects)
				c.Set("user_id", user["id"])
				c.Set("username", user["name"])

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
