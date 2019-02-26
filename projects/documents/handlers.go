package documents

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/anothrnick/machinable/dsi"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	localModels "github.com/anothrnick/machinable/projects/models"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new `Documents` struct
func New(db interfaces.ResourcesDatastore) *Documents {
	return &Documents{
		store: db,
	}
}

// Documents contains the datastore and any HTTP handlers for project resource documents
type Documents struct {
	store interfaces.ResourcesDatastore
}

// AddObject creates a new document of the resource definition
func (h *Documents) AddObject(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	projectSlug := c.MustGet("project").(string)
	creator := c.MustGet("authID").(string)
	creatorType := c.MustGet("authType").(string)

	fieldValues := models.ResourceObject{}

	err := c.BindJSON(&fieldValues)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meta := models.NewMetaData(creator, creatorType)

	newID, dsiErr := h.store.AddDefDocument(projectSlug, resourcePathName, fieldValues, meta)
	if dsiErr != nil {
		c.JSON(dsiErr.Code(), gin.H{"error": "failed to save " + resourcePathName, "errors": strings.Split(dsiErr.Error(), ",")})
		return
	}

	// Set the inserted ID for the response
	fieldValues["id"] = newID

	c.JSON(http.StatusCreated, fieldValues)
}

// ListObjects returns the list of objects for a resource
func (h *Documents) ListObjects(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	projectSlug := c.MustGet("project").(string)
	authFilters := c.MustGet("filters").(map[string]interface{})

	if resourcePathName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource cannot be empty"})
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

	// Parse and validate pagination
	il, err := strconv.Atoi(limit)
	if err != nil || il > localModels.MaxLimit || il <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	iLimit := int64(il)
	io, err := strconv.Atoi(offset)
	if err != nil || io < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}
	iOffset := int64(io)

	// Format query parameters
	filter := make(map[string]interface{})
	sort := make(map[string]int)

	var resourceDefinition *models.ResourceDefinition
	validProperties := map[string]models.Property{}
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

		if resourceDefinition == nil || validProperties == nil {
			// get resource definition if we do not already have it
			resourceDefinition, err := h.store.GetDefinitionByPathName(projectSlug, resourcePathName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve resource definition to validate query parameters"})
				return
			}

			// get property types
			var pErr error
			validProperties, pErr = resourceDefinition.GetProperties()
			if pErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting property types"})
				return
			}
		}

		prop, ok := validProperties[k]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to filter on '%s', field does not exist", k)})
			return
		}

		// cast filters to their actual types, based on definition
		trueValue, err := dsi.CastInterfaceToType(prop.Type, v[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter[k] = trueValue
	}

	// get accurate count based on auth filters
	docCount, countErr := h.store.CountDefDocuments(projectSlug, resourcePathName, authFilters)

	if countErr != nil {
		c.JSON(countErr.Code(), gin.H{"error": countErr.Error()})
		return
	}

	pageMax := (docCount % iLimit) + docCount
	if (iLimit+iOffset) > pageMax && iOffset >= docCount && docCount != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	// Apply authorization filters
	for k, v := range authFilters {
		filter[k] = v
	}

	documents, dsiErr := h.store.ListDefDocuments(projectSlug, resourcePathName, iLimit, iOffset, filter, sort)

	if dsiErr != nil {
		c.JSON(dsiErr.Code(), gin.H{"error": dsiErr.Error()})
		return
	}

	links := localModels.NewLinks(c.Request, iLimit, iOffset, docCount)

	c.PureJSON(http.StatusOK, gin.H{"items": documents, "links": links, "count": docCount})
}

// GetObject returns a single object with the resourceID for this resource
func (h *Documents) GetObject(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	resourceID := c.Param("resourceID")
	projectSlug := c.MustGet("project").(string)

	document, err := h.store.GetDefDocument(projectSlug, resourcePathName, resourceID)

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
	projectSlug := c.MustGet("project").(string)

	err := h.store.DeleteDefDocument(projectSlug, resourcePathName, resourceID)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
