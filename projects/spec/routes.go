package spec

import (
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// Handler is an interface to the API Key HTTP handler functions.
type Handler interface {
	GetSpec(c *gin.Context)
}

// SetRoutes sets all of the appropriate routes to handlers for project users
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// di for testing
	return setRoutes(
		engine,
		handler,
		datastore,
		middleware.ProjectIDAuthzMiddleware(datastore),
	)
}

func setRoutes(engine *gin.Engine, handler Handler, datastore interfaces.ResourcesDatastore, mw ...gin.HandlerFunc) error {
	// Only app users have access to api key management
	keys := engine.Group("/spec")
	keys.Use(mw...)

	keys.GET("/", handler.GetSpec) // get openapi spec

	return nil
}
