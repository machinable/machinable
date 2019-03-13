package apikeys

import (
	"net/http"

	"github.com/anothrnick/machinable/auth"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// New returns a pointer to a new `APIKeys` struct
func New(db interfaces.ProjectAPIKeysDatastore) *APIKeys {
	return &APIKeys{
		store: db,
	}
}

// APIKeys wraps the datastore and any HTTP handlers for project api keys
type APIKeys struct {
	store interfaces.ProjectAPIKeysDatastore
}

// UpdateKey updates api key role and access
func (k *APIKeys) UpdateKey(c *gin.Context) {
	var newKey NewProjectKey
	keyID := c.Param("keyID")
	projectSlug := c.MustGet("project").(string)

	c.BindJSON(&newKey)

	err := newKey.ValidateRoleAccess()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = k.store.UpdateAPIKey(
		projectSlug,
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
	projectSlug := c.MustGet("project").(string)

	c.BindJSON(&newKey)

	err := newKey.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate sha1 key
	keyHash := auth.SHA1(newKey.Key)
	newKey.Key = ""

	key, err := k.store.CreateAPIKey(
		projectSlug,
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
	projectSlug := c.MustGet("project").(string)

	keys, err := k.store.ListAPIKeys(projectSlug)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": keys})
}

// GenerateKey retrieves a single api token of this project by ID
func (k *APIKeys) GenerateKey(c *gin.Context) {
	UUID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": UUID.String()})
}

// DeleteKey removes an api token by ID
func (k *APIKeys) DeleteKey(c *gin.Context) {
	keyID := c.Param("keyID")
	projectSlug := c.MustGet("project").(string)

	err := k.store.DeleteAPIKey(projectSlug, keyID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
