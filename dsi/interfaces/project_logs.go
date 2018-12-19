package interfaces

import "bitbucket.org/nsjostrom/machinable/dsi/models"

// ProjectLogsDatastore exposes functions to the project access logs
type ProjectLogsDatastore interface {
	AddProjectLog(project string, log *models.Log) error
	GetProjectLogsForLastHours(project string, hours int) ([]*models.Log, error)
}
