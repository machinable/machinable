package management

import (
	"net/http"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/config"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/management/projects"
	"github.com/machinable/machinable/management/users"
	"github.com/machinable/machinable/middleware"
)

// Meta contains meta handlers
type Meta struct {
	config *config.AppConfig
}

// Health the health of the app
func (m *Meta) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Us? We're fine... how are you?"})
}

// Version current version of the api
func (m *Meta) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": m.config.Version})
}

// CreateRoutes creates a gin.Engine with routes to the application management resources
func CreateRoutes(datastore interfaces.Datastore, cache redis.UniversalClient, config *config.AppConfig) *gin.Engine {
	metaHandler := &Meta{config: config}
	router := gin.Default()

	router.Use(middleware.OpenCORSMiddleware())

	// meta endpoint for health and version
	meta := router.Group("/meta")
	meta.GET("/health", metaHandler.Health)
	meta.GET("/version", metaHandler.Version)

	// user endpoints
	users.SetRoutes(router, datastore, cache, config)
	projects.SetRoutes(router, datastore, config)

	return router
}
