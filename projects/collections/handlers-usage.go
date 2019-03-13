package collections

import (
	"net/http"

	"github.com/anothrnick/machinable/dsi/models"
	"github.com/gin-gonic/gin"
)

// GetStats returns the size of the collections
func (h *Collections) GetStats(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	// retrieve the list of collections
	collections, err := h.store.GetCollections(projectSlug)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	totalStats := &models.Stats{}
	collectionStats := map[string]*models.Stats{}
	for _, col := range collections {
		stats, err := h.store.GetCollectionStats(projectSlug, col.Name)
		if err != nil {
			c.JSON(err.Code(), gin.H{"error": err.Error()})
			return
		}

		collectionStats[col.Name] = stats
		totalStats.Size += stats.Size
		totalStats.Count += stats.Count
	}

	c.JSON(http.StatusOK, gin.H{"total": totalStats, "collections": collectionStats})
}
