package interfaces

import (
	"github.com/machinable/machinable/dsi/errors"
	"github.com/machinable/machinable/dsi/models"
)

// ProjectHooksDatastore defines the functions required to interact with the project web hooks datastore
type ProjectHooksDatastore interface {
	AddHook(projectID string, hook *models.WebHook) *errors.DatastoreError
	ListHooks(projectID string) ([]*models.WebHook, *errors.DatastoreError)
	GetHook(projectID, hookID string) (*models.WebHook, *errors.DatastoreError)
	UpdateHook(projectID, hookID string, hook *models.WebHook) *errors.DatastoreError
	DeleteHook(projectID, hookID string) *errors.DatastoreError

	AddResult(result *models.HookResult) *errors.DatastoreError
	ListResults(projectID, hookID string) ([]*models.HookResult, *errors.DatastoreError)
}
