package hooks

import (
	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/config"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/middleware"
)

// Handler is an interface to the API Key HTTP handler functions.
type Handler interface {
	UpdateHook(c *gin.Context)
	AddHook(c *gin.Context)
	ListHooks(c *gin.Context)
	GetHook(c *gin.Context)
	DeleteHook(c *gin.Context)
	ListResults(c *gin.Context)
}

// SetRoutes sets all of the appropriate routes to handlers for project users
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, config *config.AppConfig) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// di for testing
	return setRoutes(
		engine,
		handler,
		datastore,
		middleware.AppUserJwtAuthzMiddleware(config),
		middleware.AppUserProjectAuthzMiddleware(datastore, config),
	)
}

func setRoutes(engine *gin.Engine, handler Handler, datastore interfaces.ProjectHooksDatastore, mw ...gin.HandlerFunc) error {
	// Only app users have access to api key management
	keys := engine.Group("/hooks")
	keys.Use(mw...)

	keys.GET("/", handler.ListHooks)                  // get list of project web hooks
	keys.POST("/", handler.AddHook)                   // create a new project web hook
	keys.DELETE("/:hookID", handler.DeleteHook)       // delete a project web hook
	keys.PUT("/:hookID", handler.UpdateHook)          // update a project web hook
	keys.GET("/:hookID", handler.GetHook)             // get a project web hook
	keys.GET("/:hookID/results", handler.ListResults) // get list of hook results

	return nil
}
