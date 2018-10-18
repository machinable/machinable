package management

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/management/handlers"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Us? We're fine... how are you?"})
}

func version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "0.0.0"})
}

// CreateManagementRoutes creates a gin.Engine with routes to the application management resources
func CreateManagementRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.OpenCORSMiddleware())

	// meta endpoint for health and version
	meta := router.Group("/meta")
	meta.GET("/health", health)
	meta.GET("/version", version)

	// user endpoints
	users := router.Group("/users")
	users.POST("/register", handlers.RegisterUser)
	users.POST("/sessions", handlers.LoginUser)
	users.DELETE("/sessions/:sessionID", middleware.AppUserJwtAuthzMiddleware(), handlers.RevokeSession)
	users.POST("/refresh", middleware.ValidateRefreshToken(), handlers.RefreshToken)
	users.POST("/password", middleware.AppUserJwtAuthzMiddleware(), handlers.ResetPassword)

	// project endpoints
	projects := router.Group("/projects")
	projects.Use(middleware.AppUserJwtAuthzMiddleware())
	projects.GET("/", handlers.ListUserProjects)
	projects.POST("/", handlers.CreateProject)
	projects.DELETE("/:projectSlug", handlers.DeleteUserProject)

	// user settings endpoints
	// TODO

	return router
}
