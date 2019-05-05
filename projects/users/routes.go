package users

import (
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/middleware"
	"github.com/gin-gonic/gin"
)

// SetRoutes sets all of the appropriate routes to handlers for project users
func SetRoutes(engine *gin.Engine, datastore interfaces.Datastore) error {
	// create new Resources handler with datastore
	handler := New(datastore)

	users := engine.Group("/users")
	users.Use(middleware.ProjectLoggingMiddleware(datastore))
	users.Use(middleware.ProjectUserRegistrationMiddleware(datastore))
	users.POST("/register", handler.AddLimitedUser) // create a new user with the role 'user'

	// Only app users have access to user management
	mgmt := engine.Group("/mgmt/users")
	mgmt.Use(middleware.AppUserJwtAuthzMiddleware())
	mgmt.Use(middleware.AppUserProjectAuthzMiddleware(datastore))

	mgmt.GET("/", handler.ListUsers)            // get list of users for this project
	mgmt.POST("/", handler.AddUser)             // create a new user of this project
	mgmt.GET("/:userID", handler.GetUser)       // get a single user of this project
	mgmt.DELETE("/:userID", handler.DeleteUser) // delete a user of this project
	mgmt.PUT("/:userID", handler.UpdateUser)    // update a user of this project

	return nil
}
