package interfaces

import (
	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
)

// ResourcesDatastore exposes functions to the resources and definitions
type ResourcesDatastore interface {
	// Project resource definitions
	AddDefinition(projectID string, def *models.ResourceDefinition) (string, *errors.DatastoreError)
	UpdateDefinition(projectID, definitionID string, def *models.ResourceDefinition) *errors.DatastoreError
	ListDefinitions(projectID string) ([]*models.ResourceDefinition, *errors.DatastoreError)
	GetDefinition(projectID, definitionID string) (*models.ResourceDefinition, *errors.DatastoreError)
	GetResourceStats(projectID, pathName string) (*models.Stats, *errors.DatastoreError)
	GetDefinitionByPathName(projectID, pathName string) (*models.ResourceDefinition, *errors.DatastoreError)
	DeleteDefinition(projectID, definitionID string) *errors.DatastoreError
	DropProjectResources(projectID string) *errors.DatastoreError

	// Project definition documents
	AddDefDocument(projectID, path string, fields models.ResourceObject, metadata *models.MetaData) (string, *errors.DatastoreError)
	UpdateDefDocument(projectID, path, documentID string, updatedFields models.ResourceObject, filter map[string]interface{}) *errors.DatastoreError
	ListDefDocuments(projectID, path string, limit, offset int64, filter map[string]interface{}, sort map[string]int) ([]map[string]interface{}, *errors.DatastoreError)
	GetDefDocument(projectID, path, documentID string, filter map[string]interface{}) (map[string]interface{}, *errors.DatastoreError)
	CountDefDocuments(projectID, path string, filter map[string]interface{}) (int64, *errors.DatastoreError)
	DeleteDefDocument(projectID, path, documentID string, filter map[string]interface{}) *errors.DatastoreError
	DropDefDocuments(projectID, path string) *errors.DatastoreError
	DropProjectDefDocuments(projectID string) *errors.DatastoreError
}
