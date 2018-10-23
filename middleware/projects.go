package middleware

import (
	"github.com/gin-gonic/gin"
)

// ProjectUserAuthzMiddleware authenticates the JWT and verifies the requesting user has access to this project
func ProjectUserAuthzMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get project slug
		// validate Authorization header
		// if no Authorization header is present, load the project and check the authn policy
	}
}
