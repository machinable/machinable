package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"bitbucket.org/nsjostrom/machinable/auth"
	"bitbucket.org/nsjostrom/machinable/projects/database"
	"bitbucket.org/nsjostrom/machinable/projects/models"
	as "bitbucket.org/nsjostrom/machinable/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

// CreateSession creates a new project user session
func CreateSession(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	// basic auth for login
	authorizationHeader, _ := c.Request.Header["Authorization"]
	if len(authorizationHeader) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no authorization header"})
		return
	}
	authzHeader := strings.SplitN(authorizationHeader[0], " ", 2)

	if len(authzHeader) != 2 || authzHeader[0] != "Basic" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "malformed authorization header"})
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(authzHeader[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no authorization header payload"})
		return
	}

	userName := pair[0]
	userPassword := strings.Trim(pair[1], "\n")

	if userName == "" {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	// get the users collection
	userCollection := database.Collection(database.UserDocs(projectSlug))

	// look up the user
	documentResult := userCollection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("username", userName),
		),
		nil,
	)

	user := &models.ProjectUser{}
	// decode user document
	err := documentResult.Decode(user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "username not found"})
		return
	}

	// compare passwords
	if !auth.CompareHashAndPassword(user.PasswordHash, userPassword) {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid password"})
		return
	}

	// create access token
	claims := jwt.MapClaims{
		"projects": []string{
			projectSlug,
		},
		"user": map[string]interface{}{
			"id":     user.ID,
			"name":   user.Username,
			"active": true,
			"read":   user.Read,
			"write":  user.Write,
			"type":   "project",
		},
	}

	accessToken, err := auth.CreateAccessToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create the access token"})
		return
	}

	// create session in database (refresh token)
	sessionCollection := database.Collection(database.SessionDocs(projectSlug))
	session, err := as.CreateSession(user.ID.Hex(), c.ClientIP(), c.Request.UserAgent(), sessionCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	refreshToken, err := auth.CreateRefreshToken(session.ID.Hex(), user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Successfully logged in",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"session_id":    session.ID.Hex(),
	})
}

// ListSessions lists all active user sessions for a project
func ListSessions(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)
	sessions := make([]*as.Session, 0)

	collection := database.Connect().Collection(database.SessionDocs(projectSlug))

	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for cursor.Next(context.Background()) {
		var session as.Session
		err := cursor.Decode(&session)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sessions = append(sessions, &session)
	}

	c.JSON(http.StatusOK, gin.H{"items": sessions})
}

// RevokeSession deletes a session from the project collection
func RevokeSession(c *gin.Context) {

}

// RefreshSession uses the refresh token to generate a new access token
func RefreshSession(c *gin.Context) {

}
