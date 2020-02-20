package management

import (
	"net/http"

	"github.com/go-redis/redis"

	"github.com/anothrnick/machinable/config"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/management/projects"
	"github.com/anothrnick/machinable/management/users"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Us? We're fine... how are you?"})
}

func version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "0.0.0"})
}

// CreateRoutes creates a gin.Engine with routes to the application management resources
func CreateRoutes(datastore interfaces.Datastore, cache redis.UniversalClient, config *config.AppConfig) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.OpenCORSMiddleware())

	// meta endpoint for health and version
	meta := router.Group("/meta")
	meta.GET("/health", health)
	meta.GET("/version", version)

	// user endpoints
	users.SetRoutes(router, datastore, cache, config)
	projects.SetRoutes(router, datastore, config)

	return router
}
