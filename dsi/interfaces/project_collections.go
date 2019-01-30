package interfaces

import (
	"bitbucket.org/nsjostrom/machinable/dsi/errors"
	"bitbucket.org/nsjostrom/machinable/dsi/models"
)

// CollectionsDatastore exposes functions to collections
type CollectionsDatastore interface {
	// Project collections
	AddCollection(project, name string) *errors.DatastoreError
	GetCollection(project, name string) (string, *errors.DatastoreError)
	GetCollections(project string) ([]*models.Collection, *errors.DatastoreError)
	DeleteCollection(project, name string) *errors.DatastoreError

	// Project collection documents
	AddCollectionDocument(project, collectionName string, document map[string]interface{}, metadata *models.MetaData) (map[string]interface{}, *errors.DatastoreError)
	UpdateCollectionDocument(project, collectionName, documentID string, updatedDocumet map[string]interface{}, metadata *models.MetaData) *errors.DatastoreError
	GetCollectionDocuments(project, collectionName string, limit, offset int64, filter map[string]interface{}) ([]map[string]interface{}, *errors.DatastoreError)
	GetCollectionDocument(project, collectionName, documentID string) (map[string]interface{}, *errors.DatastoreError)
	CountCollectionDocuments(project, collectionName string) (int64, *errors.DatastoreError)
	DeleteCollectionDocument(project, collectionName, documentID string) *errors.DatastoreError
	DropAllCollectionDocuments(project, collectionName string) *errors.DatastoreError
}
