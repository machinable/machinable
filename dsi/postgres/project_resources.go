package postgres

import (
	"database/sql"
	"encoding/json"
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
			"INSERT INTO %s (project_id, name, path_name, parallel_read, parallel_write, \"create\", \"read\", \"update\", \"delete\", schema, created) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
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
			"UPDATE %s SET parallel_read=$1, parallel_write=$2, \"create\"=$3, \"read\"=$4, \"update\"=$5, \"delete\"=$6 WHERE id=$7",
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
			"SELECT id, project_id, name, path_name, parallel_read, parallel_write, \"create\", \"read\", \"update\", \"delete\", schema, created FROM %s WHERE project_id=$1",
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
			"SELECT id, project_id, name, path_name, parallel_read, parallel_write, \"create\", \"read\", \"update\", \"delete\", schema, created FROM %s WHERE id=$1",
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

	// TODO get stats for all objects...

	return &stats, nil
}

// GetDefinitionByPathName returns a definition based on `PathName`
func (d *Database) GetDefinitionByPathName(projectID, pathName string) (*models.ResourceDefinition, *errors.DatastoreError) {
	def := models.ResourceDefinition{}
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, project_id, name, path_name, parallel_read, parallel_write, \"create\", \"read\", \"update\", \"delete\", schema, created FROM %s WHERE path_name=$1",
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

	// get resource to delete objects
	resource, dErr := d.GetDefinition(projectID, definitionID)
	if dErr != nil {
		return dErr
	}

	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1",
			tableProjectResourceDefinitions,
		),
		definitionID,
	)
	if err != nil {
		return errors.New(errors.UnknownError, err)
	}

	// delete all objects for resource
	dErr = d.DropDefDocuments(projectID, resource.PathName)

	return dErr
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
	if err != nil {
		return errors.New(errors.UnknownError, err)
	}

	// delete all objects for each resource
	dErr := d.DropProjectDefDocuments(projectID)

	return dErr
}

/******************************/
/* PROJECT RESOURCE DOCUMENTS */
/******************************/

// AddDefDocument creates a new document for the existing resource, specified by the path.
func (d *Database) AddDefDocument(projectID, pathName string, fields models.ResourceObject, metadata *models.MetaData) (string, *errors.DatastoreError) {
	var id string
	var userID, apiKeyID interface{}

	if metadata.CreatorType == models.CreatorAPIKey {
		apiKeyID = metadata.Creator
	} else if metadata.CreatorType == models.CreatorUser {
		userID = metadata.Creator
	} else {
		apiKeyID = nil
		userID = nil
	}

	// Get field definitions for this resource
	resourceDefinition, defErr := d.GetDefinitionByPathName(projectID, pathName)
	if defErr != nil {
		return "", errors.New(errors.NotFound, fmt.Errorf("resource does not exist"))
	}

	// validate schema
	schemaErr := fields.Validate(resourceDefinition)
	if schemaErr != nil {
		return "", errors.New(errors.BadParameter, schemaErr)
	}

	data, der := json.Marshal(fields)
	if der != nil {
		return "", errors.New(errors.UnknownError, der)
	}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, resource_path, user_id, apikey_id, created, data) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			tableProjectResourceObjects,
		),
		projectID,
		pathName,
		userID,
		apiKeyID,
		time.Now(),
		data,
	).Scan(&id)

	return id, errors.New(errors.UnknownError, err)
}

// UpdateDefDocument updates an existing document if it exists
func (d *Database) UpdateDefDocument(projectID, pathName, documentID string, updatedFields models.ResourceObject, filter map[string]interface{}) *errors.DatastoreError {
	// Get field definitions for this resource
	resourceDefinition, defErr := d.GetDefinitionByPathName(projectID, pathName)
	if defErr != nil {
		return errors.New(errors.NotFound, fmt.Errorf("resource does not exist"))
	}

	// validate schema
	schemaErr := updatedFields.Validate(resourceDefinition)
	if schemaErr != nil {
		return errors.New(errors.BadParameter, schemaErr)
	}

	data, der := json.Marshal(updatedFields)
	if der != nil {
		return errors.New(errors.UnknownError, der)
	}

	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET data=$1 WHERE id=$2",
			tableProjectResourceObjects,
		),
		data,
		documentID,
	)

	return errors.New(errors.UnknownError, err)
}

// ListDefDocuments retrieves all definition documents for the give project and path
func (d *Database) ListDefDocuments(projectID, pathName string, limit, offset int64, filter map[string]interface{}, sort map[string]int) ([]map[string]interface{}, *errors.DatastoreError) {

	// TODO: filter

	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, user_id, apikey_id, created, data FROM %s WHERE resource_path=$1 AND project_id=$2 LIMIT $3 OFFSET $4",
			tableProjectResourceObjects,
		),
		pathName,
		projectID,
		limit,
		offset,
	)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}
	defer rows.Close()

	objects := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, creator, creatorType string
		var userID, apikeyID sql.NullString
		var created time.Time
		obj := make(map[string]interface{})
		byt := make([]byte, 0)

		err = rows.Scan(
			&id,
			&userID,
			&apikeyID,
			&created,
			&byt,
		)

		if err != nil {
			return nil, errors.New(errors.UnknownError, err)
		}

		err = json.Unmarshal(byt, &obj)
		if err != nil {
			return nil, errors.New(errors.UnknownError, err)
		}

		if userID.Valid {
			creator = userID.String
			creatorType = models.CreatorUser
		} else if apikeyID.Valid {
			creator = apikeyID.String
			creatorType = models.CreatorAPIKey
		}

		obj["_metadata"] = models.MetaData{
			Created:     created.Unix(),
			Creator:     creator,
			CreatorType: creatorType,
		}
		obj["id"] = id

		objects = append(objects, obj)
	}

	return objects, errors.New(errors.UnknownError, rows.Err())
}

// GetDefDocument retrieves a single document
func (d *Database) GetDefDocument(projectID, path, documentID string, filter map[string]interface{}) (map[string]interface{}, *errors.DatastoreError) {
	var id, creator, creatorType string
	var userID, apikeyID sql.NullString
	var created time.Time

	obj := make(map[string]interface{})
	byt := make([]byte, 0)

	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT id, user_id, apikey_id, created, data FROM %s WHERE id=$1",
			tableProjectResourceObjects,
		),
		documentID,
	).Scan(
		&id,
		&userID,
		&apikeyID,
		&created,
		&byt,
	)

	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	err = json.Unmarshal(byt, &obj)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	if userID.Valid {
		creator = userID.String
		creatorType = models.CreatorUser
	} else if apikeyID.Valid {
		creator = apikeyID.String
		creatorType = models.CreatorAPIKey
	}

	obj["_metadata"] = models.MetaData{
		Created:     created.Unix(),
		Creator:     creator,
		CreatorType: creatorType,
	}
	obj["id"] = id

	return obj, nil
}

// CountDefDocuments returns the count of all documents for a project resource
func (d *Database) CountDefDocuments(projectID, pathName string, filter map[string]interface{}) (int64, *errors.DatastoreError) {
	var count int64
	err := d.db.QueryRow(
		fmt.Sprintf(
			"SELECT count(id) FROM %s WHERE resource_path=$1 AND project_id=$2",
			tableProjectResourceObjects,
		),
		pathName,
		projectID,
	).Scan(
		&count,
	)

	if err != nil {
		return 0, errors.New(errors.UnknownError, err)
	}

	return count, nil
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

// DropDefDocuments drops documents for a resource
func (d *Database) DropDefDocuments(projectID, path string) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE resource_path=$1 AND project_id=$2",
			tableProjectResourceObjects,
		),
		path,
		projectID,
	)

	return errors.New(errors.UnknownError, err)
}

// DropProjectDefDocuments drops the entire collection of documents for a project
func (d *Database) DropProjectDefDocuments(projectID string) *errors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project_id=$1",
			tableProjectResourceObjects,
		),
		projectID,
	)

	return errors.New(errors.UnknownError, err)
}
