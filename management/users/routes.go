package users

import (
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// SetRoutes sets all of the appropriate routes to handlers for users
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, cache redis.UniversalClient) error {
	handler := New(datastore, cache)

	// user endpoints
	users := engine.Group("/users")
	users.GET("/", middleware.AppUserJwtAuthzMiddleware(), handler.GetUser)
	users.GET("/tiers", middleware.AppUserJwtAuthzMiddleware(), handler.ListTiers)
	users.GET("/usage", middleware.AppUserJwtAuthzMiddleware(), handler.GetUsage)
	users.POST("/register", handler.RegisterUser)
	users.POST("/sessions", handler.LoginUser)
	users.GET("/sessions", middleware.AppUserJwtAuthzMiddleware(), handler.ListUserSessions)
	users.GET("/sessions/:sessionID", middleware.AppUserJwtAuthzMiddleware(), handler.GetSession)
	users.DELETE("/sessions/:sessionID", middleware.AppUserJwtAuthzMiddleware(), handler.RevokeSession)
	users.POST("/refresh", middleware.ValidateRefreshToken(), handler.RefreshToken)
	users.POST("/password", middleware.AppUserJwtAuthzMiddleware(), handler.ResetPassword)

	return nil
}
