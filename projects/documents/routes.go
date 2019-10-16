package documents

import (
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// SetRoutes sets all of the appropriate routes to handlers for project collections
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, cache redis.UniversalClient) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// project/user routes
	api := engine.Group("/api")
	api.Use(middleware.ResourceStatsMiddleware(datastore))
	api.Use(middleware.ProjectUserAuthzMiddleware(datastore))
	api.Use(middleware.RequestRateLimit(datastore, cache))
	api.Use(middleware.ProjectAuthzBuildFiltersMiddleware(datastore))

	api.POST("/:resourcePathName", handler.AddObject)
	api.GET("/:resourcePathName", handler.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handler.GetObject)
	api.PUT("/:resourcePathName/:resourceID", handler.PutObject)
	api.DELETE("/:resourcePathName/:resourceID", handler.DeleteObject)

	// App mgmt routes with different authz policy
	mgmt := engine.Group("/mgmt")

	mgmtStats := mgmt.Group("/resourceUsage")
	mgmtStats.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmtStats.Use(middleware.AppUserProjectAuthzMiddleware(datastore))
	mgmtStats.GET("/", handler.ListCollectionUsage)
	mgmtStats.GET("/stats", handler.GetStats)

	mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware(datastore))

	mgmtAPI := mgmt.Group("/api")
	mgmtAPI.GET("/:resourcePathName", handler.ListObjects)

	return nil
}
