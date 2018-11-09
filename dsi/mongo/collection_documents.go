package mongo

import (
	"context"

	"bitbucket.org/nsjostrom/machinable/dsi/mongo/database"
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

	insertedID, ok := result.InsertedID.(objectid.ObjectID)
	if ok {
		document["id"] = insertedID.Hex()
	} else {
		document["id"] = result.InsertedID
	}

	return document, nil
}

func (c *CollectionDocs) UpdateDocument(project, collectionName, documentID string, updatedDocumet map[string]interface{}) error {
	return nil
}

func (c *CollectionDocs) GetDocuments(project, collectionName string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, error) {
	return nil, nil
}

func (c *CollectionDocs) GetDocument(project, collectionName, documentID string) (map[string]interface{}, error) {
	return nil, nil
}

func (c *CollectionDocs) DropAll(project, collectionName string) error {
	return nil
}
