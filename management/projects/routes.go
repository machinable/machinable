package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/config"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/middleware"
)

// SetRoutes sets all of the appropriate routes to handlers for projects
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, config *config.AppConfig) error {
	handler := New(datastore)

	// project endpoints
	projects := engine.Group("/projects")
	projects.Use(middleware.AppUserJwtAuthzMiddleware(config))
	projects.GET("/", handler.ListUserProjects)
	projects.POST("/", handler.CreateProject)
	projects.PUT("/:projectSlug", handler.UpdateProject)
	projects.DELETE("/:projectSlug", handler.DeleteUserProject)

	return nil
}
