package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// ProjectLogsDatastore exposes functions to the project access logs
type ProjectLogsDatastore interface {
	AddProjectLog(projectID string, log *models.Log) error
	ListProjectLogs(projectID string, limit, offset int64, filter *models.Filters, sort map[string]int) ([]*models.Log, error)
	CountProjectLogs(projectID string, filter *models.Filters) (int64, error)
	DropProjectLogs(projectID string) error
}
