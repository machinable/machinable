package mongo

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/nsjostrom/machinable/dsi"
	"bitbucket.org/nsjostrom/machinable/dsi/errors"
	"bitbucket.org/nsjostrom/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// AddDefinition creates a new definition
func (d *Datastore) AddDefinition(project string, def *models.ResourceDefinition) (string, *errors.DatastoreError) {
	resDefCollection := d.db.ResourceDefinitions(project)
	if exists := d.definitionExists(def, resDefCollection); exists == true {
		return "", errors.New(errors.BadParameter, fmt.Errorf("resource already exists"))
	}

	// Process the resource fields into bson
	propertyElements, err := d.processProperties(def.Properties, 0)
	if err != nil {
		return "", errors.New(errors.BadParameter, err)
	}

	// Create document
	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.String("title", def.Title))
	resourceElements = append(resourceElements, bson.EC.String("path_name", def.PathName))
	resourceElements = append(resourceElements, bson.EC.Time("created", time.Now()))
	resourceElements = append(resourceElements, bson.EC.SubDocumentFromElements("properties", propertyElements...))

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

// processProperties processes a map[string]string to a slice of *bson.Element for storing in mongo
func (d *Datastore) processProperties(properties map[string]models.Property, layer int) ([]*bson.Element, error) {
	propertyElements := make([]*bson.Element, 0)

	// this object goes deeper than supported
	if layer > dsi.MaxRecursion {
		return propertyElements, nil
	}

	for key, prop := range properties {
		if !dsi.SupportedType(prop.Type) {
			return nil, fmt.Errorf("'%s' is not a supported property type", prop.Type)
		}
		if dsi.ReservedField(key) {
			return nil, fmt.Errorf("'%s' is a reserved property name", key)
		}
		itemsType := ""
		if prop.Items != nil {
			itemsType = prop.Items.Type
			if !dsi.SupportedArrayType(itemsType) {
				return nil, fmt.Errorf("'%s' is not a supported property items.type", itemsType)
			}

			// TODO: prop.Items.Type is 'object', process prop.Items.Properties (recursive)

			// NOTE: how do we support arrays of arrays?
		}

		properties := make([]*bson.Element, 0)

		if prop.Properties != nil {
			// Process the resource fields into bson
			var err error
			properties, err = d.processProperties(prop.Properties, layer+1)
			if err != nil {
				return nil, fmt.Errorf("could not process property's properties")
			}
		}

		propertyElements = append(
			propertyElements,
			bson.EC.SubDocument(key, bson.NewDocument(
				bson.EC.String("type", prop.Type),
				bson.EC.String("description", prop.Description),
				bson.EC.String("format", prop.Format),
				bson.EC.SubDocument("items", bson.NewDocument(
					bson.EC.String("type", itemsType),
				)),
				bson.EC.SubDocumentFromElements("properties", properties...),
			)))
	}

	return propertyElements, nil
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
