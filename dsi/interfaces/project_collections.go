package interfaces

import (
	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
)

// CollectionsDatastore exposes functions to collections
type CollectionsDatastore interface {
	// Project collections
	AddCollection(project string, collection *models.Collection) *errors.DatastoreError
	GetCollection(project, name string) (*models.Collection, *errors.DatastoreError)
	GetCollectionStats(project, name string) (*models.Stats, *errors.DatastoreError)
	UpdateCollection(project, id string, collection *models.Collection) *errors.DatastoreError
	GetCollections(project string) ([]*models.Collection, *errors.DatastoreError)
	DeleteCollection(project, id string) *errors.DatastoreError
	DropProjectCollections(project string) *errors.DatastoreError

	// Project collection documents
	AddCollectionDocument(project, collectionName string, document map[string]interface{}, metadata *models.MetaData) (map[string]interface{}, *errors.DatastoreError)
	UpdateCollectionDocument(project, collectionName, documentID string, updatedDocument map[string]interface{}, metadata *models.MetaData, filter map[string]interface{}) *errors.DatastoreError
	GetCollectionDocuments(project, collectionName string, limit, offset int64, filter map[string]interface{}, sort map[string]int) ([]map[string]interface{}, *errors.DatastoreError)
	GetCollectionDocument(project, collectionName, documentID string, filter map[string]interface{}) (map[string]interface{}, *errors.DatastoreError)
	CountCollectionDocuments(project, collectionName string, filter map[string]interface{}) (int64, *errors.DatastoreError)
	DeleteCollectionDocument(project, collectionName, documentID string, filter map[string]interface{}) *errors.DatastoreError
	DropAllCollectionDocuments(project, collectionName string) *errors.DatastoreError
}
