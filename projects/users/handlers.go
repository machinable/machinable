package users

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/auth"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/dsi/models"
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

// UpdateUser updates the role and access of a project user
func (u *Users) UpdateUser(c *gin.Context) {
	var newUser NewProjectUser
	userID := c.Param("userID")
	projectID := c.MustGet("projectId").(string)

	c.BindJSON(&newUser)

	// validate user access and role
	err := newUser.ValidateAccessRole()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update project user object
	user := &models.ProjectUser{
		Read:  newUser.Read,
		Write: newUser.Write,
		Role:  newUser.Role,
	}

	u.store.UpdateUser(projectID, userID, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// AddLimitedUser creates a new user for this project from the unauthenticated route
func (u *Users) AddLimitedUser(c *gin.Context) {
	var newUser NewProjectUser
	projectID := c.MustGet("projectId").(string)

	c.BindJSON(&newUser)
	// override role
	newUser.Role = "user"

	// validate user
	err := newUser.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := u.store.GetUserByUsername(projectID, newUser.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
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
		Created:      time.Now(),
		PasswordHash: passwordHash,
		Username:     newUser.Username,
		Read:         newUser.Read,
		Write:        newUser.Write,
		Role:         newUser.Role,
	}

	u.store.CreateUser(projectID, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// AddUser creates a new user for this project
func (u *Users) AddUser(c *gin.Context) {
	var newUser NewProjectUser
	projectID := c.MustGet("projectId").(string)

	c.BindJSON(&newUser)

	// validate user
	err := newUser.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := u.store.GetUserByUsername(projectID, newUser.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
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
		Created:      time.Now(),
		PasswordHash: passwordHash,
		Username:     newUser.Username,
		Read:         newUser.Read,
		Write:        newUser.Write,
		Role:         newUser.Role,
	}

	u.store.CreateUser(projectID, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// ListUsers lists all users of this project
func (u *Users) ListUsers(c *gin.Context) {
	projectID := c.MustGet("projectId").(string)

	users, err := u.store.ListUsers(projectID)

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
	projectID := c.MustGet("projectId").(string)

	err := u.store.DeleteUser(projectID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
