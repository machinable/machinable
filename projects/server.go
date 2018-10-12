package projects

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/middleware"
	"bitbucket.org/nsjostrom/machinable/projects/handlers"
	"github.com/gin-gonic/gin"
)

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// CreateProjectRoutes creates a gin.Engine for the project routes
func CreateProjectRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.OpenCORSMiddleware())
	router.Use(middleware.SubDomainMiddleware())

	collections := router.Group("/collections")
	collections.GET("/", handlers.GetCollections)
	collections.POST("/", handlers.AddCollection)
	collections.POST("/:collectionName", handlers.AddObjectToCollection)
	collections.GET("/:collectionName", handlers.GetObjectsFromCollection)
	collections.GET("/:collectionName/:objectID", notImplemented)
	collections.DELETE("/:collectionName/:objectID", notImplemented)

	// TODO JSON Tree with any layer accessible via HTTP URL Path
	//collections.GET("/:collectionName/*collectionKeys", notImplemented)

	resources := router.Group("/resources")
	resources.POST("/", handlers.AddResourceDefinition)
	resources.GET("/", handlers.ListResourceDefinitions)
	resources.GET("/:resourceDefinitionID", handlers.GetResourceDefinition)
	resources.DELETE("/:resourceDefinitionID", handlers.DeleteResourceDefinition)

	api := router.Group("/api")
	api.POST("/:resourcePathName", handlers.AddObject)
	api.GET("/:resourcePathName", handlers.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handlers.GetObject)
	api.DELETE("/:resourcePathName/:resourceID", handlers.DeleteObject)

	// TODO
	users := router.Group("/users")
	users.GET("/", notImplemented)          // get list of users for this project
	users.POST("/", notImplemented)         // create a new user of this project
	users.POST("/sessions", notImplemented) // create a new session for a user

	tokens := router.Group("/tokens")
	tokens.GET("/", notImplemented)  // get list of api tokens for this project
	tokens.POST("/", notImplemented) // create a new api token for this project

	return router
}
