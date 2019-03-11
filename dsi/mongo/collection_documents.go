package mongo

// Functions related to the Project Collection Documents.

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo/findopt"

	"github.com/anothrnick/machinable/dsi"
	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// AddCollectionDocument creates a new arbitrary argument in the given collection for the project
func (d *Datastore) AddCollectionDocument(project, collectionName string, document map[string]interface{}, metadata *models.MetaData) (map[string]interface{}, *errors.DatastoreError) {
	// Append metadata
	document["_metadata"] = metadata.Map()

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
func (d *Datastore) UpdateCollectionDocument(project, collectionName, documentID string, updatedDocument map[string]interface{}, metadata *models.MetaData, filter map[string]interface{}) *errors.DatastoreError {
	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil || documentID == "" {
		return errors.New(errors.BadParameter, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error()))
	}

	// get object, copy metadata
	_, getErr := d.GetCollectionDocument(project, collectionName, documentID, filter)
	if getErr != nil {
		return errors.New(errors.NotFound, fmt.Errorf("object does not exist '%s'", documentID))
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

	// TODO update meta data with last updated

	// Get a connection and update the document
	collection := d.db.CollectionDocs(project, collectionName)
	_, updateErr := collection.UpdateOne(
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

	return errors.New(errors.UnknownError, updateErr)
}

// GetCollectionDocuments retrieves the entire list of documents for a collection
func (d *Datastore) GetCollectionDocuments(project, collectionName string, limit, offset int64, filter map[string]interface{}, sort map[string]int) ([]map[string]interface{}, *errors.DatastoreError) {
	collection := d.db.CollectionDocs(project, collectionName)

	limitOpt := findopt.Limit(limit)
	offsetOpt := findopt.Skip(offset)
	sortOpt := findopt.Sort(sort)

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
func (d *Datastore) GetCollectionDocument(project, collectionName, documentID string, filter map[string]interface{}) (map[string]interface{}, *errors.DatastoreError) {
	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil && documentID == "" {
		return nil, errors.New(errors.BadParameter, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error()))
	}

	collection := d.db.CollectionDocs(project, collectionName)

	filters := map[string]interface{}{
		"_id": objectID,
	}

	// apply filters
	if filter != nil {
		for k, v := range filter {
			filters[k] = v
		}
	}

	// Find object based on ID and decode result into document
	doc := bson.NewDocument()
	err = collection.FindOne(
		nil,
		filters,
		nil,
	).Decode(doc)

	if err != nil {
		return nil, errors.New(errors.NotFound, fmt.Errorf("document does not exist"))
	}

	// Lookup  definitions for this resource
	object, err := parseUnknownDocumentToMap(doc, 0)

	return object, errors.New(errors.UnknownError, err)
}

// CountCollectionDocuments returns the count of all documents for a project collection
func (d *Datastore) CountCollectionDocuments(project, collectionName string, filter map[string]interface{}) (int64, *errors.DatastoreError) {
	collection := d.db.CollectionDocs(project, collectionName)
	cnt, err := collection.CountDocuments(nil, filter, nil)

	return cnt, errors.New(errors.UnknownError, err)
}

// DeleteCollectionDocument removes a single document from the provided collection by `ID`
func (d *Datastore) DeleteCollectionDocument(project, collectionName, documentID string, filter map[string]interface{}) *errors.DatastoreError {
	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(documentID)
	if err != nil {
		return errors.New(errors.BadParameter, fmt.Errorf("invalid identifier '%s': %s", documentID, err.Error()))
	}

	collection := d.db.CollectionDocs(project, collectionName)

	filters := map[string]interface{}{
		"_id": objectID,
	}

	// apply filters
	if filter != nil {
		for k, v := range filter {
			filters[k] = v
		}
	}

	// Delete the object
	_, err = collection.DeleteOne(
		context.Background(),
		filters,
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
