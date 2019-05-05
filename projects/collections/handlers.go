package collections

import (
	"net/http"

	"github.com/anothrnick/machinable/dsi"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/anothrnick/machinable/query"
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

	// set authentication required for now
	newCollection.Create = true
	newCollection.Read = true
	newCollection.Update = true
	newCollection.Delete = true

	// validate collection name
	if newCollection.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	if err := newCollection.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// UpdateCollection updates the parallel_read and parallel_write operations of the collection
func (h *Collections) UpdateCollection(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	collectionID := c.Param("collectionName") // actually uses ID
	var updatedCollection models.Collection
	c.BindJSON(&updatedCollection)

	// add collection and return error if anything goes wrong
	err := h.store.UpdateCollection(projectSlug, collectionID, &updatedCollection)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, gin.H{})
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
	authFilters := c.MustGet("filters").(map[string]interface{})

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
	updateErr := h.store.UpdateCollectionDocument(projectSlug, collectionName, objectIDStr, bdoc, meta, authFilters)

	if updateErr != nil {
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
	requestedCollection := &models.Collection{Name: collectionName}

	bdoc := make(map[string]interface{})

	err := c.BindJSON(&bdoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dsi.ContainsReservedField(bdoc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	if err := requestedCollection.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get or create the project collection
	if _, err := h.store.GetCollection(projectSlug, collectionName); err != nil {
		// set authentication to default
		requestedCollection.Create = true
		requestedCollection.Read = true
		requestedCollection.Update = true
		requestedCollection.Delete = true
		if err := h.store.AddCollection(projectSlug, requestedCollection); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get collection"})
			return
		}
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
	authFilters := c.MustGet("filters").(map[string]interface{})

	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	// Get pagination parameters
	values := c.Request.URL.Query()

	iLimit, err := query.GetLimit(&values)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	iOffset, err := query.GetOffset(&values)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format query parameters
	filter := make(map[string]interface{})
	sort := make(map[string]int)

	for k, v := range values {
		if k == dsi.LimitKey || k == dsi.OffsetKey {
			continue
		}

		if k == dsi.SortKey {
			sortField := v[0]
			firstChar := string(sortField[0])
			order := 1
			if firstChar == "-" {
				order = -1
				sortField = sortField[1:]
			}
			sort[sortField] = order
			continue
		}
		filter[k] = v[0]
	}

	// Apply authorization filters
	for k, v := range authFilters {
		filter[k] = v
	}

	// Get the total count of documents for pagination. Use injected auth filters and query filters to get accurate count relevant to the requester
	docCount, colErr := h.store.CountCollectionDocuments(projectSlug, collectionName, filter)
	if colErr != nil {
		c.JSON(colErr.Code(), gin.H{"error": colErr.Error()})
		return
	}

	pageMax := (docCount % iLimit) + docCount
	if (iLimit+iOffset) > pageMax && iOffset >= docCount && docCount != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	// Retrieve documents for the page
	documents, colErr := h.store.GetCollectionDocuments(projectSlug, collectionName, iLimit, iOffset, filter, sort)

	if colErr != nil {
		c.JSON(colErr.Code(), gin.H{"error": err.Error()})
		return
	}

	links := query.NewLinks(c.Request, iLimit, iOffset, docCount)

	c.PureJSON(http.StatusOK, gin.H{"items": documents, "links": links, "count": docCount})
}

// GetObjectFromCollection returns a single object with the ID for this resource
func (h *Collections) GetObjectFromCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	objectID := c.Param("objectID")
	projectSlug := c.MustGet("project").(string)
	authFilters := c.MustGet("filters").(map[string]interface{})

	object, err := h.store.GetCollectionDocument(projectSlug, collectionName, objectID, authFilters)

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
	authFilters := c.MustGet("filters").(map[string]interface{})

	err := h.store.DeleteCollectionDocument(projectSlug, collectionName, objectID, authFilters)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
