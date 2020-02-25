package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/config"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/middleware"
)

// SetRoutes sets all of the appropriate routes to handlers for project collections
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, config *config.AppConfig) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// admin/mgmt routes
	// Only application users have access to resource definitions
	resources := engine.Group("/resources")
	resources.Use(middleware.AppUserJwtAuthzMiddleware(config))
	resources.Use(middleware.AppUserProjectAuthzMiddleware(datastore, config))

	resources.POST("/", handler.AddResourceDefinition)
	resources.GET("/", handler.ListResourceDefinitions)
	resources.GET("/:resourceDefinitionID", handler.GetResourceDefinition)
	resources.PUT("/:resourceDefinitionID", handler.UpdateResourceDefinition)
	resources.DELETE("/:resourceDefinitionID", handler.DeleteResourceDefinition)

	return nil
}
