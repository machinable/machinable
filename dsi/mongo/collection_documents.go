package mongo

import (
	"context"
	"fmt"

	"bitbucket.org/nsjostrom/machinable/dsi/mongo/database"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// CollectionDocs is the mongoDB implementation of the ProjectCollectionDocuments interface
type CollectionDocs struct {
	db *database.Database
}

// AddDocument creates a new arbitrary argument in the given collection for the project
func (c *CollectionDocs) AddDocument(project, collectionName string, document map[string]interface{}) (map[string]interface{}, error) {
	// Get a connection and insert the new document
	collection := c.db.CollectionDocs(project, collectionName)
	result, err := collection.InsertOne(
		context.Background(),
		document,
	)

	if err != nil {
		return nil, err
	}

	// Grab and set the new document ID
	insertedID, ok := result.InsertedID.(objectid.ObjectID)
	if ok {
		document["id"] = insertedID.Hex()
	} else {
		document["id"] = result.InsertedID
	}

	return document, nil
}

// UpdateDocument updates the entire document for the documentID
func (c *CollectionDocs) UpdateDocument(project, collectionName, documentID string, updatedDocumet map[string]interface{}) error {
	return nil
}

// GetDocuments retrieves the entire list of documents for a collection
func (c *CollectionDocs) GetDocuments(project, collectionName string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, error) {
	collection := c.db.CollectionDocs(project, collectionName)

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		return nil, err
	}

	documents := make([]map[string]interface{}, 0)
	doc := bson.NewDocument()
	for cursor.Next(context.Background()) {
		doc.Reset()
		err := cursor.Decode(doc)
		if err == nil {
			document, err := parseUnknownDocumentToMap(doc, 0)
			if err != nil {

			}
			documents = append(documents, document)
		}
	}

	return documents, nil
}

// GetDocument retrieves a single document from the collection
func (c *CollectionDocs) GetDocument(project, collectionName, documentID string) (map[string]interface{}, error) {
	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil && documentID == "" {
		return nil, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error())
	}

	collection := c.db.CollectionDocs(project, collectionName)

	// Find object based on ID and decode result into document
	doc := bson.NewDocument()
	err = collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		nil,
	).Decode(doc)

	if err != nil {
		return nil, nil
	}

	// Lookup  definitions for this resource
	object, err := parseUnknownDocumentToMap(doc, 0)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (c *CollectionDocs) DropAll(project, collectionName string) error {
	return nil
}
