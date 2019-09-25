package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// ProjectAPIKeysDatastore exposes functions to manage project api keys
type ProjectAPIKeysDatastore interface {
	GetAPIKeyByKey(projectID, hash string) (*models.ProjectAPIKey, error)
	CreateAPIKey(projectID, hash, description string, read, write bool, role string) (*models.ProjectAPIKey, error)
	UpdateAPIKey(projectID, keyID string, read, write bool, role string) error
	ListAPIKeys(projectID string) ([]*models.ProjectAPIKey, error)
	DeleteAPIKey(projectID, keyID string) error
	DropProjectKeys(projectID string) error
}
