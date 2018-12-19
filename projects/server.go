package projects

import (
	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"

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

// CreateRoutes creates a gin.Engine for the project routes
func CreateRoutes(datastore interfaces.Datastore) *gin.Engine {

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
