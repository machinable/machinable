package users

import (
	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"bitbucket.org/nsjostrom/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets all of the appropriate routes to handlers for project users
func SetRoutes(engine *gin.Engine, datastore interfaces.ProjectUsersDatastore) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	// Only app users have access to user management
	users := engine.Group("/users")
	users.Use(middleware.AppUserJwtAuthzMiddleware())
	users.Use(middleware.AppUserProjectAuthzMiddleware())

	users.GET("/", handler.ListUsers)            // get list of users for this project
	users.POST("/", handler.AddUser)             // create a new user of this project
	users.GET("/:userID", handler.GetUser)       // get a single user of this project
	users.DELETE("/:userID", handler.DeleteUser) // delete a user of this project

	return nil
}
