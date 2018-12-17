package projects

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/dsi/mongo"
	"bitbucket.org/nsjostrom/machinable/dsi/mongo/database"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"bitbucket.org/nsjostrom/machinable/projects/apikeys"
	"bitbucket.org/nsjostrom/machinable/projects/collections"
	"bitbucket.org/nsjostrom/machinable/projects/documents"
	"bitbucket.org/nsjostrom/machinable/projects/logs"
	"bitbucket.org/nsjostrom/machinable/projects/resources"
	"bitbucket.org/nsjostrom/machinable/projects/sessions"
	"bitbucket.org/nsjostrom/machinable/projects/users"
	"github.com/gin-gonic/gin"
)

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
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
	apikeys.SetRoutes(router, datastore)

	return router
}
