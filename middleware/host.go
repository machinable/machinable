package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SubDomainMiddleware controls the cross origin policies.
func SubDomainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		hostParts := strings.Split(c.Request.Host, ".")
		if len(hostParts) < 3 {
			c.AbortWithStatus(404)
			return
		}

		subDomain := hostParts[0]
		if subDomain == "" {
			respondWithError(http.StatusUnauthorized, "invalid project", c)
			return
		}

		c.Set("project", subDomain)

		c.Next()
	}
}
