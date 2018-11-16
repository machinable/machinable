package mongo

import (
	"context"
	"fmt"

	"bitbucket.org/nsjostrom/machinable/dsi/errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// AddDefDocument creates a new document for the existing resource, specified by the path.
func (d *Datastore) AddDefDocument(project, path string, fields map[string]interface{}) (string, *errors.DatastoreError) {
	resDefCollection := d.db.ResourceDefinitions(project)
	// Get field definitions for this resource
	resourceDefinition, err := getDefinition(path, resDefCollection)
	if err != nil {
		return "", errors.New(errors.NotFound, fmt.Errorf("resource does not exist"))
	}

	// TODO: Create openAPI schema and validate submitted data against it

	// Create document for this resource based on the field definitions
	objectDocument, err := createPropertyDocument(fields, resourceDefinition.Properties, 0)
	if err != nil {
		return "", errors.New(errors.BadParameter, err)
	}

	// Get the resources.{resourcePathName} collection
	rc := d.db.ResourceDocs(project, path)
	result, err := rc.InsertOne(
		context.Background(),
		objectDocument,
	)

	if err != nil {
		return "", errors.New(errors.UnknownError, err)
	}

	return result.InsertedID.(objectid.ObjectID).Hex(), nil
}

// ListDefDocuments retrieves all definition documents for the give project and path
// TODO pagination and filtering
func (d *Datastore) ListDefDocuments(project, path string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, *errors.DatastoreError) {
	collection := d.db.ResourceDocs(project, path)
	documents := make([]map[string]interface{}, 0)

	// Find all objects for this resource
	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	// Create response from documents
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

func (d *Datastore) GetDefDocument(project, path, documentID string) (map[string]interface{}, *errors.DatastoreError) {
	return nil, nil
}

func (d *Datastore) DeleteDefDocument(project, path, documentID string) *errors.DatastoreError {
	return nil
}

func (d *Datastore) DropAllDefDocuments(project, path string) *errors.DatastoreError {
	return nil
}
