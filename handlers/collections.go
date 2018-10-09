package handlers

import (
	"context"
	"net/http"

	"bitbucket.org/nsjostrom/machinable/database"
	"bitbucket.org/nsjostrom/machinable/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

// AddCollection creates a new collection
func AddCollection(c *gin.Context) {
	var newCollection models.Collection
	c.BindJSON(&newCollection)

	if newCollection.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	// Create document
	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.String("name", newCollection.Name))

	// Get a connection and insert the new document
	collection := database.Connect().Collection(database.Collections)
	_, err := collection.InsertOne(
		context.Background(),
		bson.NewDocument(resourceElements...),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// GetCollections returns the list of collections for a user
func GetCollections(c *gin.Context) {
	collections := make([]*models.Collection, 0)

	collection := database.Connect().Collection(database.Collections)

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
		collections = append(collections, &models.Collection{Name: doc.Lookup("name").StringValue()})
	}
	c.JSON(http.StatusOK, gin.H{"items": collections})
}
