package interfaces

import "github.com/anothrnick/machinable/dsi/models"

// ProjectJSONDatastore exposes functions to the project json trees
type ProjectJSONDatastore interface {
	GetRootKey(projectID, rootKey string) (*models.RootKey, error)
	ListRootKeys(projectID string) ([]*models.RootKey, error)
	CreateRootKey(projectID, rootKey string, data []byte) error
	UpdateRootKey(projectID string, rootKey *models.RootKey) error
	DeleteRootKey(projectID, rootKey string) error

	GetJSONKey(projectID, rootKey string, keys ...string) ([]byte, error)
	CreateJSONKey(projectID, rootKey string, data []byte, keys ...string) error
	UpdateJSONKey(projectID, rootKey string, data []byte, keys ...string) error
	DeleteJSONKey(projectID, rootKey string, keys ...string) error
}
