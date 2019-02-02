package resources

import (
	"net/http"

	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new `Resources` struct
func New(db interfaces.ResourcesDatastore) *Resources {
	return &Resources{
		store: db,
	}
}

// Resources contains the datastore and any HTTP handlers for project resource definitions and documents
type Resources struct {
	store interfaces.ResourcesDatastore
}

// AddResourceDefinition creates a new resource definition in the users' collection
func (h *Resources) AddResourceDefinition(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	// Marshal JSON into ResourceDefinition
	var resourceDefinition models.ResourceDefinition
	c.BindJSON(&resourceDefinition)

	// Validate the definition
	if err := resourceDefinition.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.store.AddDefinition(projectSlug, &resourceDefinition)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	// Set the inserted ID for the response
	resourceDefinition.ID = id
	c.JSON(http.StatusCreated, resourceDefinition)
}

// ListResourceDefinitions returns the list of all resource definitions
func (h *Resources) ListResourceDefinitions(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	definitions, err := h.store.ListDefinitions(projectSlug)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": definitions})
}

// GetResourceDefinition returns a single resource definition
func (h *Resources) GetResourceDefinition(c *gin.Context) {
	resourceID := c.Param("resourceDefinitionID")
	projectSlug := c.MustGet("project").(string)

	def, err := h.store.GetDefinition(projectSlug, resourceID)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, def)
}

// DeleteResourceDefinition deletes the definition and drops the resource collection
func (h *Resources) DeleteResourceDefinition(c *gin.Context) {
	resourceID := c.Param("resourceDefinitionID")
	projectSlug := c.MustGet("project").(string)

	err := h.store.DeleteDefinition(projectSlug, resourceID)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
