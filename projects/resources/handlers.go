package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/dsi/models"
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
	// projectSlug := c.MustGet("project").(string)
	projectID := c.MustGet("projectId").(string)
	// Marshal JSON into ResourceDefinition
	var resourceDefinition models.ResourceDefinition
	c.BindJSON(&resourceDefinition)

	// set authentication required for now
	resourceDefinition.Create = true
	resourceDefinition.Read = true
	resourceDefinition.Update = true
	resourceDefinition.Delete = true

	// validate the definition
	if err := resourceDefinition.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify uniqueness
	if _, dne := h.store.GetDefinitionByPathName(projectID, resourceDefinition.PathName); dne == nil {
		// path name already exists for project and path name
		c.JSON(http.StatusBadRequest, gin.H{"error": "path name already in use for project"})
		return
	}

	id, err := h.store.AddDefinition(projectID, &resourceDefinition)

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
	projectID := c.MustGet("projectId").(string)

	definitions, err := h.store.ListDefinitions(projectID)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": definitions})
}

// GetResourceDefinition returns a single resource definition
func (h *Resources) GetResourceDefinition(c *gin.Context) {
	resourceID := c.Param("resourceDefinitionID")
	projectID := c.MustGet("projectId").(string)

	def, err := h.store.GetDefinition(projectID, resourceID)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, def)
}

// UpdateResourceDefinition updates the parallel_read and parallel_write operations of the definition
func (h *Resources) UpdateResourceDefinition(c *gin.Context) {
	projectID := c.MustGet("projectId").(string)
	resourceDefinitionID := c.Param("resourceDefinitionID") // actually uses ID
	var updatedDefinition models.ResourceDefinition
	c.BindJSON(&updatedDefinition)

	// add collection and return error if anything goes wrong
	err := h.store.UpdateDefinition(projectID, resourceDefinitionID, &updatedDefinition)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteResourceDefinition deletes the definition and drops the resource collection
func (h *Resources) DeleteResourceDefinition(c *gin.Context) {
	resourceID := c.Param("resourceDefinitionID")
	projectID := c.MustGet("projectId").(string)

	err := h.store.DeleteDefinition(projectID, resourceID)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
