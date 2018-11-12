package mongo

// Functions related to the Project Collection Documents.

import (
	"context"
	"fmt"

	"bitbucket.org/nsjostrom/machinable/dsi"
	"bitbucket.org/nsjostrom/machinable/dsi/errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// AddCollectionDocument creates a new arbitrary argument in the given collection for the project
func (d *Datastore) AddCollectionDocument(project, collectionName string, document map[string]interface{}) (map[string]interface{}, *errors.DatastoreError) {
	// Get a connection and insert the new document
	collection := d.db.CollectionDocs(project, collectionName)
	result, err := collection.InsertOne(
		context.Background(),
		document,
	)

	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
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

// UpdateCollectionDocument updates the entire document for the documentID, removing any reserved fields with `dsi.ReservedField`
func (d *Datastore) UpdateCollectionDocument(project, collectionName, documentID string, updatedDocument map[string]interface{}) *errors.DatastoreError {
	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil || documentID == "" {
		return errors.New(errors.BadParameter, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error()))
	}

	// iterate over root keys for reserved fields
	updatedElements := make([]*bson.Element, 0)
	for key := range updatedDocument {
		if dsi.ReservedField(key) {
			// remove field for PUT
			delete(updatedDocument, key)
		} else {
			// append to element slice
			updatedElements = append(updatedElements, bson.EC.Interface(key, updatedDocument[key]))
		}
	}

	// Get a connection and update the document
	collection := d.db.CollectionDocs(project, collectionName)
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

// GetCollectionDocuments retrieves the entire list of documents for a collection
func (d *Datastore) GetCollectionDocuments(project, collectionName string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, *errors.DatastoreError) {
	collection := d.db.CollectionDocs(project, collectionName)

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
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

// GetCollectionDocument retrieves a single document from the collection
func (d *Datastore) GetCollectionDocument(project, collectionName, documentID string) (map[string]interface{}, *errors.DatastoreError) {
	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil && documentID == "" {
		return nil, errors.New(errors.BadParameter, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error()))
	}

	collection := d.db.CollectionDocs(project, collectionName)

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

	return object, errors.New(errors.UnknownError, err)
}

// CountCollectionDocuments returns the count of all documents for a project collection
func (d *Datastore) CountCollectionDocuments(project, collectionName string) (int64, *errors.DatastoreError) {
	collection := d.db.CollectionDocs(project, collectionName)
	cnt, err := collection.CountDocuments(nil, nil, nil)

	return cnt, errors.New(errors.UnknownError, err)
}

// DeleteCollectionDocument removes a single document from the provided collection by `ID`
func (d *Datastore) DeleteCollectionDocument(project, collectionName, documentID string) *errors.DatastoreError {
	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil {
		return errors.New(errors.BadParameter, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error()))
	}

	collection := d.db.CollectionDocs(project, collectionName)

	// Delete the object
	_, err = collection.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
	)

	return errors.New(errors.UnknownError, err)
}

// DropAllCollectionDocuments drops the entire collection
func (d *Datastore) DropAllCollectionDocuments(project, collectionName string) *errors.DatastoreError {
	// drop the collection docs (actual data for this collection)
	collection := d.db.CollectionDocs(project, collectionName)

	err := collection.Drop(nil, nil)

	return errors.New(errors.UnknownError, err)
}
