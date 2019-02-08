package collections

import (
	"net/http"
	"strconv"

	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	localModels "github.com/anothrnick/machinable/projects/models"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new `Collections`
func New(db interfaces.CollectionsDatastore) *Collections {
	return &Collections{
		store: db,
	}
}

// Collections contains the datastore and any HTTP handlers to collections
type Collections struct {
	store interfaces.CollectionsDatastore
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
	err := h.store.AddCollection(projectSlug, &newCollection)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
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
		c.JSON(err.Code(), gin.H{"error": err.Error()})
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
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// PutObjectInCollection alters an existing collection document
func (h *Collections) PutObjectInCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	objectIDStr := c.Param("objectID")
	projectSlug := c.MustGet("project").(string)
	creator := c.MustGet("authID").(string)
	creatorType := c.MustGet("authType").(string)

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

	meta := models.NewMetaData(creator, creatorType)

	// updated the document
	updateErr := h.store.UpdateCollectionDocument(projectSlug, collectionName, objectIDStr, bdoc, meta)

	if err != nil {
		c.JSON(updateErr.Code(), gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, bdoc)
}

// AddObjectToCollection adds a new document to the collection
func (h *Collections) AddObjectToCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	projectSlug := c.MustGet("project").(string)
	creator := c.MustGet("authID").(string)
	creatorType := c.MustGet("authType").(string)

	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	// get or create the project collection
	if _, err := h.store.GetCollection(projectSlug, collectionName); err != nil {
		if err := h.store.AddCollection(projectSlug, &models.Collection{Name: collectionName}); err != nil {
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

	meta := models.NewMetaData(creator, creatorType)

	// add the new document
	newDocument, docErr := h.store.AddCollectionDocument(projectSlug, collectionName, bdoc, meta)
	if err != nil {
		c.JSON(docErr.Code(), gin.H{"error": err.Error()})
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

	// Get pagination parameters
	values := c.Request.URL.Query()
	limit := values.Get("_limit")
	offset := values.Get("_offset")

	// Set defaults if necessary
	if limit == "" {
		limit = localModels.Limit
	}

	if offset == "" {
		offset = "0"
	}

	// Clear reserved query parameters

	// Format query parameters
	filter := make(map[string]interface{})

	for k, v := range values {
		if k == "_limit" || k == "_offset" {
			continue
		}
		filter[k] = v[0]
	}

	// Parse and validate pagination
	il, err := strconv.Atoi(limit)
	if err != nil || il > localModels.MaxLimit {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	iLimit := int64(il)
	io, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}
	iOffset := int64(io)

	// Get the total count of documents for pagination
	docCount, colErr := h.store.CountCollectionDocuments(projectSlug, collectionName)

	if colErr != nil {
		c.JSON(colErr.Code(), gin.H{"error": err.Error()})
		return
	}

	pageMax := (docCount % iLimit) + docCount
	if (iLimit+iOffset) > pageMax && iOffset >= docCount && docCount != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	// Retrieve documents for the page
	documents, colErr := h.store.GetCollectionDocuments(projectSlug, collectionName, iLimit, iOffset, filter)

	if colErr != nil {
		c.JSON(colErr.Code(), gin.H{"error": err.Error()})
		return
	}

	links := localModels.NewLinks(c.Request, iLimit, iOffset, docCount)

	c.PureJSON(http.StatusOK, gin.H{"items": documents, "links": links, "count": docCount})
}

// GetObjectFromCollection returns a single object with the ID for this resource
func (h *Collections) GetObjectFromCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	objectID := c.Param("objectID")
	projectSlug := c.MustGet("project").(string)

	object, err := h.store.GetCollectionDocument(projectSlug, collectionName, objectID)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
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
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
