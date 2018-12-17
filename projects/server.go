package projects

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/dsi/mongo"
	"bitbucket.org/nsjostrom/machinable/dsi/mongo/database"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"bitbucket.org/nsjostrom/machinable/projects/collections"
	"bitbucket.org/nsjostrom/machinable/projects/documents"
	"bitbucket.org/nsjostrom/machinable/projects/handlers"
	"bitbucket.org/nsjostrom/machinable/projects/logs"
	"bitbucket.org/nsjostrom/machinable/projects/resources"
	"bitbucket.org/nsjostrom/machinable/projects/sessions"
	"bitbucket.org/nsjostrom/machinable/projects/users"
	"github.com/gin-gonic/gin"
)

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func setupMgmtRoutes(engine *gin.Engine) {
	// Only app users have access to api key management
	keys := engine.Group("/keys")
	keys.Use(middleware.AppUserJwtAuthzMiddleware())
	keys.Use(middleware.AppUserProjectAuthzMiddleware())

	keys.GET("/generate", handlers.GenerateKey) // get list of api keys for this project
	keys.GET("/", handlers.ListKeys)            // get list of api keys for this project
	keys.POST("/", handlers.AddKey)             // create a new api key for this project
	keys.DELETE("/:keyID", handlers.DeleteKey)  // get list of api keys for this project
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
	documents.SetRoutes(router, datastore)
	logs.SetRoutes(router, datastore)
	users.SetRoutes(router, datastore)
	sessions.SetRoutes(router, datastore)

	setupMgmtRoutes(router)

	return router
}
