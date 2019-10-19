package jsontree

import (
	"github.com/anothrNick/machinable/dsi/interfaces"
	"github.com/gin-gonic/gin"
)

// Handler is an interface to the JSON key/val HTTP handler functions.
type Handler interface {
	ListRootKeys(c *gin.Context)
	CreateRootKey(c *gin.Context)
	ReadRootKey(c *gin.Context)
	DeleteRootKey(c *gin.Context)
	ReadJSONKey(c *gin.Context)
	CreateJSONKey(c *gin.Context)
	UpdateJSONKey(c *gin.Context)
	DeleteJSONKey(c *gin.Context)
}

// SetRoutes sets all of the appropriate routes to handlers for the application
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore) error {
	handler := NewHandlers(datastore)

	return setRoutes(engine, handler)
}

func setRoutes(engine *gin.Engine, h Handler) error {
	jsonKeys := engine.Group("/json")

	// TODO middleware

	jsonKeys.GET("/", h.ListRootKeys)             // returns entire root tree
	jsonKeys.GET("/:rootKey", h.ReadRootKey)      // returns entire root tree
	jsonKeys.POST("/:rootKey", h.CreateRootKey)   // create a new tree at `rootKey`
	jsonKeys.DELETE("/:rootKey", h.DeleteRootKey) // root tree must be empty to delete

	jsonKeys.GET("/:rootKey/*keys", h.ReadJSONKey)
	jsonKeys.POST("/:rootKey/*keys", h.CreateJSONKey)
	jsonKeys.PUT("/:rootKey/*keys", h.UpdateJSONKey)
	jsonKeys.DELETE("/:rootKey/*keys", h.DeleteJSONKey)

	return nil
}
