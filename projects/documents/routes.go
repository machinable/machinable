package documents

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/machinable/machinable/config"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/events"
	"github.com/machinable/machinable/middleware"
)

// SetRoutes sets all of the appropriate routes to handlers for project collections
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, cache redis.UniversalClient, processor *events.Processor, config *config.AppConfig) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// project/user routes
	api := engine.Group("/api")
	api.Use(middleware.ResourceStatsMiddleware(datastore, processor))
	api.Use(middleware.ProjectUserAuthzMiddleware(datastore, config))
	api.Use(middleware.RequestRateLimit(datastore, cache))
	api.Use(middleware.ProjectAuthzBuildFiltersMiddleware(datastore))

	api.POST("/:resourcePathName", handler.AddObject)
	api.GET("/:resourcePathName", handler.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handler.GetObject)
	api.PUT("/:resourcePathName/:resourceID", handler.PutObject)
	api.DELETE("/:resourcePathName/:resourceID", handler.DeleteObject)

	// App mgmt routes with different authz policy
	mgmt := engine.Group("/mgmt")
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware(config))
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware(datastore, config))

	// mgmt resource usage
	mgmtStats := mgmt.Group("/resourceUsage")
	mgmtStats.GET("/", handler.ListCollectionUsage)
	mgmtStats.GET("/stats", handler.GetStats)

	// mgmt get objects
	mgmtAPI := mgmt.Group("/api")
	mgmtAPI.GET("/:resourcePathName", handler.ListObjects)

	return nil
}
