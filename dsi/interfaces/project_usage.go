package interfaces

import (
	"github.com/anothrnick/machinable/dsi/models"
)

// ProjectCollectionUsageDatastore exposes functions to manage project user sessions
type ProjectCollectionUsageDatastore interface {
	SaveResponseTimes(project string, timestamp int64, responseTimes *models.ResponseTimes) error
	ListResponseTimes(project string, filter *models.Filters) ([]*models.ResponseTimes, error)

	SaveStatusCode(project string, timestamp int64, statusCode *models.StatusCode) error
	ListStatusCode(project string, filter *models.Filters) ([]*models.StatusCode, error)
}
