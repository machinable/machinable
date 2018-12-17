package apikeys

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/auth"
	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"bitbucket.org/nsjostrom/machinable/projects/models"
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

// AddKey creates a new api key for this project
func (k *APIKeys) AddKey(c *gin.Context) {
	var newKey models.NewProjectKey
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
		newKey.Write)

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
