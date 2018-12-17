package interfaces

import (
	"bitbucket.org/nsjostrom/machinable/dsi/errors"
	"bitbucket.org/nsjostrom/machinable/dsi/models"
)

// Datastore exposes the necessary functions to interact with the Machinable datastore.
// Functions are grouped logically based on their purpose and the collections they interact with.
// implemented connectors: MongoDB
// potential connectors: InfluxDB, Postgres JSON, Redis, CouchDB, etc.
type Datastore interface {
	// Project resources/definitions
	ResourcesDatastore
	// Project collections
	CollectionsDatastore
	// Project users
	ProjectUsersDatastore
	// Project apikeys
	ProjectAPIKeysDatastore
	// Project logs
	ProjectLogsDatastore
	// Project sessions
	ProjectSessionsDatastore
	// Projects

	// Users

	// Sessions

}

// ProjectAPIKeysDatastore exposes functions to manage project api keys
type ProjectAPIKeysDatastore interface {
	CreateAPIKey(project, hash, description string, read, write bool) (*models.ProjectAPIKey, error)
	ListAPIKeys(project string) ([]*models.ProjectAPIKey, error)
	DeleteAPIKey(project, keyID string) error
}

// ProjectSessionsDatastore exposes functions to manage project user sessions
type ProjectSessionsDatastore interface {
	CreateSession(project string, session *models.Session) error
	GetSession(project, sessionID string) (*models.Session, error)
	ListSessions(project string) ([]*models.Session, error)
	DeleteSession(project, sessionID string) error
}

// ProjectUsersDatastore exposes functions to manage project users
type ProjectUsersDatastore interface {
	GetUserByUsername(project, userName string) (*models.ProjectUser, error)
	CreateUser(project string, user *models.ProjectUser) error
	ListUsers(project string) ([]*models.ProjectUser, error)
	DeleteUser(project, userID string) error
}

// ProjectLogsDatastore exposes functions to the project access logs
type ProjectLogsDatastore interface {
	AddProjectLog(project string, log *models.Log) error
	GetProjectLogsForLastHours(project string, hours int) ([]*models.Log, error)
}

// ResourcesDatastore exposes functions to the resources and definitions
type ResourcesDatastore interface {
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
}

// CollectionsDatastore exposes functions to collections
type CollectionsDatastore interface {
	// Project collections
	AddCollection(project, name string) *errors.DatastoreError
	GetCollection(project, name string) (string, *errors.DatastoreError)
	GetCollections(project string) ([]*models.Collection, *errors.DatastoreError)
	DeleteCollection(project, name string) *errors.DatastoreError

	// Project collection documents
	AddCollectionDocument(project, collectionName string, document map[string]interface{}) (map[string]interface{}, *errors.DatastoreError)
	UpdateCollectionDocument(project, collectionName, documentID string, updatedDocumet map[string]interface{}) *errors.DatastoreError
	GetCollectionDocuments(project, collectionName string, limit, offset int64, filter map[string]interface{}) ([]map[string]interface{}, *errors.DatastoreError)
	GetCollectionDocument(project, collectionName, documentID string) (map[string]interface{}, *errors.DatastoreError)
	CountCollectionDocuments(project, collectionName string) (int64, *errors.DatastoreError)
	DeleteCollectionDocument(project, collectionName, documentID string) *errors.DatastoreError
	DropAllCollectionDocuments(project, collectionName string) *errors.DatastoreError
}
