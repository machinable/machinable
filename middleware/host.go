package middleware

import (
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
		if len(hostParts) != 3 {
			c.AbortWithStatus(404)
			return
		}

		subDomain := hostParts[0]
		c.Set("project", subDomain)

		c.Next()
	}
}
