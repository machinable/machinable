package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// GetCollectionStats returns stats for a collection
func (d *Datastore) GetCollectionStats(project, name string) (*models.Stats, *errors.DatastoreError) {
	collection := d.db.CollectionDocs(project, name)

	reader, err := d.db.GetMongoDatabase().RunCommand(nil, bson.NewDocument(bson.EC.String("collStats", collection.Name())))

	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	stats := &struct {
		Size  int64 `bson:"size"`
		Count int64 `bson:"count"`
	}{}

	err = bson.Unmarshal(reader, stats)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	return &models.Stats{
		Size:  stats.Size,
		Count: stats.Count,
	}, nil
}

// AddCollection creates a new project collection
func (d *Datastore) AddCollection(project string, newCol *models.Collection) *errors.DatastoreError {
	// Create document
	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.String("name", newCol.Name))
	resourceElements = append(resourceElements, bson.EC.Boolean("parallel_read", newCol.ParallelRead))
	resourceElements = append(resourceElements, bson.EC.Boolean("parallel_write", newCol.ParallelWrite))
	resourceElements = append(resourceElements, bson.EC.Time("created", time.Now()))

	// Get a connection and insert the new document
	collection := d.db.CollectionNames(project)
	_, err := collection.InsertOne(
		context.Background(),
		bson.NewDocument(resourceElements...),
	)

	return errors.New(errors.UnknownError, err)
}

// GetCollection retrieves a collection by name, confirming the collection exists
func (d *Datastore) GetCollection(project, name string) (*models.Collection, *errors.DatastoreError) {
	collection := d.db.CollectionNames(project)
	doc := bson.Document{}
	colModel := &models.Collection{}

	// Find the resource definition for this object
	err := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("name", name),
		),
		nil,
	).Decode(&doc)

	if err == nil {
		colModel = &models.Collection{
			Name:          doc.Lookup("name").StringValue(),
			ID:            doc.Lookup("_id").ObjectID().Hex(),
			ParallelRead:  doc.Lookup("parallel_read").Boolean(),
			ParallelWrite: doc.Lookup("parallel_write").Boolean(),
			Created:       doc.Lookup("created").Time(),
			Items:         0,
		}
	}

	return colModel, errors.New(errors.UnknownError, err)
}

// UpdateCollection updates a collection by id
func (d *Datastore) UpdateCollection(project, collectionID string, read, write bool) *errors.DatastoreError {
	collection := d.db.CollectionNames(project)
	// Get the object id for collection name
	objectID, err := objectid.FromHex(collectionID)
	if err != nil {
		return errors.New(errors.BadParameter, err)
	}

	// Find object based on ID
	// decode result into document
	doc := bson.NewDocument()
	err = collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		nil,
	).Decode(doc)

	if err != nil {
		return errors.New(errors.NotFound, fmt.Errorf("collection not found, '%s'", collectionID))
	}

	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.Boolean("parallel_read", read))
	resourceElements = append(resourceElements, bson.EC.Boolean("parallel_write", write))

	_, err = collection.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				resourceElements...,
			),
		),
	)

	return errors.New(errors.UnknownError, err)
}

// GetCollections retrieves the full list of collections, by name, for a project. Each item has the count of documents within the collection.
func (d *Datastore) GetCollections(project string) ([]*models.Collection, *errors.DatastoreError) {
	collections := make([]*models.Collection, 0)

	collection := d.db.CollectionNames(project)

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	doc := bson.NewDocument()
	for cursor.Next(context.Background()) {
		doc.Reset()
		err := cursor.Decode(doc)
		if err != nil {
			return nil, errors.New(errors.UnknownError, err)
		}
		// get count
		docName := doc.Lookup("name").StringValue()
		cnt, _ := d.CountCollectionDocuments(project, docName, nil)

		collections = append(collections,
			&models.Collection{
				Name:          doc.Lookup("name").StringValue(),
				ID:            doc.Lookup("_id").ObjectID().Hex(),
				ParallelRead:  doc.Lookup("parallel_read").Boolean(),
				ParallelWrite: doc.Lookup("parallel_write").Boolean(),
				Created:       doc.Lookup("created").Time(),
				Items:         cnt,
			})
	}

	return collections, nil
}

// DeleteCollection deletes the collection document, as well as all documents for the collection
func (d *Datastore) DeleteCollection(project, collectionID string) *errors.DatastoreError {
	collections := d.db.CollectionNames(project)
	// Get the object id for collection name
	objectID, err := objectid.FromHex(collectionID)
	if err != nil {
		return errors.New(errors.BadParameter, err)
	}

	// Find object based on ID
	// decode result into document
	doc := bson.NewDocument()
	err = collections.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		nil,
	).Decode(doc)

	if err != nil {
		return errors.New(errors.NotFound, fmt.Errorf("collection not found, '%s'", collectionID))
	}

	// get collection name
	collectionName := doc.Lookup("name").StringValue()

	// delete the collection
	_, err = collections.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
	)
	if err != nil {
		return errors.New(errors.UnknownError, err)
	}

	// drop dollection
	dropErr := d.DropAllCollectionDocuments(project, collectionName)

	return dropErr
}

// DropProjectCollections deletes all project collections and their data
func (d *Datastore) DropProjectCollections(project string) *errors.DatastoreError {
	collections, err := d.GetCollections(project)
	if err != nil {
		return err
	}

	for _, collection := range collections {
		err := d.DeleteCollection(project, collection.ID)
		if err != nil {
			return err
		}
	}

	// drop collection
	collection := d.db.CollectionNames(project)

	dropErr := collection.Drop(nil, nil)
	return errors.New(errors.UnknownError, dropErr)
}
