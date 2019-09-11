package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	dsiErrors "github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
)

const (
	tableProjectResourceDefinitions = "project_resource_definitions"
	tableProjectResourceObjects     = "project_resource_objects"
)

var objectFilterTranslation = map[string]string{
	"_metadata.creator":      "creator",
	"_metadata.creator_type": "creator_type",
	"_metadata.created":      "created",
}

// AddDefinition creates a new definition
func (d *Database) AddDefinition(projectID string, definition *models.ResourceDefinition) (string, *dsiErrors.DatastoreError) {
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

	return definition.ID, dsiErrors.New(dsiErrors.UnknownError, err)
}

// UpdateDefinition updates the access fields of a definition
func (d *Database) UpdateDefinition(projectID, definitionID string, definition *models.ResourceDefinition) *dsiErrors.DatastoreError {
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

	return dsiErrors.New(dsiErrors.UnknownError, err)
}

// ListDefinitions lists all definitions for a project
func (d *Database) ListDefinitions(projectID string) ([]*models.ResourceDefinition, *dsiErrors.DatastoreError) {
	rows, err := d.db.Query(
		fmt.Sprintf(
			"SELECT id, project_id, name, path_name, parallel_read, parallel_write, \"create\", \"read\", \"update\", \"delete\", schema, created FROM %s WHERE project_id=$1",
			tableProjectResourceDefinitions,
		),
		projectID,
	)
	if err != nil {
		return nil, dsiErrors.New(dsiErrors.UnknownError, err)
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
			return nil, dsiErrors.New(dsiErrors.UnknownError, err)
		}

		definitions = append(definitions, &def)
	}

	return definitions, dsiErrors.New(dsiErrors.UnknownError, rows.Err())
}

// GetDefinition returns a single definition by ID.
func (d *Database) GetDefinition(projectID, definitionID string) (*models.ResourceDefinition, *dsiErrors.DatastoreError) {
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
		return nil, dsiErrors.New(dsiErrors.UnknownError, err)
	}

	return &def, nil
}

// GetResourceStats returns stats for a resource collection
func (d *Database) GetResourceStats(projectID, definitionID string) (*models.Stats, *dsiErrors.DatastoreError) {
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
		return nil, dsiErrors.New(dsiErrors.UnknownError, err)
	}

	// TODO get stats for all objects...

	return &stats, nil
}

// GetDefinitionByPathName returns a definition based on `PathName`
func (d *Database) GetDefinitionByPathName(projectID, pathName string) (*models.ResourceDefinition, *dsiErrors.DatastoreError) {
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
		return nil, dsiErrors.New(dsiErrors.UnknownError, err)
	}

	return &def, nil
}

// DeleteDefinition deletes a definition as well as any data stored for that definition
func (d *Database) DeleteDefinition(projectID, definitionID string) *dsiErrors.DatastoreError {

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
		return dsiErrors.New(dsiErrors.UnknownError, err)
	}

	// delete all objects for resource
	dErr = d.DropDefDocuments(projectID, resource.PathName)

	return dErr
}

// DropProjectResources drops all resource data as well as the definition
func (d *Database) DropProjectResources(projectID string) *dsiErrors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project_id=$1",
			tableProjectResourceDefinitions,
		),
		projectID,
	)
	if err != nil {
		return dsiErrors.New(dsiErrors.UnknownError, err)
	}

	// delete all objects for each resource
	dErr := d.DropProjectDefDocuments(projectID)

	return dErr
}

/******************************/
/* PROJECT RESOURCE DOCUMENTS */
/******************************/

// AddDefDocument creates a new document for the existing resource, specified by the path.
func (d *Database) AddDefDocument(projectID, pathName string, fields models.ResourceObject, metadata *models.MetaData) (string, *dsiErrors.DatastoreError) {
	var id string
	var creatorID interface{}

	if metadata.CreatorType == models.CreatorAPIKey || metadata.CreatorType == models.CreatorUser {
		creatorID = metadata.Creator
	}

	// Get field definitions for this resource
	resourceDefinition, defErr := d.GetDefinitionByPathName(projectID, pathName)
	if defErr != nil {
		return "", dsiErrors.New(dsiErrors.NotFound, fmt.Errorf("resource does not exist"))
	}

	// validate schema
	schemaErr := fields.Validate(resourceDefinition)
	if schemaErr != nil {
		return "", dsiErrors.New(dsiErrors.BadParameter, schemaErr)
	}

	data, der := json.Marshal(fields)
	if der != nil {
		return "", dsiErrors.New(dsiErrors.UnknownError, der)
	}

	err := d.db.QueryRow(
		fmt.Sprintf(
			"INSERT INTO %s (project_id, resource_path, creator_type, creator, created, data) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			tableProjectResourceObjects,
		),
		projectID,
		pathName,
		metadata.CreatorType,
		creatorID,
		time.Now(),
		data,
	).Scan(&id)

	return id, dsiErrors.New(dsiErrors.UnknownError, err)
}

// UpdateDefDocument updates an existing document if it exists
func (d *Database) UpdateDefDocument(projectID, pathName, documentID string, updatedFields models.ResourceObject, filter map[string]interface{}) *dsiErrors.DatastoreError {
	// Get field definitions for this resource
	resourceDefinition, defErr := d.GetDefinitionByPathName(projectID, pathName)
	if defErr != nil {
		return dsiErrors.New(dsiErrors.NotFound, fmt.Errorf("resource does not exist"))
	}

	// validate schema
	schemaErr := updatedFields.Validate(resourceDefinition)
	if schemaErr != nil {
		return dsiErrors.New(dsiErrors.BadParameter, schemaErr)
	}

	data, der := json.Marshal(updatedFields)
	if der != nil {
		return dsiErrors.New(dsiErrors.UnknownError, der)
	}

	_, err := d.db.Exec(
		fmt.Sprintf(
			"UPDATE %s SET data=$1 WHERE id=$2",
			tableProjectResourceObjects,
		),
		data,
		documentID,
	)

	return dsiErrors.New(dsiErrors.UnknownError, err)
}

// ListDefDocuments retrieves all definition documents for the give project and path
func (d *Database) ListDefDocuments(projectID, pathName string, limit, offset int64, filter map[string]interface{}, sort map[string]int) ([]map[string]interface{}, *dsiErrors.DatastoreError) {
	// translate filters
	for key, value := range filter {
		if translated, ok := objectFilterTranslation[key]; ok {
			if _, ok := filter[translated]; !ok {
				filter[translated] = value
			}
			delete(filter, key)
		}
	}

	args := make([]interface{}, 0)
	index := 1

	// query builders
	filterString := make([]string, 0)
	sortString := make([]string, 0)
	pageString := ""

	// path name
	args = append(args, pathName)
	filterString = append(filterString, fmt.Sprintf("resource_path=$%d", index))
	index++
	// projectID
	args = append(args, projectID)
	filterString = append(filterString, fmt.Sprintf("project_id=$%d", index))
	index++

	// valid sort/filter
	validFields := map[string]bool{"creator": true}

	// filters
	filterErr := d.mapToQuery(filter, validFields, &filterString, &args, &index)
	if filterErr != nil {
		return nil, dsiErrors.New(dsiErrors.UnknownError, filterErr)
	}

	// sort
	for key, val := range sort {
		// validate fields
		if _, ok := validFields[key]; !ok {
			// not a valid field, move on
			continue
		}
		direction := "DESC"
		if val > 0 {
			direction = "ASC"
		}
		sortString = append(sortString, fmt.Sprintf("%s %s", key, direction))
	}

	// paginate
	if limit >= 0 {
		args = append(args, limit)
		pageString += fmt.Sprintf(" LIMIT $%d", index)
		index++
	}

	if offset >= 0 {
		args = append(args, offset)
		pageString += fmt.Sprintf(" OFFSET $%d", index)
		index++
	}

	queryFields := "id, creator, creator_type, created, data"
	orderBy := ""
	if len(sortString) > 0 {
		orderBy = fmt.Sprintf(" ORDER BY %s", strings.Join(sortString, ", "))
	}
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s%s%s",
		queryFields,
		tableProjectResourceObjects,
		strings.Join(filterString, " AND "),
		orderBy,
		pageString,
	)

	rows, err := d.db.Query(
		query,
		args...,
	)

	if err != nil {
		return nil, dsiErrors.New(dsiErrors.UnknownError, err)
	}
	defer rows.Close()

	objects := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, creatorType string
		var creatorID sql.NullString
		var created time.Time
		obj := make(map[string]interface{})
		byt := make([]byte, 0)

		err = rows.Scan(
			&id,
			&creatorID,
			&creatorType,
			&created,
			&byt,
		)

		if err != nil {
			return nil, dsiErrors.New(dsiErrors.UnknownError, err)
		}

		err = json.Unmarshal(byt, &obj)
		if err != nil {
			return nil, dsiErrors.New(dsiErrors.UnknownError, err)
		}

		obj["_metadata"] = models.MetaData{
			Created:     created.Unix(),
			Creator:     creatorID.String,
			CreatorType: creatorType,
		}
		obj["id"] = id

		objects = append(objects, obj)
	}

	return objects, dsiErrors.New(dsiErrors.UnknownError, rows.Err())
}

// GetDefDocument retrieves a single document
func (d *Database) GetDefDocument(projectID, path, documentID string, filter map[string]interface{}) (map[string]interface{}, *dsiErrors.DatastoreError) {
	// translate filters
	for key, value := range filter {
		if translated, ok := objectFilterTranslation[key]; ok {
			if _, ok := filter[translated]; !ok {
				filter[translated] = value
			}
			delete(filter, key)
		}
	}

	args := make([]interface{}, 0)
	index := 1

	// query builders
	filterString := make([]string, 0)

	// document id
	args = append(args, documentID)
	filterString = append(filterString, fmt.Sprintf("id=$%d", index))
	index++

	// valid sort/filter
	validFields := map[string]bool{"creator": true}

	// filters
	filterErr := d.mapToQuery(filter, validFields, &filterString, &args, &index)
	if filterErr != nil {
		return nil, dsiErrors.New(dsiErrors.UnknownError, filterErr)
	}

	queryFields := "id, creator, creator_type, created, data"

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		queryFields,
		tableProjectResourceObjects,
		strings.Join(filterString, " AND "),
	)

	var id, creatorType string
	var creatorID sql.NullString
	var created time.Time

	obj := make(map[string]interface{})
	byt := make([]byte, 0)

	err := d.db.QueryRow(
		query,
		args...,
	).Scan(
		&id,
		&creatorID,
		&creatorType,
		&created,
		&byt,
	)

	if err != nil {
		return nil, dsiErrors.New(dsiErrors.NotFound, errors.New("not found"))
	}

	err = json.Unmarshal(byt, &obj)
	if err != nil {
		return nil, dsiErrors.New(dsiErrors.UnknownError, err)
	}

	obj["_metadata"] = models.MetaData{
		Created:     created.Unix(),
		Creator:     creatorID.String,
		CreatorType: creatorType,
	}
	obj["id"] = id

	return obj, nil
}

// CountDefDocuments returns the count of all documents for a project resource
func (d *Database) CountDefDocuments(projectID, pathName string, filter map[string]interface{}) (int64, *dsiErrors.DatastoreError) {
	// translate filters
	for key, value := range filter {
		if translated, ok := objectFilterTranslation[key]; ok {
			if _, ok := filter[translated]; !ok {
				filter[translated] = value
			}
			delete(filter, key)
		}
	}

	args := make([]interface{}, 0)
	index := 1

	// query builders
	filterString := make([]string, 0)

	// path name
	args = append(args, pathName)
	filterString = append(filterString, fmt.Sprintf("resource_path=$%d", index))
	index++
	// projectID
	args = append(args, projectID)
	filterString = append(filterString, fmt.Sprintf("project_id=$%d", index))
	index++

	// valid sort/filter
	validFields := map[string]bool{"creator": true}

	// filters
	filterErr := d.mapToQuery(filter, validFields, &filterString, &args, &index)
	if filterErr != nil {
		return 0, dsiErrors.New(dsiErrors.UnknownError, filterErr)
	}

	queryFields := "count(id)"

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		queryFields,
		tableProjectResourceObjects,
		strings.Join(filterString, " AND "),
	)

	var count int64
	err := d.db.QueryRow(
		query,
		args...,
	).Scan(
		&count,
	)

	if err != nil {
		return 0, dsiErrors.New(dsiErrors.UnknownError, err)
	}

	return count, nil
}

// DeleteDefDocument deletes a single document
func (d *Database) DeleteDefDocument(projectID, path, documentID string, filter map[string]interface{}) *dsiErrors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1",
			tableProjectResourceObjects,
		),
		documentID,
	)

	return dsiErrors.New(dsiErrors.UnknownError, err)
}

// DropDefDocuments drops documents for a resource
func (d *Database) DropDefDocuments(projectID, path string) *dsiErrors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE resource_path=$1 AND project_id=$2",
			tableProjectResourceObjects,
		),
		path,
		projectID,
	)

	return dsiErrors.New(dsiErrors.UnknownError, err)
}

// DropProjectDefDocuments drops the entire collection of documents for a project
func (d *Database) DropProjectDefDocuments(projectID string) *dsiErrors.DatastoreError {
	_, err := d.db.Exec(
		fmt.Sprintf(
			"DELETE FROM %s WHERE project_id=$1",
			tableProjectResourceObjects,
		),
		projectID,
	)

	return dsiErrors.New(dsiErrors.UnknownError, err)
}
