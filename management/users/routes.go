package users

import (
	"github.com/anothrnick/machinable/config"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// SetRoutes sets all of the appropriate routes to handlers for users
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore, cache redis.UniversalClient, config *config.AppConfig) error {
	handler := New(datastore, cache, config)

	appUserMiddleware := middleware.AppUserJwtAuthzMiddleware(config)
	// user endpoints
	users := engine.Group("/users")
	users.GET("/", appUserMiddleware, handler.GetUser)
	users.GET("/tiers", appUserMiddleware, handler.ListTiers)
	users.GET("/usage", appUserMiddleware, handler.GetUsage)
	users.POST("/register", handler.RegisterUser)
	users.POST("/sessions", handler.LoginUser)
	users.GET("/sessions", appUserMiddleware, handler.ListUserSessions)
	users.GET("/sessions/:sessionID", appUserMiddleware, handler.GetSession)
	users.DELETE("/sessions/:sessionID", appUserMiddleware, handler.RevokeSession)
	users.POST("/refresh", middleware.ValidateRefreshToken(config), handler.RefreshToken)
	users.POST("/password", appUserMiddleware, handler.ResetPassword)

	return nil
}
