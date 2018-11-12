package collections

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"bitbucket.org/nsjostrom/machinable/projects/models"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new `Collections`
func New(db interfaces.Datastore) *Collections {
	return &Collections{
		store: db,
	}
}

// Collections contains the datastore and any HTTP handlers to collections
type Collections struct {
	store interfaces.Datastore
}

// AddCollection creates a new collection
func (h *Collections) AddCollection(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	var newCollection models.Collection
	c.BindJSON(&newCollection)

	// validate collection name
	if newCollection.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	// add collection and return error if anything goes wrong
	err := h.store.AddCollection(projectSlug, newCollection.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusCreated, gin.H{})
}

// GetCollections returns the list of collections for a user
func (h *Collections) GetCollections(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	// retrieve the list of collections
	collections, err := h.store.GetCollections(projectSlug)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": collections})
}

// DeleteCollection deletes a project collection along with all of it's data
func (h *Collections) DeleteCollection(c *gin.Context) {
	collectionID := c.Param("collectionName")
	projectSlug := c.MustGet("project").(string)

	err := h.store.DeleteCollection(projectSlug, collectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// PutObjectInCollection alters an existing collection document
func (h *Collections) PutObjectInCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	objectIDStr := c.Param("objectID")
	projectSlug := c.MustGet("project").(string)

	bdoc := make(map[string]interface{})

	// get body
	err := c.BindJSON(&bdoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate collection name path parameter
	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	// updated the document
	err = h.store.UpdateCollectionDocument(projectSlug, collectionName, objectIDStr, bdoc)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, bdoc)
}

// AddObjectToCollection adds a new document to the collection
func (h *Collections) AddObjectToCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	projectSlug := c.MustGet("project").(string)
	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	// get or create the project collection
	if _, err := h.store.GetCollection(projectSlug, collectionName); err != nil {
		if err := h.store.AddCollection(projectSlug, collectionName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get collection"})
			return
		}
	}

	bdoc := make(map[string]interface{})

	err := c.BindJSON(&bdoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// add the new document
	newDocument, err := h.store.AddCollectionDocument(projectSlug, collectionName, bdoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newDocument)
}

// GetObjectsFromCollection returns the full list of documents
func (h *Collections) GetObjectsFromCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	projectSlug := c.MustGet("project").(string)
	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	documents, err := h.store.GetCollectionDocuments(projectSlug, collectionName, 0, 0, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"items": documents})
}

// GetObjectFromCollection returns a single object with the ID for this resource
func (h *Collections) GetObjectFromCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	objectID := c.Param("objectID")
	projectSlug := c.MustGet("project").(string)

	object, err := h.store.GetCollectionDocument(projectSlug, collectionName, objectID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, object)
}

// DeleteObjectFromCollection deletes the object from the collection
func (h *Collections) DeleteObjectFromCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	objectID := c.Param("objectID")
	projectSlug := c.MustGet("project").(string)

	err := h.store.DeleteCollectionDocument(projectSlug, collectionName, objectID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
