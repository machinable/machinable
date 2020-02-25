package projects

import (
	"github.com/go-redis/redis"
	"github.com/machinable/machinable/config"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/events"

	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/middleware"
	"github.com/machinable/machinable/projects/apikeys"
	"github.com/machinable/machinable/projects/documents"
	"github.com/machinable/machinable/projects/hooks"
	"github.com/machinable/machinable/projects/jsontree"
	"github.com/machinable/machinable/projects/logs"
	"github.com/machinable/machinable/projects/resources"
	"github.com/machinable/machinable/projects/sessions"
	"github.com/machinable/machinable/projects/spec"
	"github.com/machinable/machinable/projects/users"
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
