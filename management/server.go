package management

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"bitbucket.org/nsjostrom/machinable/management/projects"
	"bitbucket.org/nsjostrom/machinable/management/users"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Us? We're fine... how are you?"})
}

func version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "0.0.0"})
}

// CreateRoutes creates a gin.Engine with routes to the application management resources
func CreateRoutes(datastore interfaces.Datastore) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.OpenCORSMiddleware())

	// meta endpoint for health and version
	meta := router.Group("/meta")
	meta.GET("/health", health)
	meta.GET("/version", version)

	// user endpoints
	users.SetRoutes(router, datastore)
	projects.SetRoutes(router, datastore)

	return router
}
