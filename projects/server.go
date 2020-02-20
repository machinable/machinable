package projects

import (
	"github.com/anothrnick/machinable/config"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/events"
	"github.com/go-redis/redis"

	"github.com/anothrnick/machinable/middleware"
	"github.com/anothrnick/machinable/projects/apikeys"
	"github.com/anothrnick/machinable/projects/documents"
	"github.com/anothrnick/machinable/projects/hooks"
	"github.com/anothrnick/machinable/projects/jsontree"
	"github.com/anothrnick/machinable/projects/logs"
	"github.com/anothrnick/machinable/projects/resources"
	"github.com/anothrnick/machinable/projects/sessions"
	"github.com/anothrnick/machinable/projects/spec"
	"github.com/anothrnick/machinable/projects/users"
	"github.com/gin-gonic/gin"
)

// CreateRoutes creates a gin.Engine for the project routes
func CreateRoutes(datastore interfaces.Datastore, cache redis.UniversalClient, processor *events.Processor, config *config.AppConfig) *gin.Engine {

	router := gin.Default()
	router.Use(middleware.OpenCORSMiddleware())
	router.Use(middleware.SubDomainMiddleware())

	// set routes -> handlers for each package
	resources.SetRoutes(router, datastore, config)
	documents.SetRoutes(router, datastore, cache, processor, config)
	logs.SetRoutes(router, datastore, config)
	users.SetRoutes(router, datastore, config)
	sessions.SetRoutes(router, datastore, config)
	apikeys.SetRoutes(router, datastore, config)
	jsontree.SetRoutes(router, datastore, cache, processor, config)
	spec.SetRoutes(router, datastore)
	hooks.SetRoutes(router, datastore, config)

	return router
}
