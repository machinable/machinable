package jsontree

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/gin-gonic/gin"
)

// Handlers contains all handler functions
type Handlers struct {
	db interfaces.Datastore
}

// NewHandlers creates and returns a new instance of `Handlers` with the datastore
func NewHandlers(datastore interfaces.Datastore) *Handlers {
	return &Handlers{
		db: datastore,
	}
}

// ListRootKeys returns the full list of root keys for a project, does not include data for the key
func (h *Handlers) ListRootKeys(c *gin.Context) {
	projectID := c.MustGet("projectId").(string)

	rootKeys, err := h.db.ListRootKeys(projectID)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": rootKeys})
}

// CreateRootKey creates a new JSON tree for a rootKey name
func (h *Handlers) CreateRootKey(c *gin.Context) {
	rootKey := c.Param("rootKey")
	projectID := c.MustGet("projectId").(string)
	b, _ := c.GetRawData()

	err := h.db.CreateRootKey(projectID, rootKey, b)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// ReadRootKey retrieves the root JSON tree
func (h *Handlers) ReadRootKey(c *gin.Context) {
	rootKey := c.Param("rootKey")
	projectID := c.MustGet("projectId").(string)

	byt, err := h.db.GetJSONKey(projectID, rootKey)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	var obj interface{}
	json.Unmarshal(byt, &obj)
	c.IndentedJSON(http.StatusOK, obj)
}

// DeleteRootKey deletes the entire rootKey
func (h *Handlers) DeleteRootKey(c *gin.Context) {
	rootKey := c.Param("rootKey")
	projectID := c.MustGet("projectId").(string)

	err := h.db.DeleteRootKey(projectID, rootKey)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// ReadJSONKey retrieves the data stored at the key path provided by the HTTP path parameters
func (h *Handlers) ReadJSONKey(c *gin.Context) {
	rootKey := c.Param("rootKey")
	projectID := c.MustGet("projectId").(string)
	keys := c.Param("keys")

	keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
	byt, err := h.db.GetJSONKey(projectID, rootKey, strings.Split(keys, "/")...)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	var obj interface{}
	json.Unmarshal(byt, &obj)
	c.IndentedJSON(http.StatusOK, obj)
}

// CreateJSONKey updates a key at the key path. An error is returned if the key already exists.
func (h *Handlers) CreateJSONKey(c *gin.Context) {
	rootKey := c.Param("rootKey")
	projectID := c.MustGet("projectId").(string)
	keys := c.Param("keys")

	b, _ := c.GetRawData()

	keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
	if keys == "" {
		c.JSON(http.StatusBadRequest, "no keys provided")
		return
	}

	err := h.db.CreateJSONKey(projectID, rootKey, b, strings.Split(keys, "/")...)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// UpdateJSONKey updates a root key at the key path. The key is created if it does not already exist.
func (h *Handlers) UpdateJSONKey(c *gin.Context) {
	rootKey := c.Param("rootKey")
	projectID := c.MustGet("projectId").(string)
	keys := c.Param("keys")

	b, _ := c.GetRawData()

	keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
	if keys == "" {
		c.JSON(http.StatusBadRequest, "no keys provided")
		return
	}

	err := h.db.UpdateJSONKey(projectID, rootKey, b, strings.Split(keys, "/")...)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// DeleteJSONKey deletes a project key at the key path
func (h *Handlers) DeleteJSONKey(c *gin.Context) {
	rootKey := c.Param("rootKey")
	projectID := c.MustGet("projectId").(string)
	keys := c.Param("keys")
	keys = strings.TrimRight(strings.TrimLeft(keys, "/"), "/")
	if keys == "" {
		c.JSON(http.StatusBadRequest, "no keys provided")
		return
	}

	err := h.db.DeleteJSONKey(projectID, rootKey, strings.Split(keys, "/")...)
	if err != nil {
		tErr := h.db.TranslateError(err)
		c.JSON(tErr.Code, tErr.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
