package users

import (
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets all of the appropriate routes to handlers for users
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore) error {
	handler := New(datastore)

	// user endpoints
	users := engine.Group("/users")
	users.POST("/register", handler.RegisterUser)
	users.POST("/sessions", handler.LoginUser)
	users.GET("/sessions", middleware.AppUserJwtAuthzMiddleware(), handler.ListUserSessions)
	users.DELETE("/sessions/:sessionID", middleware.AppUserJwtAuthzMiddleware(), handler.RevokeSession)
	users.POST("/refresh", middleware.ValidateRefreshToken(), handler.RefreshToken)
	users.POST("/password", middleware.AppUserJwtAuthzMiddleware(), handler.ResetPassword)

	return nil
}
