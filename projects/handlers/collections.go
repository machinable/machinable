package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/nsjostrom/machinable/projects/database"
	"bitbucket.org/nsjostrom/machinable/projects/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type errorItem struct {
	ID objectid.ObjectID `bson:"_id"`
}

func getOrCreateCollection(name, projectSlug string) error {
	collection := database.Connect().Collection(database.CollectionNames(projectSlug))

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
		err := createNewCollection(name, projectSlug)
		return err
	}

	return err
}

func createNewCollection(name, projectSlug string) error {
	// Create document
	resourceElements := make([]*bson.Element, 0)
	resourceElements = append(resourceElements, bson.EC.String("name", name))
	resourceElements = append(resourceElements, bson.EC.Time("created", time.Now()))

	// Get a connection and insert the new document
	collection := database.Connect().Collection(database.CollectionNames(projectSlug))
	_, err := collection.InsertOne(
		context.Background(),
		bson.NewDocument(resourceElements...),
	)

	return err
}

// AddCollection creates a new collection
func AddCollection(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	var newCollection models.Collection
	c.BindJSON(&newCollection)

	if newCollection.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	err := createNewCollection(newCollection.Name, projectSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// GetCollections returns the list of collections for a user
func GetCollections(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	collections := make([]*models.Collection, 0)

	collection := database.Connect().Collection(database.CollectionNames(projectSlug))

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
		// get count
		docName := doc.Lookup("name").StringValue()
		collectPath := database.CollectionDocs(projectSlug, docName)
		collection := database.Connect().Collection(collectPath)
		cnt, _ := collection.CountDocuments(nil, nil, nil)

		collections = append(collections,
			&models.Collection{
				Name:    doc.Lookup("name").StringValue(),
				ID:      doc.Lookup("_id").ObjectID().Hex(),
				Created: doc.Lookup("created").Time(),
				Items:   cnt,
			})
	}
	c.JSON(http.StatusOK, gin.H{"items": collections})
}

// DeleteCollection deletes a project collection along with all of it's data
func DeleteCollection(c *gin.Context) {
	collectionID := c.Param("collectionName")
	projectSlug := c.MustGet("project").(string)

	collections := database.Collection(database.CollectionNames(projectSlug))
	// Get the object id for collection name
	objectID, err := objectid.FromHex(collectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find object based on ID
	documentResult := collections.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		nil,
	)

	if documentResult == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "collection not found"})
	}

	// decode result into document
	doc := bson.NewDocument()
	err = documentResult.Decode(doc)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("collection not found, '%s'", collectionID)})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// drop the collection docs (actual data for this collection)
	collectionData := database.Collection(database.CollectionDocs(projectSlug, collectionName))

	err = collectionData.Drop(nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// PutObjectInCollection alters an existing collection document
func PutObjectInCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	objectIDStr := c.Param("objectID")
	projectSlug := c.MustGet("project").(string)

	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(objectIDStr)
	if err != nil && objectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid identifier '%s': %s", objectIDStr, err.Error())})
		return
	}

	bdoc := make(map[string]interface{})

	err = c.BindJSON(&bdoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// iterate over root keys for reserved fields
	els := make([]*bson.Element, 0)
	for key := range bdoc {
		if reservedField(key) {
			// remove field for PUT
			delete(bdoc, key)
		} else {
			// append to element slice
			els = append(els, bson.EC.Interface(key, bdoc[key]))
		}
	}

	// Get a connection and update the document
	collection := database.Connect().Collection(database.CollectionDocs(projectSlug, collectionName))
	_, err = collection.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				els...,
			),
		),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bdoc)
}

// AddObjectToCollection adds a new document to the collection
func AddObjectToCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	projectSlug := c.MustGet("project").(string)
	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	if err := getOrCreateCollection(collectionName, projectSlug); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get collection"})
		return
	}

	bdoc := make(map[string]interface{})

	err := c.BindJSON(&bdoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// iterate over root keys for reserved fields
	for key := range bdoc {
		if reservedField(key) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("'%s' is a reserved field", key)})
			return
		}
	}

	// Get a connection and insert the new document
	collection := database.Connect().Collection(database.CollectionDocs(projectSlug, collectionName))
	result, err := collection.InsertOne(
		context.Background(),
		bdoc,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	insertedID, ok := result.InsertedID.(objectid.ObjectID)
	if ok {
		bdoc["id"] = insertedID.Hex()
	} else {
		bdoc["id"] = result.InsertedID.(interface{})
	}

	c.JSON(http.StatusCreated, bdoc)
}

// GetObjectsFromCollection returns the full list of documents
func GetObjectsFromCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	projectSlug := c.MustGet("project").(string)
	if collectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection name cannot be empty"})
		return
	}

	collection := database.Connect().Collection(database.CollectionDocs(projectSlug, collectionName))

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
	c.IndentedJSON(http.StatusOK, gin.H{"items": documents})
}

// GetObjectFromCollection returns a single object with the ID for this resource
func GetObjectFromCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	objectIDStr := c.Param("objectID")
	projectSlug := c.MustGet("project").(string)
	collection := database.Collection(database.CollectionDocs(projectSlug, collectionName))

	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(objectIDStr)
	if err != nil && objectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid identifier '%s': %s", objectIDStr, err.Error())})
		return
	}

	// Find object based on ID
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		nil,
	)

	if documentResult == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no documents for collection"})
	}

	// Decode result into document
	doc := bson.NewDocument()
	err = documentResult.Decode(doc)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("object does not exist, '%s'", objectIDStr)})
		return
	}

	// Lookup  definitions for this resource
	object, err := parseUnknownDocumentToMap(doc, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, object)
}

// DeleteObjectFromCollection deletes the object from the collection
func DeleteObjectFromCollection(c *gin.Context) {
	collectionName := c.Param("collectionName")
	objectIDStr := c.Param("objectID")
	projectSlug := c.MustGet("project").(string)
	collection := database.Collection(database.CollectionDocs(projectSlug, collectionName))

	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(objectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid identifier '%s': %s", objectIDStr, err.Error())})
		return
	}

	// Delete the object
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

	c.JSON(http.StatusNoContent, gin.H{})
}
