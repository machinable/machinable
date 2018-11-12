package mongo

import "bitbucket.org/nsjostrom/machinable/dsi/models"

// Project resource definitions
func (d *Datastore) AddDefinition(project string, def *models.ResourceDefinition) (string, error) {
	return "", nil
}

func (d *Datastore) ListDefinitions(project string) ([]*models.ResourceDefinition, error) {
	return nil, nil
}

func (d *Datastore) GetDefinition(project, definitionID string) (*models.ResourceDefinition, error) {
	return nil, nil
}

func (d *Datastore) DeleteDefinition(project, definitionID string) error {
	return nil
}
