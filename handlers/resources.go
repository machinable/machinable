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
	var resourceDefinition models.ResourceDefinition

	c.BindJSON(&resourceDefinition)

	//db := database.Connect()

	if err := validateResourceDefinition(&resourceDefinition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fieldElements, err := processFields(resourceDefinition.Fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.String("name", resourceDefinition.Name))
	resourceElements = append(resourceElements, bson.EC.String("path_name", resourceDefinition.PathName))
	resourceElements = append(resourceElements, bson.EC.SubDocumentFromElements("fields", fieldElements...))

	collection := database.Connect().Collection(database.ResourceDefinitions)

	result, err := collection.InsertOne(
		context.Background(),
		bson.NewDocument(resourceElements...),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resourceDefinition.ID = result.InsertedID.(objectid.ObjectID).Hex()
	c.JSON(http.StatusCreated, resourceDefinition)
}
