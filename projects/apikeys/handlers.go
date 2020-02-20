package apikeys

import (
	"net/http"

	"github.com/anothrnick/machinable/auth"
	"github.com/anothrnick/machinable/config"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// New returns a pointer to a new `APIKeys` struct
func New(db interfaces.ProjectAPIKeysDatastore, config *config.AppConfig) *APIKeys {
	return &APIKeys{
		store:  db,
		config: config,
	}
}

// APIKeys wraps the datastore and any HTTP handlers for project api keys
type APIKeys struct {
	store  interfaces.ProjectAPIKeysDatastore
	config *config.AppConfig
}

// UpdateKey updates api key role and access
func (k *APIKeys) UpdateKey(c *gin.Context) {
	var newKey NewProjectKey
	keyID := c.Param("keyID")
	projectID := c.MustGet("projectId").(string)

	c.BindJSON(&newKey)

	err := newKey.ValidateRoleAccess()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = k.store.UpdateAPIKey(
		projectID,
		keyID,
		newKey.Read,
		newKey.Write,
		newKey.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// AddKey creates a new api key for this project
func (k *APIKeys) AddKey(c *gin.Context) {
	var newKey NewProjectKey
	projectID := c.MustGet("projectId").(string)

	c.BindJSON(&newKey)

	err := newKey.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate sha1 key
	keyHash := auth.SHA1(newKey.Key, k.config.AppSecret)
	newKey.Key = ""

	key, err := k.store.CreateAPIKey(
		projectID,
		keyHash,
		newKey.Description,
		newKey.Read,
		newKey.Write,
		newKey.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, key)
}

// ListKeys lists all api tokens of this project
func (k *APIKeys) ListKeys(c *gin.Context) {
	projectID := c.MustGet("projectId").(string)

	keys, err := k.store.ListAPIKeys(projectID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": keys})
}

// GenerateKey retrieves a single api token of this project by ID
func (k *APIKeys) GenerateKey(c *gin.Context) {
	UUID := uuid.NewV4()
	c.JSON(http.StatusOK, gin.H{"key": UUID.String()})
}

// DeleteKey removes an api token by ID
func (k *APIKeys) DeleteKey(c *gin.Context) {
	keyID := c.Param("keyID")
	projectID := c.MustGet("projectId").(string)

	err := k.store.DeleteAPIKey(projectID, keyID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
