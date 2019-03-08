package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// ProjectLogsDatastore exposes functions to the project access logs
type ProjectLogsDatastore interface {
	AddProjectLog(project string, log *models.Log) error
	ListProjectLogs(project string, limit, offset int64, filter *models.Filters, sort map[string]int) ([]*models.Log, error)
	CountProjectLogs(project string, filter *models.Filters) (int64, error)
	DropProjectLogs(project string) error
}
