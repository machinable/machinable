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

func getOrCreateCollection(name string) error {
	collection := database.Collection(database.Collections)

	// Find the resource definition for this object
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("name", name),
		),
		nil,
	)

	// Decode result into document
	doc := bson.Document{}
	err := documentResult.Decode(&doc)
	if err != nil {
		err := createNewCollection(name)
		return err
	}

	return err
}

func createNewCollection(name string) error {
	// Create document
	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.String("name", name))

	// Get a connection and insert the new document
	collection := database.Connect().Collection(database.Collections)
	_, err := collection.InsertOne(
		context.Background(),
		bson.NewDocument(resourceElements...),
	)

	return err
}

// AddCollection creates a new collection
func AddCollection(c *gin.Context) {
	var newCollection models.Collection
	c.BindJSON(&newCollection)

	if newCollection.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	err := createNewCollection(newCollection.Name)
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

// AddObjectToCollection adds a new document to the collection
func AddObjectToCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	if err := getOrCreateCollection(collectionName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get collection"})
		return
	}

	bdoc := make(map[string]interface{})

	c.BindJSON(&bdoc)

	// Get a connection and insert the new document
	collection := database.Connect().Collection(fmt.Sprintf(database.CollectionFormat, collectionName))
	result, err := collection.InsertOne(
		context.Background(),
		bdoc,
	)

	// TODO
	// Load result id and try to decode. If an error occurs, delete the document and return the error message with a 400 status code

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bdoc["id"] = result.InsertedID.(objectid.ObjectID).Hex()
	c.JSON(http.StatusCreated, bdoc)
}

// GetObjectsFromCollection returns the full list of documents
func GetObjectsFromCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}
	documents := make([]map[string]interface{}, 0)

	collection := database.Connect().Collection(fmt.Sprintf(database.CollectionFormat, collectionName))

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//doc := bson.NewDocument()
	for cursor.Next(context.Background()) {
		//doc.Reset()
		doc := make(map[string]interface{})
		err := cursor.Decode(doc)
		if err != nil {
			var errid string
			item := &errorItem{}
			anotherErr := cursor.Decode(item)
			if anotherErr == nil {
				errid = item.ID.Hex()
			}
			documents = append(documents, map[string]interface{}{
				"id":    errid,
				"error": err.Error(),
			})
		} else {
			// get stringified version of the ID
			objectID, ok := doc["_id"].(objectid.ObjectID)
			if ok {
				doc["id"] = objectID.Hex()
				delete(doc, "_id")
			}

			documents = append(documents, doc)
		}
	}
	c.JSON(http.StatusOK, gin.H{"items": documents})
}
