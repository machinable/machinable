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
	collections.Use(middleware.ProjectUserAuthzMiddleware())
	collections.GET("/", handlers.GetCollections)
	collections.POST("/", handlers.AddCollection)
	collections.POST("/:collectionName", handlers.AddObjectToCollection)
	collections.GET("/:collectionName", handlers.GetObjectsFromCollection)
	collections.GET("/:collectionName/:objectID", handlers.GetObjectFromCollection)
	collections.PUT("/:collectionName/:objectID", handlers.PutObjectInCollection)
	collections.DELETE("/:collectionName/:objectID", handlers.DeleteObjectFromCollection)

	api := engine.Group("/api")
	api.Use(middleware.ProjectUserAuthzMiddleware())
	api.POST("/:resourcePathName", handlers.AddObject)
	api.GET("/:resourcePathName", handlers.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handlers.GetObject)
	api.DELETE("/:resourcePathName/:resourceID", handlers.DeleteObject)

	// sessions have a mixed authz policy so there is a route here and at /mgmt/sessions
	sessions := engine.Group("/sessions")
	//sessions.Use(middleware.ProjectUserAuthzMiddleware())
	sessions.POST("/", handlers.CreateSession)             // create a new session
	sessions.DELETE("/:sessionID", handlers.RevokeSession) // delete this user's session TODO: AUTH
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
	keys := engine.Group("/keys")
	keys.Use(middleware.AppUserJwtAuthzMiddleware())
	keys.Use(middleware.AppUserProjectAuthzMiddleware())

	keys.GET("/generate", handlers.GenerateKey) // get list of api keys for this project
	keys.GET("/", handlers.ListKeys)            // get list of api keys for this project
	keys.POST("/", handlers.AddKey)             // create a new api key for this project

	// stats := engine.Group("/stats")
	// stats.Use(middleware.AppUserJwtAuthzMiddleware())
	// stats.Use(middleware.AppUserProjectAuthzMiddleware())
	// stats.GET("/collections/:collectionID", handlers.GetCollection)

	// App mgmt routes with different authz policy
	mgmt := engine.Group("/mgmt")
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware())

	collections := mgmt.Group("/collections")
	collections.GET("/", handlers.GetCollections)
	collections.POST("/", handlers.AddCollection)
	collections.POST("/:collectionName", handlers.AddObjectToCollection)
	collections.GET("/:collectionName", handlers.GetObjectsFromCollection)
	collections.DELETE("/:collectionName", handlers.DeleteCollection) // this is actually uses collection ID as the parameter, gin does not allow different wildcard names
	collections.GET("/:collectionName/:objectID", handlers.GetObjectFromCollection)
	collections.DELETE("/:collectionName/:objectID", handlers.DeleteObjectFromCollection)

	api := mgmt.Group("/api")
	api.POST("/:resourcePathName", handlers.AddObject)
	api.GET("/:resourcePathName", handlers.ListObjects)
	api.GET("/:resourcePathName/:resourceID", handlers.GetObject)
	api.DELETE("/:resourcePathName/:resourceID", handlers.DeleteObject)

	sessions := mgmt.Group("/sessions")
	sessions.GET("/", handlers.ListSessions)
	sessions.DELETE("/:sessionID", handlers.RevokeSession)
}

// CreateProjectRoutes creates a gin.Engine for the project routes
func CreateProjectRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.OpenCORSMiddleware())
	router.Use(middleware.SubDomainMiddleware())

	setupMgmtRoutes(router)
	setupProjectUserRoutes(router)

	return router
}
