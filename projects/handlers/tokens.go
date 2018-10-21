package handlers

import (
	"context"
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

// AddToken creates a new api token for this project
func AddToken(c *gin.Context) {
	var newToken models.NewProjectToken
	projectSlug := c.MustGet("project").(string)

	c.BindJSON(&newToken)

	err := newToken.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate hashed token
	tokenHash, err := auth.HashPassword(newToken.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newToken.Token = ""

	token := &models.ProjectAPIToken{
		ID:          objectid.New(), // I don't like this
		Created:     time.Now(),
		TokenHash:   tokenHash,
		Description: newToken.Description,
		Read:        newToken.Read,
		Write:       newToken.Write,
	}

	// get the tokens collection
	rc := database.Collection(database.TokenDocs(projectSlug))
	// save token
	_, err = rc.InsertOne(
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

// GenerateToken retrieves a single api token of this project by ID
func GenerateToken(c *gin.Context) {
	UUID, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": UUID.String()})
}

// DeleteToken removes an api token by ID
func DeleteToken(c *gin.Context) {
	//projectSlug := c.MustGet("project").(string)
	c.JSON(http.StatusNotImplemented, gin.H{})
}
