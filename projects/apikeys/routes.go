package apikeys

import (
	"github.com/anothrnick/machinable/config"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// Handler is an interface to the API Key HTTP handler functions.
type Handler interface {
	UpdateKey(c *gin.Context)
	AddKey(c *gin.Context)
	ListKeys(c *gin.Context)
	GenerateKey(c *gin.Context)
	DeleteKey(c *gin.Context)
}

// SetRoutes sets all of the appropriate routes to handlers for project users
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, config *config.AppConfig) error {
	// create new Resources handler with datastore
	handler := New(datastore, config)

	// di for testing
	return setRoutes(
		engine,
		handler,
		datastore,
		middleware.AppUserJwtAuthzMiddleware(config),
		middleware.AppUserProjectAuthzMiddleware(datastore, config),
	)
}

func setRoutes(engine *gin.Engine, handler Handler, datastore interfaces.ProjectAPIKeysDatastore, mw ...gin.HandlerFunc) error {
	// Only app users have access to api key management
	keys := engine.Group("/keys")
	keys.Use(mw...)

	keys.GET("/generate", handler.GenerateKey) // generate a valid uuid for a potential api key
	keys.GET("/", handler.ListKeys)            // get list of api keys for this project
	keys.POST("/", handler.AddKey)             // create a new api key for this project
	keys.DELETE("/:keyID", handler.DeleteKey)  // delete an api key
	keys.PUT("/:keyID", handler.UpdateKey)     // update an api key

	return nil
}
