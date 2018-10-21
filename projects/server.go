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

func setupProjectUserRoutes(engine *gin.Engine) {
	collections := engine.Group("/collections")
	collections.GET("/", handlers.GetCollections)
	collections.POST("/", handlers.AddCollection)
	collections.POST("/:collectionName", handlers.AddObjectToCollection)
	collections.GET("/:collectionName", handlers.GetObjectsFromCollection)
	collections.GET("/:collectionName/:objectID", handlers.GetObjectFromCollection)
	collections.DELETE("/:collectionName/:objectID", handlers.DeleteObjectFromCollection)

	api := engine.Group("/api")
	api.POST("/:resourcePathName", handlers.AddObject)
	api.GET("/:resourcePathName", handlers.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handlers.GetObject)
	api.DELETE("/:resourcePathName/:resourceID", handlers.DeleteObject)
}

func setupMgmtRoutes(engine *gin.Engine) {
	// Only application users have access to resource definitions
	resources := engine.Group("/resources")
	resources.Use(middleware.AppUserJwtAuthzMiddleware())
	resources.Use(middleware.AppUserProjectAuthzMiddleware())

	resources.POST("/", handlers.AddResourceDefinition)
	resources.GET("/", handlers.ListResourceDefinitions)
	resources.GET("/:resourceDefinitionID", handlers.GetResourceDefinition)
	resources.DELETE("/:resourceDefinitionID", handlers.DeleteResourceDefinition)

	// Only app users have access to user management
	users := engine.Group("/users")
	users.Use(middleware.AppUserJwtAuthzMiddleware())
	users.Use(middleware.AppUserProjectAuthzMiddleware())

	users.GET("/", handlers.ListUsers)            // get list of users for this project
	users.POST("/", handlers.AddUser)             // create a new user of this project
	users.GET("/:userID", handlers.GetUser)       // get a single user of this project
	users.DELETE("/:userID", handlers.DeleteUser) // delete a user of this project

	// Only app users have access to api key management
	tokens := engine.Group("/tokens")
	tokens.Use(middleware.AppUserJwtAuthzMiddleware())
	tokens.Use(middleware.AppUserProjectAuthzMiddleware())

	tokens.GET("/generate", handlers.GenerateToken) // get list of api tokens for this project
	tokens.GET("/", handlers.ListTokens)            // get list of api tokens for this project
	tokens.POST("/", handlers.AddToken)             // create a new api token for this project

	// App mgmt routes with different authz policy
	mgmt := engine.Group("/mgmt")
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware())

	collections := mgmt.Group("/collections")
	collections.GET("/", handlers.GetCollections)
	collections.POST("/", handlers.AddCollection)
	collections.POST("/:collectionName", handlers.AddObjectToCollection)
	collections.GET("/:collectionName", handlers.GetObjectsFromCollection)
	collections.GET("/:collectionName/:objectID", handlers.GetObjectFromCollection)
	collections.DELETE("/:collectionName/:objectID", handlers.DeleteObjectFromCollection)

	api := mgmt.Group("/api")
	api.POST("/:resourcePathName", handlers.AddObject)
	api.GET("/:resourcePathName", handlers.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handlers.GetObject)
	api.DELETE("/:resourcePathName/:resourceID", handlers.DeleteObject)
}

// CreateProjectRoutes creates a gin.Engine for the project routes
func CreateProjectRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.OpenCORSMiddleware())
	router.Use(middleware.SubDomainMiddleware())

	setupMgmtRoutes(router)
	setupProjectUserRoutes(router)

	// sessions has a mixed authz policy
	sessions := router.Group("/sessions")
	sessions.POST("/", handlers.CreateSession)                                       // create a new session
	sessions.GET("/", middleware.AppUserJwtAuthzMiddleware(), handlers.ListSessions) // list sessions of a project

	return router
}
