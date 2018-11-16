package projects

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/dsi/mongo"
	"bitbucket.org/nsjostrom/machinable/dsi/mongo/database"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"bitbucket.org/nsjostrom/machinable/projects/collections"
	"bitbucket.org/nsjostrom/machinable/projects/handlers"
	"bitbucket.org/nsjostrom/machinable/projects/resources"
	"github.com/gin-gonic/gin"
)

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func setupProjectUserRoutes(engine *gin.Engine) {
	api := engine.Group("/api")
	api.Use(middleware.ProjectLoggingMiddleware())
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
	// resources := engine.Group("/resources")
	// resources.Use(middleware.AppUserJwtAuthzMiddleware())
	// resources.Use(middleware.AppUserProjectAuthzMiddleware())

	// resources.POST("/", handlers.AddResourceDefinition)
	// resources.GET("/", handlers.ListResourceDefinitions)
	// resources.GET("/:resourceDefinitionID", handlers.GetResourceDefinition)
	// resources.DELETE("/:resourceDefinitionID", handlers.DeleteResourceDefinition)

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
	keys.DELETE("/:keyID", handlers.DeleteKey)  // get list of api keys for this project

	// stats := engine.Group("/stats")
	// stats.Use(middleware.AppUserJwtAuthzMiddleware())
	// stats.Use(middleware.AppUserProjectAuthzMiddleware())
	// stats.GET("/collections/:collectionID", handlers.GetCollection)

	logs := engine.Group("/logs")
	logs.Use(middleware.AppUserJwtAuthzMiddleware())
	logs.Use(middleware.AppUserProjectAuthzMiddleware())
	logs.GET("/", handlers.GetProjectLogs)

	// App mgmt routes with different authz policy
	mgmt := engine.Group("/mgmt")
	mgmt.Use(middleware.ProjectLoggingMiddleware())
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware())

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
	// use mongoDB connector
	// if another connector is needed, the Datastore interface can be implemented and these 2 lines changes to instantiate the new connector
	// potential connectors: InfluxDB, Postgres JSON, Redis, CouchDB
	mongoDB := database.Connect()
	datastore := mongo.New(mongoDB)

	router := gin.Default()
	router.Use(middleware.OpenCORSMiddleware())
	router.Use(middleware.SubDomainMiddleware())

	// set routes -> handlers for each package
	collections.SetRoutes(router, datastore)
	resources.SetRoutes(router, datastore)

	setupMgmtRoutes(router)
	setupProjectUserRoutes(router)

	return router
}
