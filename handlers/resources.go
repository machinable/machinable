package handlers

import (
	"context"
	"fmt"
	"net/http"

	"bitbucket.org/nsjostrom/machinable/database"
	"bitbucket.org/nsjostrom/machinable/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

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
	definitions := make([]*models.ResourceDefinition, 0)

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
		def, _ := parseDefinition(doc)
		definitions = append(definitions, def)
	}
	c.JSON(http.StatusOK, gin.H{"items": definitions})
}

// GetResourceDefinition returns a single resource definition
func GetResourceDefinition(c *gin.Context) {
	resourceID := c.Param("resourceDefinitionID")
	def, err := getDefinitionByID(resourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, def)
}

// DeleteResourceDefinition deletes the definition and drops the resource collection
func DeleteResourceDefinition(c *gin.Context) {
	resourceID := c.Param("resourceDefinitionID")

	// Get definition for resource name
	def, err := getDefinitionByID(resourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	resourcePathName := def.PathName

	// Get the object id
	objectID, err := objectid.FromHex(resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	collection := database.Collection(database.ResourceDefinitions)

	// Delete the definition
	_, err = collection.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resourceCollection := database.Collection(fmt.Sprintf(database.ResourceFormat, resourcePathName))
	resourceCollection.Drop(nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
