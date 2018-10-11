package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// SubDomainMiddleware controls the cross origin policies.
func SubDomainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.Host)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
