package users

import (
	"net/http"
	"time"

	"github.com/anothrnick/machinable/auth"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// New returns a pointer to a new `Users` struct
func New(db interfaces.ProjectUsersDatastore) *Users {
	return &Users{
		store: db,
	}
}

// Users wraps the datastore and any HTTP handlers for project users
type Users struct {
	store interfaces.ProjectUsersDatastore
}

// AddUser creates a new user for this project
func (u *Users) AddUser(c *gin.Context) {
	var newUser NewProjectUser
	projectSlug := c.MustGet("project").(string)

	c.BindJSON(&newUser)

	// validate user
	err := newUser.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate hashed password
	passwordHash, err := auth.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newUser.Password = ""

	// create project user object
	user := &models.ProjectUser{
		ID:           objectid.New(), // I don't like this
		Created:      time.Now(),
		PasswordHash: passwordHash,
		Username:     newUser.Username,
		Read:         newUser.Read,
		Write:        newUser.Write,
		Role:         newUser.Role,
	}

	u.store.CreateUser(projectSlug, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// ListUsers lists all users of this project
func (u *Users) ListUsers(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	users, err := u.store.ListUsers(projectSlug)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": users})
}

// GetUser retrieves a single user of this project by ID
func (u *Users) GetUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// DeleteUser removes a user by ID
func (u *Users) DeleteUser(c *gin.Context) {
	userID := c.Param("userID")
	projectSlug := c.MustGet("project").(string)

	err := u.store.DeleteUser(projectSlug, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
