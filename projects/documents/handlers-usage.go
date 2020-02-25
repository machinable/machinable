package documents

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/machinable/machinable/dsi/models"
)

type Usage struct {
	RequestCount      int64         `json:"request_count"`
	TotalResponseTime int64         `json:"-"`
	AvgResponse       int64         `json:"avg_response"`
	StatusCodes       map[int]int64 `json:"status_codes"`
}

// ListCollectionUsage returns the list of activity logs for a project
func (d *Documents) ListCollectionUsage(c *gin.Context) {
	projectID := c.MustGet("projectId").(string)

	// filter anything within x hours
	old := time.Now().Add(-time.Hour * time.Duration(1))
	filter := &models.Filters{
		"created": models.Value{
			models.GTE: old,
		},
		"endpoint_type": models.Value{
			models.EQ: models.EndpointResource,
		},
	}

	// TODO: base this on the api limit for the customer tier
	iLimit := int64(10000)
	logs, err := d.store.ListProjectLogs(projectID, iLimit, 0, filter, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make(map[int64]Usage)

	// transform logs
	for _, log := range logs {
		aligned := log.AlignedCreated

		data, ok := response[aligned]
		if !ok {
			data = Usage{
				StatusCodes: make(map[int]int64),
			}
		}

		data.RequestCount++
		data.TotalResponseTime += log.ResponseTime
		data.StatusCodes[log.StatusCode]++

		response[aligned] = data
	}

	// get average response time
	for key, usage := range response {
		usage.AvgResponse = (usage.TotalResponseTime / usage.RequestCount)
		response[key] = usage
	}

	c.PureJSON(http.StatusOK, gin.H{"items": response})
}

// GetStats returns the size of the collections
func (d *Documents) GetStats(c *gin.Context) {
	projectID := c.MustGet("projectId").(string)

	// retrieve the list of resources
	collections, err := d.store.ListDefinitions(projectID)

	if err != nil {
		c.JSON(err.Code(), gin.H{"error": err.Error()})
		return
	}

	totalStats := &models.Stats{}
	resourceStats := map[string]*models.Stats{}
	for _, col := range collections {
		stats, err := d.store.GetResourceStats(projectID, col.PathName)
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
