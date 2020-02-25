package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/config"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/middleware"
)

// SetRoutes sets all of the appropriate routes to handlers for project collections
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, config *config.AppConfig) error {
	// create new Logs handler with datastore
	handler := New(datastore)

	logs := engine.Group("/logs")
	logs.Use(middleware.AppUserJwtAuthzMiddleware(config))
	logs.Use(middleware.AppUserProjectAuthzMiddleware(datastore, config))
	logs.GET("/", handler.ListProjectLogs)

	return nil
}
