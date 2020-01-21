package interfaces

import (
	"errors"

	"github.com/anothrnick/machinable/dsi/models"
)

// ProjectAPIKeysDatastore exposes functions to manage project api keys
type ProjectAPIKeysDatastore interface {
	GetAPIKeyByKey(projectID, hash string) (*models.ProjectAPIKey, error)
	CreateAPIKey(projectID, hash, description string, read, write bool, role string) (*models.ProjectAPIKey, error)
	UpdateAPIKey(projectID, keyID string, read, write bool, role string) error
	ListAPIKeys(projectID string) ([]*models.ProjectAPIKey, error)
	DeleteAPIKey(projectID, keyID string) error
	DropProjectKeys(projectID string) error
}

// MockProjectAPIKeysDatastore mocks the datastore functions for ProjectAPIKeysDatastore testing
type MockProjectAPIKeysDatastore struct {
	GetAPIKeyByKeyFunc  func(projectID, hash string) (*models.ProjectAPIKey, error)
	CreateAPIKeyFunc    func(projectID, hash, description string, read, write bool, role string) (*models.ProjectAPIKey, error)
	UpdateAPIKeyFunc    func(projectID, keyID string, read, write bool, role string) error
	ListAPIKeysFunc     func(projectID string) ([]*models.ProjectAPIKey, error)
	DeleteAPIKeyFunc    func(projectID, keyID string) error
	DropProjectKeysFunc func(projectID string) error
}

// GetAPIKeyByKey mock function, calls field if not nil
func (m *MockProjectAPIKeysDatastore) GetAPIKeyByKey(projectID, hash string) (*models.ProjectAPIKey, error) {
	if m.GetAPIKeyByKeyFunc != nil {
		return m.GetAPIKeyByKeyFunc(projectID, hash)
	}
	return nil, errors.New("not implemented")
}

// CreateAPIKey mock function, calls field if not nil
func (m *MockProjectAPIKeysDatastore) CreateAPIKey(projectID, hash, description string, read, write bool, role string) (*models.ProjectAPIKey, error) {
	if m.CreateAPIKeyFunc != nil {
		return m.CreateAPIKeyFunc(projectID, hash, description, read, write, role)
	}
	return nil, errors.New("not implemented")
}

// UpdateAPIKey mock function, calls field if not nil
func (m *MockProjectAPIKeysDatastore) UpdateAPIKey(projectID, keyID string, read, write bool, role string) error {
	if m.UpdateAPIKeyFunc != nil {
		return m.UpdateAPIKeyFunc(projectID, keyID, read, write, role)
	}
	return errors.New("not implemented")
}

// ListAPIKeys mock function, calls field if not nil
func (m *MockProjectAPIKeysDatastore) ListAPIKeys(projectID string) ([]*models.ProjectAPIKey, error) {
	if m.ListAPIKeysFunc != nil {
		return m.ListAPIKeysFunc(projectID)
	}
	return nil, errors.New("not implemented")
}

// DeleteAPIKey mock function, calls field if not nil
func (m *MockProjectAPIKeysDatastore) DeleteAPIKey(projectID, keyID string) error {
	if m.DeleteAPIKeyFunc != nil {
		return m.DeleteAPIKeyFunc(projectID, keyID)
	}
	return errors.New("not implemented")
}

// DropProjectKeys mock function, calls field if not nil
func (m *MockProjectAPIKeysDatastore) DropProjectKeys(projectID string) error {
	if m.DropProjectKeysFunc != nil {
		return m.DropProjectKeysFunc(projectID)
	}
	return errors.New("not implemented")
}
