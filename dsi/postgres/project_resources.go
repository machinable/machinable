package postgres

import (
	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
)

const (
	tableProjectResources       = "project_resources"
	tableProjectResourceObjects = "project_resource_objects"
)

// AddDefinition creates a new definition
func (d *Database) AddDefinition(project string, def *models.ResourceDefinition) (string, *errors.DatastoreError) {
	return "", nil
}

// UpdateDefinition updates the PARALLEL_READ and PARALLEL_WRITE fields of a definition
func (d *Database) UpdateDefinition(project, definitionID string, def *models.ResourceDefinition) *errors.DatastoreError {
	return nil
}

// ListDefinitions lists all definitions for a project
func (d *Database) ListDefinitions(project string) ([]*models.ResourceDefinition, *errors.DatastoreError) {
	return nil, nil
}

// GetDefinition returns a single definition by ID.
func (d *Database) GetDefinition(project, definitionID string) (*models.ResourceDefinition, *errors.DatastoreError) {
	return nil, nil
}

// GetResourceStats returns stats for a resource collectiond
func (d *Database) GetResourceStats(project, name string) (*models.Stats, *errors.DatastoreError) {
	return nil, nil
}

// GetDefinitionByPathName returns a definition based on `PathName`
func (d *Database) GetDefinitionByPathName(project, pathName string) (*models.ResourceDefinition, *errors.DatastoreError) {
	return nil, nil
}

// DeleteDefinition deletes a definition as well as any data stored for that definition
func (d *Database) DeleteDefinition(project, definitionID string) *errors.DatastoreError {
	return nil
}

// DropProjectResources drops all resource data as well as the definition
func (d *Database) DropProjectResources(project string) *errors.DatastoreError {
	return nil
}

/******************************/
/* PROJECT RESOURCE DOCUMENTS */
/******************************/

// AddDefDocument creates a new document for the existing resource, specified by the path.
func (d *Database) AddDefDocument(project, path string, fields models.ResourceObject, metadata *models.MetaData) (string, *errors.DatastoreError) {
	return "", nil
}

// UpdateDefDocument updates an existing document if it exists
func (d *Database) UpdateDefDocument(project, path, documentID string, updatedFields models.ResourceObject, filter map[string]interface{}) *errors.DatastoreError {
	return nil
}

// ListDefDocuments retrieves all definition documents for the give project and path
func (d *Database) ListDefDocuments(project, path string, limit, offset int64, filter map[string]interface{}, sort map[string]int) ([]map[string]interface{}, *errors.DatastoreError) {
	return nil, nil
}

// GetDefDocument retrieves a single document
func (d *Database) GetDefDocument(project, path, documentID string, filter map[string]interface{}) (map[string]interface{}, *errors.DatastoreError) {
	return nil, nil
}

// CountDefDocuments returns the count of all documents for a project resource
func (d *Database) CountDefDocuments(project, path string, filter map[string]interface{}) (int64, *errors.DatastoreError) {
	return 0, nil
}

// DeleteDefDocument deletes a single document
func (d *Database) DeleteDefDocument(project, path, documentID string, filter map[string]interface{}) *errors.DatastoreError {
	return nil
}

// DropAllDefDocuments drops the entire collection of documents
func (d *Database) DropAllDefDocuments(project, path string) *errors.DatastoreError {
	return nil
}
