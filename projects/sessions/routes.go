package sessions

import (
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets all of the appropriate routes to handlers for project sessions
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// sessions have a mixed authz policy so there is a route here and at /mgmt/sessions
	sessions := engine.Group("/sessions")
	sessions.Use(middleware.ProjectLoggingMiddleware(datastore))
	sessions.POST("/", handler.CreateSession)             // create a new session
	sessions.DELETE("/:sessionID", handler.RevokeSession) // delete this user's session TODO: AUTH

	// App mgmt routes with different authz policy
	mgmt := engine.Group("/mgmt")
	mgmt.Use(middleware.ProjectLoggingMiddleware(datastore))
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware(datastore))

	mgmtSessions := mgmt.Group("/sessions")
	mgmtSessions.GET("/", handler.ListSessions)
	mgmtSessions.DELETE("/:sessionID", handler.RevokeSession)

	return nil
}
