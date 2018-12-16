package logs

import (
	"net/http"

	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"github.com/gin-gonic/gin"
)

func New(db interfaces.ProjectLogsDatastore) *Logs {
	return &Logs{
		store: db,
	}
}

type Logs struct {
	store interfaces.ProjectLogsDatastore
}

// GetProjectLogs returns the list of activity logs for a project
func (l *Logs) GetProjectLogs(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	logs, err := l.store.GetProjectLogsForLastHours(projectSlug, 24)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"items": logs})
}
