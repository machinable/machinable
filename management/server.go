package management

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Us? We're fine... how are you?"})
}

func version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "0.0.0"})
}

// CreateManagementRoutes creates a gin.Engine with routes to the application management resources
func CreateManagementRoutes() *gin.Engine {
	router := gin.Default()
	meta := router.Group("/meta")
	meta.GET("/health", health)
	meta.GET("/version", version)

	return router
}
