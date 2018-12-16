package documents

import (
	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets all of the appropriate routes to handlers for project collections
func SetRoutes(engine *gin.Engine, datastore interfaces.ResourcesDatastore) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// project/user routes
	api := engine.Group("/api")
	api.Use(middleware.ProjectLoggingMiddleware())
	api.Use(middleware.ProjectUserAuthzMiddleware())
	api.POST("/:resourcePathName", handler.AddObject)
	api.GET("/:resourcePathName", handler.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handler.GetObject)
	api.DELETE("/:resourcePathName/:resourceID", handler.DeleteObject)

	// App mgmt routes with different authz policy
	mgmt := engine.Group("/mgmt")
	mgmt.Use(middleware.ProjectLoggingMiddleware())
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware())

	mgmtAPI := mgmt.Group("/api")
	mgmtAPI.POST("/:resourcePathName", handler.AddObject)
	mgmtAPI.GET("/:resourcePathName", handler.ListObjects)
	mgmtAPI.GET("/:resourcePathName/:resourceID", handler.GetObject)
	mgmtAPI.DELETE("/:resourcePathName/:resourceID", handler.DeleteObject)

	return nil
}
