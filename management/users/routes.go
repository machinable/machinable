package users

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/machinable/machinable/config"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/middleware"
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
	users.POST("/verify/:verificationCode", handler.EmailVerification)
	users.GET("/sessions", appUserMiddleware, handler.ListUserSessions)
	users.GET("/sessions/:sessionID", appUserMiddleware, handler.GetSession)
	users.DELETE("/sessions/:sessionID", appUserMiddleware, handler.RevokeSession)
	users.POST("/refresh", middleware.ValidateRefreshToken(datastore, config), handler.RefreshToken)
	users.POST("/password", appUserMiddleware, handler.ResetPassword)

	return nil
}
