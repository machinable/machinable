package interfaces

import (
	"github.com/anothrnick/machinable/dsi/models"
)

// ProjectCollectionUsageDatastore exposes functions to manage project collection usage
type ProjectCollectionUsageDatastore interface {
	SaveCollectionResponseTimes(project string, timestamp int64, responseTimes *models.ResponseTimes) error
	ListCollectionResponseTimes(project string, filter *models.Filters) ([]*models.ResponseTimes, error)

	SaveCollectionStatusCode(project string, timestamp int64, statusCode *models.StatusCode) error
	ListCollectionStatusCode(project string, filter *models.Filters) ([]*models.StatusCode, error)
}

// ProjectResourceUsageDatastore exposes functions to manage project resource usage
type ProjectResourceUsageDatastore interface {
	SaveResourceResponseTimes(project string, timestamp int64, responseTimes *models.ResponseTimes) error
	ListResourceResponseTimes(project string, filter *models.Filters) ([]*models.ResponseTimes, error)

	SaveResourceStatusCode(project string, timestamp int64, statusCode *models.StatusCode) error
	ListResourceStatusCode(project string, filter *models.Filters) ([]*models.StatusCode, error)
}
