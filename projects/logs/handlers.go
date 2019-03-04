package logs

import (
	"net/http"
	"time"

	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/anothrnick/machinable/query"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new instance of the Logs handler
func New(db interfaces.ProjectLogsDatastore) *Logs {
	return &Logs{
		store: db,
	}
}

// Logs wraps handler access to project logs
type Logs struct {
	store interfaces.ProjectLogsDatastore
}

// ListProjectLogs returns the list of activity logs for a project
func (l *Logs) ListProjectLogs(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

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

	// sort by created, descending
	sort := map[string]int{
		"created": -1,
	}

	// filter anything within x hours
	old := time.Now().Add(-time.Hour * time.Duration(24))
	filter := &models.Filters{
		"created": models.Value{
			models.GTE: old.Unix(),
		},
	}

	// TODO: should have total count and filtered count..
	logCount, err := l.store.CountProjectLogs(projectSlug, filter)

	pageMax := (logCount % iLimit) + logCount
	if (iLimit+iOffset) > pageMax && iOffset >= logCount && logCount != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	logs, err := l.store.ListProjectLogs(projectSlug, iLimit, iOffset, filter, sort)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	links := query.NewLinks(c.Request, iLimit, iOffset, logCount)

	c.PureJSON(http.StatusOK, gin.H{"items": logs, "links": links, "count": logCount})
}
