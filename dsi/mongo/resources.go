package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/anothrnick/machinable/dsi"
	"github.com/anothrnick/machinable/dsi/errors"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// GetResourceStats returns stats for a resource collectiond
func (d *Datastore) GetResourceStats(project, name string) (*models.Stats, *errors.DatastoreError) {
	collection := d.db.ResourceDocs(project, name)

	reader, err := d.db.GetMongoDatabase().RunCommand(nil, bson.NewDocument(bson.EC.String("collStats", collection.Name())))

	if err != nil {
		// log error, collection does not exist yet
		return &models.Stats{
			Size:  0,
			Count: 0,
		}, nil
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

// AddDefinition creates a new definition
func (d *Datastore) AddDefinition(project string, def *models.ResourceDefinition) (string, *errors.DatastoreError) {
	resDefCollection := d.db.ResourceDefinitions(project)
	if exists := d.definitionExists(def, resDefCollection); exists == true {
		return "", errors.New(errors.BadParameter, fmt.Errorf("resource already exists"))
	}

	// Create document
	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.String("title", def.Title))
	resourceElements = append(resourceElements, bson.EC.String("path_name", def.PathName))
	resourceElements = append(resourceElements, bson.EC.Boolean("parallel_read", def.ParallelRead))
	resourceElements = append(resourceElements, bson.EC.Boolean("parallel_write", def.ParallelWrite))
	resourceElements = append(resourceElements, bson.EC.Boolean("create", def.Create))
	resourceElements = append(resourceElements, bson.EC.Boolean("read", def.Read))
	resourceElements = append(resourceElements, bson.EC.Boolean("update", def.Update))
	resourceElements = append(resourceElements, bson.EC.Boolean("delete", def.Delete))
	resourceElements = append(resourceElements, bson.EC.Time("created", time.Now()))
	resourceElements = append(resourceElements, bson.EC.String("properties", def.Properties))

	// Get a connection and insert the new document
	collection := d.db.ResourceDefinitions(project)
	result, err := collection.InsertOne(
		context.Background(),
		bson.NewDocument(resourceElements...),
	)

	if err != nil {
		return "", errors.New(errors.UnknownError, err)
	}

	return result.InsertedID.(objectid.ObjectID).Hex(), nil
}

// ListDefinitions lists all definitions for a project
func (d *Datastore) ListDefinitions(project string) ([]*models.ResourceDefinition, *errors.DatastoreError) {
	definitions := make([]*models.ResourceDefinition, 0)

	collection := d.db.ResourceDefinitions(project)

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
		def, _ := parseDefinition(doc)
		definitions = append(definitions, def)
	}

	return definitions, nil
}

// GetDefinition returns a single definition by ID.
func (d *Datastore) GetDefinition(project, definitionID string) (*models.ResourceDefinition, *errors.DatastoreError) {
	resDefCollection := d.db.ResourceDefinitions(project)
	def, err := d.getDefinitionByID(definitionID, resDefCollection)
	if err != nil {
		return nil, errors.New(errors.NotFound, err)
	}

	return def, nil
}

// GetDefinitionByPathName returns a definition based on `PathName`
func (d *Datastore) GetDefinitionByPathName(project, pathName string) (*models.ResourceDefinition, *errors.DatastoreError) {
	resDefCollection := d.db.ResourceDefinitions(project)
	// Find the resource definition for this object
	// Decode result into document
	doc := bson.Document{}
	err := resDefCollection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("path_name", pathName),
		),
		nil,
	).Decode(&doc)

	if err != nil {
		return nil, errors.New(errors.NotFound, fmt.Errorf("resource not found"))
	}

	// Lookup field definitions for this resource
	resourceDefinition, err := parseDefinition(&doc)
	if err != nil {
		return nil, errors.New(errors.UnknownError, err)
	}

	return resourceDefinition, nil
}

// UpdateDefinition updates the PARALLEL_READ and PARALLEL_WRITE fields of a definition
func (d *Datastore) UpdateDefinition(project, definitionID string, def *models.ResourceDefinition) *errors.DatastoreError {
	resDefCollection := d.db.ResourceDefinitions(project)

	// Get the object id for collection name
	objectID, err := objectid.FromHex(definitionID)
	if err != nil {
		return errors.New(errors.BadParameter, err)
	}

	// Find object based on ID
	// decode result into document
	doc := bson.NewDocument()
	err = resDefCollection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		nil,
	).Decode(doc)

	if err != nil {
		return errors.New(errors.NotFound, fmt.Errorf("resource not found, '%s'", definitionID))
	}

	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.Boolean("parallel_read", def.ParallelRead))
	resourceElements = append(resourceElements, bson.EC.Boolean("parallel_write", def.ParallelWrite))
	resourceElements = append(resourceElements, bson.EC.Boolean("create", def.Create))
	resourceElements = append(resourceElements, bson.EC.Boolean("read", def.Read))
	resourceElements = append(resourceElements, bson.EC.Boolean("update", def.Update))
	resourceElements = append(resourceElements, bson.EC.Boolean("delete", def.Delete))

	_, err = resDefCollection.UpdateOne(
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

// DeleteDefinition deletes a definition as well as any data stored for that definition
func (d *Datastore) DeleteDefinition(project, definitionID string) *errors.DatastoreError {
	resDefCollection := d.db.ResourceDefinitions(project)
	// Get definition for resource name
	def, err := d.getDefinitionByID(definitionID, resDefCollection)
	if err != nil {
		return errors.New(errors.NotFound, err)
	}
	resourcePathName := def.PathName

	// Get the object id
	objectID, err := objectid.FromHex(definitionID)
	if err != nil {
		return errors.New(errors.UnknownError, err)
	}

	// Delete the definition
	_, err = resDefCollection.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
	)
	if err != nil {
		return errors.New(errors.UnknownError, err)
	}

	resourceCollection := d.db.ResourceDocs(project, resourcePathName)
	err = resourceCollection.Drop(nil, nil)
	if err != nil {
		return errors.New(errors.UnknownError, err)
	}

	return nil
}

// DropProjectResources drops all resource data as well as the definition
func (d *Datastore) DropProjectResources(project string) *errors.DatastoreError {
	defs, err := d.ListDefinitions(project)
	if err != nil {
		return err
	}

	for _, def := range defs {
		err := d.DeleteDefinition(project, def.ID)
		if err != nil {
			return err
		}
	}

	// drop definition storage
	collection := d.db.ResourceDefinitions(project)

	dropErr := collection.Drop(nil, nil)

	return errors.New(errors.UnknownError, dropErr)
}

// definitionExists returns true if a definition already exists with path_name or name
func (d *Datastore) definitionExists(definition *models.ResourceDefinition, collection *mongo.Collection) bool {
	// Find the resource definition for this object
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ArrayFromElements("$or",
				bson.VC.DocumentFromElements(
					bson.EC.String("path_name", definition.PathName),
				),
				bson.VC.DocumentFromElements(
					bson.EC.String("title", definition.Title),
				),
			),
		),
		nil,
	)

	// Decode result into document
	doc := bson.Document{}
	err := documentResult.Decode(&doc)
	if err != nil {
		return false
	}

	return true
}

// getDefinitionByID returns the *model.ResourceDefinition by resource definition ID
func (d *Datastore) getDefinitionByID(resourceID string, collection *mongo.Collection) (*models.ResourceDefinition, error) {
	objectID, err := objectid.FromHex(resourceID)
	if err != nil {
		return nil, err
	}

	// Find the resource definition for this object
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID(dsi.DocumentIDKey, objectID),
		),
		nil,
	)

	if documentResult == nil {
		return nil, fmt.Errorf("no documents for resource")
	}

	// Decode result into document
	doc := bson.Document{}
	err = documentResult.Decode(&doc)
	if err != nil {
		return nil, fmt.Errorf("no documents for resource")
	}
	// Lookup field definitions for this resource
	resourceDefinition, err := parseDefinition(&doc)
	if err != nil {
		return nil, err
	}

	return resourceDefinition, nil
}
