package documents

import (
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets all of the appropriate routes to handlers for project collections
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// project/user routes
	api := engine.Group("/api")
	api.Use(middleware.ResourceStatsMiddleware(datastore))
	api.Use(middleware.ProjectLoggingMiddleware(datastore))
	api.Use(middleware.ProjectUserAuthzMiddleware(datastore))
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
	mgmtStats.GET("/stats", handler.GetStats)
	mgmtStats.GET("/responseTimes", handler.ListResponseTimes)
	mgmtStats.GET("/statusCodes", handler.ListStatusCodes)

	mgmt.Use(middleware.ProjectLoggingMiddleware(datastore))
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware(datastore))

	mgmtAPI := mgmt.Group("/api")
	// mgmtAPI.POST("/:resourcePathName", handler.AddObject)
	mgmtAPI.GET("/:resourcePathName", handler.ListObjects)
	// mgmtAPI.GET("/:resourcePathName/:resourceID", handler.GetObject)
	// mgmtAPI.PUT("/:resourcePathName/:resourceID", handler.PutObject)
	// mgmtAPI.DELETE("/:resourcePathName/:resourceID", handler.DeleteObject)

	return nil
}
