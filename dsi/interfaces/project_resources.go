package interfaces

import (
	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
)

// ResourcesDatastore exposes functions to the resources and definitions
type ResourcesDatastore interface {
	// Project resource definitions
	AddDefinition(project string, def *models.ResourceDefinition) (string, *errors.DatastoreError)
	UpdateDefinition(project, definitionID string, read, write bool) *errors.DatastoreError
	ListDefinitions(project string) ([]*models.ResourceDefinition, *errors.DatastoreError)
	GetDefinition(project, definitionID string) (*models.ResourceDefinition, *errors.DatastoreError)
	GetDefinitionByPathName(project, pathName string) (*models.ResourceDefinition, *errors.DatastoreError)
	DeleteDefinition(project, definitionID string) *errors.DatastoreError

	// Project definition documents
	AddDefDocument(project, path string, fields models.ResourceObject, metadata *models.MetaData) (string, *errors.DatastoreError)
	UpdateDefDocument(project, path, documentID string, updatedFields models.ResourceObject, filter map[string]interface{}) *errors.DatastoreError
	ListDefDocuments(project, path string, limit, offset int64, filter map[string]interface{}, sort map[string]int) ([]map[string]interface{}, *errors.DatastoreError)
	GetDefDocument(project, path, documentID string) (map[string]interface{}, *errors.DatastoreError)
	CountDefDocuments(project, path string, filter map[string]interface{}) (int64, *errors.DatastoreError)
	DeleteDefDocument(project, path, documentID string) *errors.DatastoreError
	DropAllDefDocuments(project, path string) *errors.DatastoreError
}
