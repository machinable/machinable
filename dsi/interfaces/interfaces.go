package interfaces

import (
	"bitbucket.org/nsjostrom/machinable/dsi/errors"
	"bitbucket.org/nsjostrom/machinable/dsi/models"
)

// Datastore exposes the necessary functions to interact with the Machinable datastore.
// Functions are grouped logically based on their purpose and the collections they interact with.
type Datastore interface {
	// Project definition documents
	AddDefDocument(project, path string, fields map[string]interface{}) (string, *errors.DatastoreError)
	ListDefDocuments(project, path string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, *errors.DatastoreError)
	GetDefDocument(project, path, documentID string) (map[string]interface{}, *errors.DatastoreError)
	DeleteDefDocument(project, path, documentID string) *errors.DatastoreError
	DropAllDefDocuments(project, path string) *errors.DatastoreError

	// Project resource definitions
	AddDefinition(project string, def *models.ResourceDefinition) (string, *errors.DatastoreError)
	ListDefinitions(project string) ([]*models.ResourceDefinition, *errors.DatastoreError)
	GetDefinition(project, definitionID string) (*models.ResourceDefinition, *errors.DatastoreError)
	DeleteDefinition(project, definitionID string) *errors.DatastoreError

	// Project collections
	AddCollection(project, name string) *errors.DatastoreError
	GetCollection(project, name string) (string, *errors.DatastoreError)
	GetCollections(project string) ([]*models.Collection, *errors.DatastoreError)
	DeleteCollection(project, name string) *errors.DatastoreError

	// Project collection documents
	AddCollectionDocument(project, collectionName string, document map[string]interface{}) (map[string]interface{}, *errors.DatastoreError)
	UpdateCollectionDocument(project, collectionName, documentID string, updatedDocumet map[string]interface{}) *errors.DatastoreError
	GetCollectionDocuments(project, collectionName string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, *errors.DatastoreError)
	GetCollectionDocument(project, collectionName, documentID string) (map[string]interface{}, *errors.DatastoreError)
	CountCollectionDocuments(project, collectionName string) (int64, *errors.DatastoreError)
	DeleteCollectionDocument(project, collectionName, documentID string) *errors.DatastoreError
	DropAllCollectionDocuments(project, collectionName string) *errors.DatastoreError

	// Project users
	// Project apikeys
	// Project logs
	// Project sessions
	// Projects
	// Users
	// Sessions
}
