package documents

import (
	"net/http"
	"time"

	"github.com/anothrnick/machinable/dsi/models"
	"github.com/gin-gonic/gin"
)

// ListResponseTimes returns HTTP response times for collections for the last 1 hour
func (d *Documents) ListResponseTimes(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	old := time.Now().Add(-time.Hour * time.Duration(1))
	filter := &models.Filters{
		"timestamp": models.Value{
			models.GTE: old.Unix(),
		},
	}

	responseTimes, err := d.store.ListResourceResponseTimes(projectSlug, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response_times": responseTimes})
}

// ListStatusCodes returns HTTP response status codes for collections for the last 1 hour
func (d *Documents) ListStatusCodes(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	old := time.Now().Add(-time.Hour * time.Duration(1))
	filter := &models.Filters{
		"timestamp": models.Value{
			models.GTE: old.Unix(),
		},
	}

	statusCodes, err := d.store.ListResourceStatusCode(projectSlug, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status_codes": statusCodes})
}

// GetStats returns the size of the collections
func (d *Documents) GetStats(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	// retrieve the list of resources
	collections, err := d.store.ListDefinitions(projectSlug)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	totalStats := &models.Stats{}
	resourceStats := map[string]*models.Stats{}
	for _, col := range collections {
		stats, err := d.store.GetResourceStats(projectSlug, col.PathName)
		if err != nil {
			c.JSON(err.Code(), gin.H{"error": err.Error()})
			return
		}

		resourceStats[col.PathName] = stats
		totalStats.Size += stats.Size
		totalStats.Count += stats.Count
	}

	c.JSON(http.StatusOK, gin.H{"total": totalStats, "resources": resourceStats})
}
