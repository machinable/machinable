package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/dsi/interfaces"
	"github.com/machinable/machinable/dsi/models"
)

// New returns a pointer to a new `APIKeys` struct
func New(db interfaces.ProjectHooksDatastore) *WebHooks {
	return &WebHooks{
		store: db,
	}
}

// WebHooks wraps the datastore and any HTTP handlers for project web hooks
type WebHooks struct {
	store interfaces.ProjectHooksDatastore
}

// UpdateHook updates an existing project webhook by id and and project id
func (w *WebHooks) UpdateHook(c *gin.Context) {
	hookID := c.Param("hookID")
	projectID := c.MustGet("projectId").(string)

	hook := models.WebHook{}
	c.BindJSON(&hook)
	hook.ProjectID = projectID

	// validate hook before storing
	if err := hook.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update hook in database for project
	err := w.store.UpdateHook(projectID, hookID, &hook)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// AddHook creates a new webhook for a project
func (w *WebHooks) AddHook(c *gin.Context) {
	projectID := c.MustGet("projectId").(string)

	hook := models.WebHook{}
	c.BindJSON(&hook)
	hook.ProjectID = projectID

	// validate hook before storing
	if err := hook.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// save hook to database for project
	err := w.store.AddHook(projectID, &hook)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"hook": &hook})
}

// ListHooks lists all webhooks for a project
func (w *WebHooks) ListHooks(c *gin.Context) {
	projectID := c.MustGet("projectId").(string)

	// retrieve all hooks for the project
	hooks, err := w.store.ListHooks(projectID)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": &hooks})
}

// ListResults returns the full list of hook results
func (w *WebHooks) ListResults(c *gin.Context) {
	hookID := c.Param("hookID")
	projectID := c.MustGet("projectId").(string)

	results, err := w.store.ListResults(projectID, hookID)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": results})
}

// GetHook gets a single webhook for a project
func (w *WebHooks) GetHook(c *gin.Context) {
	hookID := c.Param("hookID")
	projectID := c.MustGet("projectId").(string)

	hook, err := w.store.GetHook(projectID, hookID)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &hook)
}

// DeleteHook deletes a webhook by id and project
func (w *WebHooks) DeleteHook(c *gin.Context) {
	hookID := c.Param("hookID")
	projectID := c.MustGet("projectId").(string)

	err := w.store.DeleteHook(projectID, hookID)
	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
