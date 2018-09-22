package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/nsjostrom/machinable/database"
	"bitbucket.org/nsjostrom/machinable/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// supportedTypes is the list of supported resource field types, this will include any other
// resource definitions that have been created ("foreign key")
var supportedTypes = []string{"int", "float", "date", "bool", "string"}

func supportedType(a string) bool {
	for _, b := range supportedTypes {
		if b == a {
			return true
		}
	}
	return false
}

func validateResourceDefinition(def *models.ResourceDefinition) error {
	if def.Name == "" {
		return errors.New("resource name cannot be empty")
	} else if def.PathName == "" {
		return errors.New("resource path_name cannot be empty")
	} else if len(def.Fields) == 0 {
		return errors.New("resource fields cannot be empty")
	}

	return nil
}

func processFields(fields map[string]string) ([]*bson.Element, error) {
	fieldElements := make([]*bson.Element, 0)
	for key, value := range fields {
		if !supportedType(value) {
			return nil, fmt.Errorf("'%s' is not a supported field type", value)
		}
		fieldElements = append(fieldElements, bson.EC.String(key, value))
	}

	return fieldElements, nil
}

// AddResourceDefinition creates a new resource definition in the users' collection
func AddResourceDefinition(c *gin.Context) {
	// Marshal JSON into ResourceDefinition
	var resourceDefinition models.ResourceDefinition
	c.BindJSON(&resourceDefinition)

	// Validate the definition
	if err := validateResourceDefinition(&resourceDefinition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process the resource fields into bson
	fieldElements, err := processFields(resourceDefinition.Fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create document
	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.String("name", resourceDefinition.Name))
	resourceElements = append(resourceElements, bson.EC.String("path_name", resourceDefinition.PathName))
	resourceElements = append(resourceElements, bson.EC.SubDocumentFromElements("fields", fieldElements...))

	// Get a connection and insert the new document
	collection := database.Connect().Collection(database.ResourceDefinitions)
	result, err := collection.InsertOne(
		context.Background(),
		bson.NewDocument(resourceElements...),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the inserted ID for the response
	resourceDefinition.ID = result.InsertedID.(objectid.ObjectID).Hex()
	c.JSON(http.StatusCreated, resourceDefinition)
}

// ListResourceDefinitions returns the list of all resource definitions
func ListResourceDefinitions(c *gin.Context) {
	definitions := make([]models.ResourceDefinition, 0)

	collection := database.Connect().Collection(database.ResourceDefinitions)

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	doc := bson.NewDocument()
	for cursor.Next(context.Background()) {
		doc.Reset()
		err := cursor.Decode(doc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		def := models.ResourceDefinition{
			Fields: make(map[string]string),
		}
		def.ID = doc.Lookup("_id").ObjectID().Hex()
		def.Name = doc.Lookup("name").StringValue()
		def.PathName = doc.Lookup("path_name").StringValue()
		fields := doc.Lookup("fields").MutableDocument()
		fieldKeys, _ := fields.Keys(false)
		for _, key := range fieldKeys {
			def.Fields[key.String()] = fields.Lookup(key.String()).StringValue()
		}

		definitions = append(definitions, def)
	}
	c.JSON(http.StatusCreated, gin.H{"resources": definitions})
}
