package interfaces

import "bitbucket.org/nsjostrom/machinable/dsi/models"

// Datastore exposes the necessary functions to interact with the Machinable datastore.
// Functions are grouped logically based on their purpose and the collections they interact with.
type Datastore interface {
	// Project definition documents
	AddDefDocument(project, path string, fields map[string]interface{}) (string, error)
	ListDefDocuments(project, path string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, error)
	GetDefDocument(project, path, documentID string) (map[string]interface{}, error)
	DeleteDefDocument(project, path, documentID string) error
	DropAllDefDocuments(project, path string) error

	// Project resource definitions
	AddDefinition(project string, def *models.ResourceDefinition) (string, error)
	ListDefinitions(project string) ([]*models.ResourceDefinition, error)
	GetDefinition(project, definitionID string) (*models.ResourceDefinition, error)
	DeleteDefinition(project, definitionID string) error

	// Project collections
	AddCollection(project, name string) error
	GetCollection(project, name string) (string, error)
	GetCollections(project string) ([]*models.Collection, error)
	DeleteCollection(project, name string) error

	// Project collection documents
	AddCollectionDocument(project, collectionName string, document map[string]interface{}) (map[string]interface{}, error)
	UpdateCollectionDocument(project, collectionName, documentID string, updatedDocumet map[string]interface{}) error
	GetCollectionDocuments(project, collectionName string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, error)
	GetCollectionDocument(project, collectionName, documentID string) (map[string]interface{}, error)
	CountCollectionDocuments(project, collectionName string) (int64, error)
	DeleteCollectionDocument(project, collectionName, documentID string) error
	DropAllCollectionDocuments(project, collectionName string) error

	// Project users
	// Project apikeys
	// Project logs
	// Project sessions
	// Projects
	// Users
	// Sessions
}
