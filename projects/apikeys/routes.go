package apikeys

import (
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets all of the appropriate routes to handlers for project users
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// Only app users have access to api key management
	keys := engine.Group("/keys")
	keys.Use(middleware.AppUserJwtAuthzMiddleware())
	keys.Use(middleware.AppUserProjectAuthzMiddleware(datastore))

	keys.GET("/generate", handler.GenerateKey) // generate a valid uuid for a potential api key
	keys.GET("/", handler.ListKeys)            // get list of api keys for this project
	keys.POST("/", handler.AddKey)             // create a new api key for this project
	keys.DELETE("/:keyID", handler.DeleteKey)  // delete an api key
	keys.PUT("/:keyID", handler.UpdateKey)     // update an api key

	return nil
}
