package postgres

import (
	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
)

const tableProjectWebHooks = "project_users"

// AddHook saves a new WebHook to the datastore
func (d *Database) AddHook(projectID string, hook *models.WebHook) *errors.DatastoreError {
	return nil
}

// ListHooks retrieves all WebHooks for a project
func (d *Database) ListHooks(projectID string) ([]*models.WebHook, *errors.DatastoreError) {
	return nil, nil
}

// GetHook retrieves a single hook by project and hook ID, if it exists
func (d *Database) GetHook(projectID, hookID string) (*models.WebHook, *errors.DatastoreError) {
	return nil, nil
}

// UpdateHook updates all fields of a WebHook by project and hook ID
func (d *Database) UpdateHook(projectID, hookID string, hook *models.WebHook) *errors.DatastoreError {
	return nil
}

// DeleteHook permanently removes a WebHook by project and hook ID
func (d *Database) DeleteHook(projectID, hookID string) *errors.DatastoreError {
	return nil
}
