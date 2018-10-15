package handlers

import (
	"context"
	"net/http"
	"time"

	"bitbucket.org/nsjostrom/machinable/projects/database"
	"bitbucket.org/nsjostrom/machinable/projects/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// AddToken creates a new api token for this project
func AddToken(c *gin.Context) {
	var newToken models.NewProjectToken
	projectSlug := c.MustGet("project").(string)

	c.BindJSON(&newToken)

	token := &models.ProjectAPIToken{
		ID:          objectid.New(), // I don't like this
		Created:     time.Now(),
		TokenHash:   newToken.Token, // salt and hash
		Description: newToken.Description,
		Read:        newToken.Read,
		Write:       newToken.Write,
	}

	// Get the resources.{resourcePathName} collection
	rc := database.Collection(database.TokenDocs(projectSlug))
	_, err := rc.InsertOne(
		context.Background(),
		token,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, token)
}

// ListTokens lists all api tokens of this project
func ListTokens(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	tokens := make([]*models.ProjectAPIToken, 0)

	collection := database.Connect().Collection(database.TokenDocs(projectSlug))

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for cursor.Next(context.Background()) {
		var token models.ProjectAPIToken
		err := cursor.Decode(&token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tokens = append(tokens, &token)
	}

	c.JSON(http.StatusOK, gin.H{"items": tokens})
}

// GetToken retrieves a single api token of this project by ID
func GetToken(c *gin.Context) {
	//projectSlug := c.MustGet("project").(string)
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// DeleteToken removes an api token by ID
func DeleteToken(c *gin.Context) {
	//projectSlug := c.MustGet("project").(string)
	c.JSON(http.StatusNotImplemented, gin.H{})
}
