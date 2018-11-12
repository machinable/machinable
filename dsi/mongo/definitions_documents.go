package mongo

import (
	"bitbucket.org/nsjostrom/machinable/dsi/errors"
	"bitbucket.org/nsjostrom/machinable/dsi/models"
)

// Project resource definitions
func (d *Datastore) AddDefinition(project string, def *models.ResourceDefinition) (string, *errors.DatastoreError) {
	return "", nil
}

func (d *Datastore) ListDefinitions(project string) ([]*models.ResourceDefinition, *errors.DatastoreError) {
	return nil, nil
}

func (d *Datastore) GetDefinition(project, definitionID string) (*models.ResourceDefinition, *errors.DatastoreError) {
	return nil, nil
}

func (d *Datastore) DeleteDefinition(project, definitionID string) *errors.DatastoreError {
	return nil
}
