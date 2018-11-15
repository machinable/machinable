package resources

import (
	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets all of the appropriate routes to handlers for project collections
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// project/user routes
	// api := engine.Group("/api")
	// api.Use(middleware.ProjectLoggingMiddleware())
	// api.Use(middleware.ProjectUserAuthzMiddleware())
	// api.POST("/:resourcePathName", handlers.AddObject)
	// api.GET("/:resourcePathName", handlers.ListObjects)
	// api.GET("/:resourcePathName/:resourceID", handlers.GetObject)
	// api.DELETE("/:resourcePathName/:resourceID", handlers.DeleteObject)

	// admin/mgmt routes
	// Only application users have access to resource definitions
	resources := engine.Group("/resources")
	resources.Use(middleware.AppUserJwtAuthzMiddleware())
	resources.Use(middleware.AppUserProjectAuthzMiddleware())

	resources.POST("/", handler.AddResourceDefinition)
	resources.GET("/", handler.ListResourceDefinitions)
	resources.GET("/:resourceDefinitionID", handler.GetResourceDefinition)
	resources.DELETE("/:resourceDefinitionID", handler.DeleteResourceDefinition)

	// App mgmt routes with different authz policy
	// mgmt := engine.Group("/mgmt")
	// mgmt.Use(middleware.ProjectLoggingMiddleware())
	// mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	// mgmt.Use(middleware.AppUserProjectAuthzMiddleware())

	// mgmtAPI := mgmt.Group("/api")
	// mgmtAPI.POST("/:resourcePathName", handlers.AddObject)
	// mgmtAPI.GET("/:resourcePathName", handlers.ListObjects)
	// mgmtAPI.GET("/:resourcePathName/:resourceID", handlers.GetObject)
	// mgmtAPI.DELETE("/:resourcePathName/:resourceID", handlers.DeleteObject)

	return nil
}
