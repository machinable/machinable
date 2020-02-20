package documents

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anothrnick/machinable/dsi"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/anothrnick/machinable/query"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new `Documents` struct
func New(db interfaces.Datastore) *Documents {
	return &Documents{
		store: db,
	}
}

// Documents contains the datastore and any HTTP handlers for project resource documents
type Documents struct {
	store interfaces.Datastore
}

// AddObject creates a new document of the resource definition
func (h *Documents) AddObject(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	projectID := c.MustGet("projectId").(string)
	creator := c.MustGet("authID").(string)
	creatorType := c.MustGet("authType").(string)

	fieldValues := models.ResourceObject{}

	err := c.BindJSON(&fieldValues)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meta := models.NewMetaData(creator, creatorType)

	newID, dsiErr := h.store.AddDefDocument(projectID, resourcePathName, fieldValues, meta)
	if dsiErr != nil {
		c.JSON(dsiErr.Code(), gin.H{"error": "failed to save " + resourcePathName, "errors": strings.Split(dsiErr.Error(), ",")})
		return
	}

	// Set the inserted ID for the response
	fieldValues["id"] = newID
	fieldValues["_metadata"] = meta

	c.JSON(http.StatusCreated, fieldValues)
}

// PutObject updates an existing document of the resource definition
func (h *Documents) PutObject(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	resourceID := c.Param("resourceID")
	projectID := c.MustGet("projectId").(string)
	authFilters := c.MustGet("filters").(map[string]interface{})

	fieldValues := models.ResourceObject{}

	err := c.BindJSON(&fieldValues)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	object, dsiErr := h.store.UpdateDefDocument(projectID, resourcePathName, resourceID, fieldValues, authFilters)
	if dsiErr != nil {
		c.JSON(dsiErr.Code(), gin.H{"error": "failed to save " + resourcePathName, "errors": strings.Split(dsiErr.Error(), ",")})
		return
	}

	c.JSON(http.StatusOK, object)
}

// ListObjects returns the list of objects for a resource
func (h *Documents) ListObjects(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	projectID := c.MustGet("projectId").(string)
	authFilters := c.MustGet("filters").(map[string]interface{})

	if resourcePathName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource cannot be empty"})
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

	var resourceDefinition *models.ResourceDefinition
	validSchema := &models.JSONSchemaObject{}
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

		if resourceDefinition == nil || validSchema == nil {
			// get resource definition if we do not already have it
			resourceDefinition, err := h.store.GetDefinitionByPathName(projectID, resourcePathName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve resource definition to validate query parameters"})
				return
			}

			// get property types
			var pErr error
			validSchema, pErr = resourceDefinition.GetSchema()
			if pErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting schema property types"})
				return
			}
		}

		_, ok := validSchema.Properties[k]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to filter on '%s'", k)})
			return
		}

		// no need to cast type, let the DSI layer do that (if needed)
		filter[k] = v[0]
	}

	// Apply authorization filters
	for k, v := range authFilters {
		filter[k] = v
	}

	// get accurate count based on auth filters and query filters
	docCount, countErr := h.store.CountDefDocuments(projectID, resourcePathName, filter)

	if countErr != nil {
		c.JSON(countErr.Code(), gin.H{"error": countErr.Error()})
		return
	}

	pageMax := (docCount % iLimit) + docCount
	if (iLimit+iOffset) > pageMax && iOffset >= docCount && docCount != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	documents, dsiErr := h.store.ListDefDocuments(projectID, resourcePathName, iLimit, iOffset, filter, sort)

	if dsiErr != nil {
		c.JSON(dsiErr.Code(), gin.H{"error": dsiErr.Error()})
		return
	}

	links := query.NewLinks(c.Request, iLimit, iOffset, docCount)

	c.PureJSON(http.StatusOK, gin.H{"items": documents, "links": links, "count": docCount})
}

// GetObject returns a single object with the resourceID for this resource
func (h *Documents) GetObject(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	resourceID := c.Param("resourceID")
	projectID := c.MustGet("projectId").(string)
	authFilters := c.MustGet("filters").(map[string]interface{})

	document, err := h.store.GetDefDocument(projectID, resourcePathName, resourceID, authFilters)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, document)
}

// DeleteObject deletes the object from the collection
func (h *Documents) DeleteObject(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	resourceID := c.Param("resourceID")
	projectID := c.MustGet("projectId").(string)
	authFilters := c.MustGet("filters").(map[string]interface{})

	err := h.store.DeleteDefDocument(projectID, resourcePathName, resourceID, authFilters)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
