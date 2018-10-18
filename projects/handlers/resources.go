package handlers

import (
	"context"
	"net/http"

	"bitbucket.org/nsjostrom/machinable/projects/database"
	"bitbucket.org/nsjostrom/machinable/projects/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// AddResourceDefinition creates a new resource definition in the users' collection
func AddResourceDefinition(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	// Marshal JSON into ResourceDefinition
	var resourceDefinition models.ResourceDefinition
	c.BindJSON(&resourceDefinition)

	// Validate the definition
	if err := validateResourceDefinition(&resourceDefinition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resDefCollection := database.Collection(database.ResourceDefinitions(projectSlug))
	if exists := definitionExists(&resourceDefinition, resDefCollection); exists == true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource already exists"})
		return
	}

	// Process the resource fields into bson
	propertyElements, err := processProperties(resourceDefinition.Properties, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create document
	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.String("title", resourceDefinition.Title))
	resourceElements = append(resourceElements, bson.EC.String("path_name", resourceDefinition.PathName))
	resourceElements = append(resourceElements, bson.EC.SubDocumentFromElements("properties", propertyElements...))

	// Get a connection and insert the new document
	collection := database.Connect().Collection(database.ResourceDefinitions(projectSlug))
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
	projectSlug := c.MustGet("project").(string)
	definitions := make([]*models.ResourceDefinition, 0)

	collection := database.Connect().Collection(database.ResourceDefinitions(projectSlug))

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
	projectSlug := c.MustGet("project").(string)
	resDefCollection := database.Collection(database.ResourceDefinitions(projectSlug))
	def, err := getDefinitionByID(resourceID, resDefCollection)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, def)
}

// DeleteResourceDefinition deletes the definition and drops the resource collection
func DeleteResourceDefinition(c *gin.Context) {
	resourceID := c.Param("resourceDefinitionID")
	projectSlug := c.MustGet("project").(string)

	resDefCollection := database.Collection(database.ResourceDefinitions(projectSlug))
	// Get definition for resource name
	def, err := getDefinitionByID(resourceID, resDefCollection)
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

	// Delete the definition
	_, err = resDefCollection.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resourceCollection := database.Collection(database.ResourceDocs(projectSlug, resourcePathName))
	resourceCollection.Drop(nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
