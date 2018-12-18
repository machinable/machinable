package collections

import (
	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets all of the appropriate routes to handlers for project collections
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore) error {
	// create new Collections handler with datastore, set routes -> handlers
	handler := New(datastore)

	// routes for http api access to collections
	collections := engine.Group("/collections")

	collections.Use(middleware.ProjectLoggingMiddleware(datastore))
	collections.Use(middleware.ProjectUserAuthzMiddleware(datastore))

	collections.GET("/", handler.GetCollections)
	collections.POST("/", handler.AddCollection)
	collections.POST("/:collectionName", handler.AddObjectToCollection)
	collections.GET("/:collectionName", handler.GetObjectsFromCollection)
	collections.GET("/:collectionName/:objectID", handler.GetObjectFromCollection)
	collections.PUT("/:collectionName/:objectID", handler.PutObjectInCollection)
	collections.DELETE("/:collectionName/:objectID", handler.DeleteObjectFromCollection)

	// routes for admin http api access to collections

	// admin routes with different authz policy
	mgmt := engine.Group("/mgmt")
	mgmt.Use(middleware.ProjectLoggingMiddleware(datastore))
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware())

	mgmtCollections := mgmt.Group("/collections")
	mgmtCollections.GET("/", handler.GetCollections)
	mgmtCollections.POST("/", handler.AddCollection)
	mgmtCollections.POST("/:collectionName", handler.AddObjectToCollection)
	mgmtCollections.GET("/:collectionName", handler.GetObjectsFromCollection)
	mgmtCollections.DELETE("/:collectionName", handler.DeleteCollection) // this actually uses collection ID as the parameter, gin does not allow different wildcard names
	mgmtCollections.GET("/:collectionName/:objectID", handler.GetObjectFromCollection)
	mgmtCollections.DELETE("/:collectionName/:objectID", handler.DeleteObjectFromCollection)

	return nil
}
