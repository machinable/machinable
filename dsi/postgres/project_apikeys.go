package postgres

import "github.com/anothrnick/machinable/dsi/models"

const tableProjectAPIKeys = "project_apikeys"

// GetAPIKeyByKey retrieves a single api key by key hash
func (d *Database) GetAPIKeyByKey(project, hash string) (*models.ProjectAPIKey, error) {
	return nil, nil
}

// CreateAPIKey creates a new api key for the project
func (d *Database) CreateAPIKey(project, hash, description string, read, write bool, role string) (*models.ProjectAPIKey, error) {
	return nil, nil
}

// UpdateAPIKey updates the role and access of an API key
func (d *Database) UpdateAPIKey(project, keyID string, read, write bool, role string) error {
	return nil
}

// ListAPIKeys retrieves all api keys for a project
func (d *Database) ListAPIKeys(project string) ([]*models.ProjectAPIKey, error) {
	return nil, nil
}

// DeleteAPIKey removes a project api key permanently
func (d *Database) DeleteAPIKey(project, keyID string) error {
	return nil
}

// DropProjectKeys drops the key collection for this project
func (d *Database) DropProjectKeys(project string) error {
	return nil
}
