package mongo

import (
	"context"
	"fmt"

	"github.com/anothrnick/machinable/dsi"
	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// AddDefDocument creates a new document for the existing resource, specified by the path.
func (d *Datastore) AddDefDocument(project, path string, fields models.ResourceObject, metadata *models.MetaData) (string, *errors.DatastoreError) {
	resDefCollection := d.db.ResourceDefinitions(project)
	// Get field definitions for this resource
	resourceDefinition, err := getDefinition(path, resDefCollection)
	if err != nil {
		return "", errors.New(errors.NotFound, fmt.Errorf("resource does not exist"))
	}

	// validate schema
	schemaErr := fields.Validate(resourceDefinition)
	if schemaErr != nil {
		return "", errors.New(errors.BadParameter, schemaErr)
	}

	// Append metadata
	fields["_metadata"] = metadata.Map()

	// Get the resources.{resourcePathName} collection
	rc := d.db.ResourceDocs(project, path)
	result, err := rc.InsertOne(
		context.Background(),
		fields,
	)

	if err != nil {
		return "", errors.New(errors.UnknownError, err)
	}

	return result.InsertedID.(objectid.ObjectID).Hex(), nil
}

// UpdateDefDocument updates an existing document if it exists
func (d *Datastore) UpdateDefDocument(project, path, documentID string, updatedFields models.ResourceObject, filter map[string]interface{}) *errors.DatastoreError {
	// create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil || documentID == "" {
		return errors.New(errors.BadParameter, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error()))
	}

	// get existing object
	collection := d.db.ResourceDocs(project, path)
	_, derr := d.getDefinitionDocument(objectID, collection)
	if derr != nil {
		return errors.New(errors.NotFound, fmt.Errorf("object does not exist '%s'", documentID))
	}

	// get definition
	resDefCollection := d.db.ResourceDefinitions(project)
	// get field definitions for this resource
	resourceDefinition, err := getDefinition(path, resDefCollection)
	if err != nil {
		return errors.New(errors.NotFound, fmt.Errorf("resource does not exist"))
	}

	// validate schema
	schemaErr := updatedFields.Validate(resourceDefinition)
	if schemaErr != nil {
		return errors.New(errors.BadParameter, schemaErr)
	}

	// create updated elements
	updatedElements := make([]*bson.Element, 0)
	for key := range updatedFields {
		if dsi.ReservedField(key) {
			// remove field for PUT
			delete(updatedFields, key)
		} else {
			// append to element slice
			updatedElements = append(updatedElements, bson.EC.Interface(key, updatedFields[key]))
		}
	}

	// TODO update meta data with last updated

	// Get the resources.{resourcePathName} collection
	_, err = collection.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				updatedElements...,
			),
		),
	)

	return errors.New(errors.UnknownError, err)
}

// ListDefDocuments retrieves all definition documents for the give project and path
func (d *Datastore) ListDefDocuments(project, path string, limit, offset int64, filter map[string]interface{}, sort map[string]int) ([]map[string]interface{}, *errors.DatastoreError) {
	collection := d.db.ResourceDocs(project, path)
	limitOpt := findopt.Limit(limit)
	offsetOpt := findopt.Skip(offset)
	sortOpt := findopt.Sort(sort)

	// Find all objects for this resource
	cursor, err := collection.Find(
		context.Background(),
		filter,
		limitOpt,
		offsetOpt,
		sortOpt,
	)

	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	// Create response from documents
	documents := make([]map[string]interface{}, 0)
	doc := bson.NewDocument()
	for cursor.Next(context.Background()) {
		doc.Reset()
		err := cursor.Decode(doc)
		if err != nil {
			return nil, errors.New(errors.UnknownError, err)
		}
		// The document is techically "known" because we have a resource definition, but because
		// we save the data as the types based on the definition, the interface values will marshal
		// to JSON just fine.
		fields, err := parseUnknownDocumentToMap(doc, 0)
		if err != nil {
			return nil, errors.New(errors.UnknownError, err)
		}

		documents = append(documents, fields)
	}

	return documents, nil
}

// GetDefDocument retrieves a single document
func (d *Datastore) GetDefDocument(project, path, documentID string) (map[string]interface{}, *errors.DatastoreError) {
	collection := d.db.ResourceDocs(project, path)
	resDefCollection := d.db.ResourceDefinitions(project)

	// Find the resource definition for this object
	_, err := getDefinition(path, resDefCollection)
	if err != nil {
		return nil, errors.New(errors.NotFound, fmt.Errorf("resource does not exist"))
	}

	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil {
		return nil, errors.New(errors.BadParameter, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error()))
	}

	// Find object based on ID
	// Decode result into document
	doc, err := d.getDefinitionDocument(objectID, collection)

	if err != nil {
		return nil, errors.New(errors.NotFound, fmt.Errorf("object does not exist, '%s'", documentID))
	}

	// Lookup  definitions for this resource
	object, err := parseUnknownDocumentToMap(doc, 0)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	return object, nil
}

// CountDefDocuments returns the count of all documents for a project resource
func (d *Datastore) CountDefDocuments(project, path string, filter map[string]interface{}) (int64, *errors.DatastoreError) {
	collection := d.db.ResourceDocs(project, path)
	cnt, err := collection.CountDocuments(nil, filter, nil)

	return cnt, errors.New(errors.UnknownError, err)
}

// DeleteDefDocument deletes a single document
func (d *Datastore) DeleteDefDocument(project, path, documentID string) *errors.DatastoreError {
	collection := d.db.ResourceDocs(project, path)

	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil {
		return errors.New(errors.BadParameter, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error()))
	}

	// Delete the object
	_, err = collection.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
	)
	if err != nil {
		return errors.New(errors.UnknownError, err)
	}

	return nil
}

// DropAllDefDocuments drops the entire collection of documents
func (d *Datastore) DropAllDefDocuments(project, path string) *errors.DatastoreError {
	return nil
}

func (d *Datastore) getDefinitionDocument(objectID objectid.ObjectID, collection *mongo.Collection) (*bson.Document, error) {
	// Find object based on ID
	// Decode result into document
	doc := bson.NewDocument()
	err := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		nil,
	).Decode(doc)

	return doc, err
}
