package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/config"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/middleware"
)

// SetRoutes sets all of the appropriate routes to handlers for project sessions
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, config *config.AppConfig) error {
	// create new Resources handler with datastore
	handler := New(datastore, config)

	// sessions have a mixed authz policy so there is a route here and at /mgmt/sessions
	sessions := engine.Group("/sessions")
	sessions.POST("/", handler.CreateSession)             // create a new session
	sessions.DELETE("/:sessionID", handler.RevokeSession) // delete this user's session TODO: AUTH
	sessions.POST("/refresh", middleware.ValidateRefreshToken(config), handler.RefreshSession)

	// App mgmt routes with different authz policy
	mgmt := engine.Group("/mgmt")
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware(config))
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware(datastore, config))

	mgmtSessions := mgmt.Group("/sessions")
	mgmtSessions.GET("/", handler.ListSessions)
	mgmtSessions.DELETE("/:sessionID", handler.RevokeSession)

	return nil
}
