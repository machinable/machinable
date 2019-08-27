package postgres

import (
	"fmt"
	"time"

	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
)

const (
	tableProjectResourceDefinitions = "project_resource_definitions"
	tableProjectResourceObjects     = "project_resource_objects"
)

// AddDefinition creates a new definition
func (d *Database) AddDefinition(projectID string, definition *models.ResourceDefinition) (string, *errors.DatastoreError) {
	err := d.db.QueryRow(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, name, path_name, parallel_read, parallel_write, create, read, update, delete, schema, created) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
			tableProjectResourceDefinitions,
		),
		projectID,
		definition.Title,
		definition.PathName,
		definition.ParallelRead,
		definition.ParallelWrite,
		definition.Create,
		definition.Read,
		definition.Update,
		definition.Delete,
		definition.Schema,
		time.Now(),
	).Scan(&definition.ID)

	return definition.ID, errors.New(errors.UnknownError, err)
}

// UpdateDefinition updates the access fields of a definition
func (d *Database) UpdateDefinition(projectID, definitionID string, definition *models.ResourceDefinition) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET parallel_read=$1, parallel_write=$2, create=$3, read=$4, update=$5, delete=$6 WHERE id=$7",
			tableProjectResourceDefinitions,
		),
		definition.ParallelRead,
		definition.ParallelWrite,
		definition.Create,
		definition.Read,
		definition.Update,
		definition.Delete,
		definitionID,
	)

	return errors.New(errors.UnknownError, err)
}

// ListDefinitions lists all definitions for a project
func (d *Database) ListDefinitions(projectID string) ([]*models.ResourceDefinition, *errors.DatastoreError) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, project_id, name, path_name, parallel_read, parallel_write, create, read, update, delete, schema, created FROM %s WHERE project_id=$1",
			tableProjectResourceDefinitions,
		),
		projectID,
	)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}
	defer rows.Close()

	definitions := make([]*models.ResourceDefinition, 0)
	for rows.Next() {
		def := models.ResourceDefinition{}
		err = rows.Scan(
			&def.ID,
			&def.ProjectID,
			&def.Title,
			&def.PathName,
			&def.ParallelRead,
			&def.ParallelWrite,
			&def.Create,
			&def.Read,
			&def.Update,
			&def.Delete,
			&def.Schema,
			&def.Created,
		)
		if err != nil {
			return nil, errors.New(errors.UnknownError, err)
		}

		definitions = append(definitions, &def)
	}

	return definitions, errors.New(errors.UnknownError, rows.Err())
}

// GetDefinition returns a single definition by ID.
func (d *Database) GetDefinition(projectID, definitionID string) (*models.ResourceDefinition, *errors.DatastoreError) {
	def := models.ResourceDefinition{}
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, project_id, name, path_name, parallel_read, parallel_write, create, read, update, delete, schema, created FROM %s WHERE id=$1",
			tableProjectResourceDefinitions,
		),
		definitionID,
	).Scan(
		&def.ID,
		&def.ProjectID,
		&def.Title,
		&def.PathName,
		&def.ParallelRead,
		&def.ParallelWrite,
		&def.Create,
		&def.Read,
		&def.Update,
		&def.Delete,
		&def.Schema,
		&def.Created,
	)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	return &def, nil
}

// GetResourceStats returns stats for a resource collection
func (d *Database) GetResourceStats(projectID, definitionID string) (*models.Stats, *errors.DatastoreError) {
	stats := models.Stats{}
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT sum(pg_column_size(%s)), count(*) FROM %s WHERE id=$1",
			tableProjectResourceDefinitions,
			tableProjectResourceDefinitions,
		),
		definitionID,
	).Scan(
		&stats.Size,
		&stats.Count,
	)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	return &stats, nil
}

// GetDefinitionByPathName returns a definition based on `PathName`
func (d *Database) GetDefinitionByPathName(projectID, pathName string) (*models.ResourceDefinition, *errors.DatastoreError) {
	def := models.ResourceDefinition{}
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, project_id, name, path_name, parallel_read, parallel_write, create, read, update, delete, schema, created FROM %s WHERE path_name=$1",
			tableProjectResourceDefinitions,
		),
		pathName,
	).Scan(
		&def.ID,
		&def.ProjectID,
		&def.Title,
		&def.PathName,
		&def.ParallelRead,
		&def.ParallelWrite,
		&def.Create,
		&def.Read,
		&def.Update,
		&def.Delete,
		&def.Schema,
		&def.Created,
	)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	return &def, nil
}

// DeleteDefinition deletes a definition as well as any data stored for that definition
func (d *Database) DeleteDefinition(projectID, definitionID string) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1",
			tableProjectResourceDefinitions,
		),
		definitionID,
	)

	// TODO: delete all objects for resource

	return errors.New(errors.UnknownError, err)
}

// DropProjectResources drops all resource data as well as the definition
func (d *Database) DropProjectResources(projectID string) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project_id=$1",
			tableProjectResourceDefinitions,
		),
		projectID,
	)

	// TODO: delete all objects for each resource

	return errors.New(errors.UnknownError, err)
}

/******************************/
/* PROJECT RESOURCE DOCUMENTS */
/******************************/

// AddDefDocument creates a new document for the existing resource, specified by the path.
func (d *Database) AddDefDocument(projectID, path string, fields models.ResourceObject, metadata *models.MetaData) (string, *errors.DatastoreError) {
	return "", nil
}

// UpdateDefDocument updates an existing document if it exists
func (d *Database) UpdateDefDocument(projectID, path, documentID string, updatedFields models.ResourceObject, filter map[string]interface{}) *errors.DatastoreError {
	return nil
}

// ListDefDocuments retrieves all definition documents for the give project and path
func (d *Database) ListDefDocuments(projectID, path string, limit, offset int64, filter map[string]interface{}, sort map[string]int) ([]map[string]interface{}, *errors.DatastoreError) {
	return nil, nil
}

// GetDefDocument retrieves a single document
func (d *Database) GetDefDocument(projectID, path, documentID string, filter map[string]interface{}) (map[string]interface{}, *errors.DatastoreError) {
	return nil, nil
}

// CountDefDocuments returns the count of all documents for a project resource
func (d *Database) CountDefDocuments(projectID, path string, filter map[string]interface{}) (int64, *errors.DatastoreError) {
	return 0, nil
}

// DeleteDefDocument deletes a single document
func (d *Database) DeleteDefDocument(projectID, path, documentID string, filter map[string]interface{}) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1",
			tableProjectResourceObjects,
		),
		documentID,
	)

	return errors.New(errors.UnknownError, err)
}

// DropAllDefDocuments drops the entire collection of documents
func (d *Database) DropAllDefDocuments(projectID, path string) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE resource_path=$1",
			tableProjectResourceObjects,
		),
		path,
	)

	return errors.New(errors.UnknownError, err)
}
