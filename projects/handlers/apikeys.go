package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/nsjostrom/machinable/auth"
	"bitbucket.org/nsjostrom/machinable/projects/database"
	"bitbucket.org/nsjostrom/machinable/projects/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	uuid "github.com/satori/go.uuid"
)

// AddKey creates a new api key for this project
func AddKey(c *gin.Context) {
	var newKey models.NewProjectKey
	projectSlug := c.MustGet("project").(string)

	c.BindJSON(&newKey)

	err := newKey.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate sha1 key
	keyHash := auth.SHA1(newKey.Key)
	newKey.Key = ""

	key := &models.ProjectAPIKey{
		ID:          objectid.New(), // I don't like this
		Created:     time.Now(),
		KeyHash:     keyHash,
		Description: newKey.Description,
		Read:        newKey.Read,
		Write:       newKey.Write,
	}

	// get the keys collection
	rc := database.Collection(database.KeyDocs(projectSlug))
	// save key
	_, err = rc.InsertOne(
		context.Background(),
		key,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, key)
}

// ListKeys lists all api tokens of this project
func ListKeys(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	tokens := make([]*models.ProjectAPIKey, 0)

	collection := database.Connect().Collection(database.KeyDocs(projectSlug))

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for cursor.Next(context.Background()) {
		var token models.ProjectAPIKey
		err := cursor.Decode(&token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tokens = append(tokens, &token)
	}

	c.JSON(http.StatusOK, gin.H{"items": tokens})
}

// GenerateKey retrieves a single api token of this project by ID
func GenerateKey(c *gin.Context) {
	UUID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": UUID.String()})
}

// DeleteKey removes an api token by ID
func DeleteKey(c *gin.Context) {
	keyID := c.Param("keyID")
	projectSlug := c.MustGet("project").(string)

	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(keyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid identifier '%s': %s", keyID, err.Error())})
		return
	}

	collection := database.Connect().Collection(database.KeyDocs(projectSlug))
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
