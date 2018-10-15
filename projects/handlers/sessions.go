package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"bitbucket.org/nsjostrom/machinable/auth"
	"bitbucket.org/nsjostrom/machinable/projects/database"
	"bitbucket.org/nsjostrom/machinable/projects/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Create creates a new project user session
func Create(c *gin.Context) {
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
	userPassword := pair[1]

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
	err := documentResult.Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	// compare passwords
	if !auth.CompareHashAndPassword(user.PasswordHash, userPassword) {
		c.JSON(http.StatusNotFound, gin.H{})
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
		},
	}

	accessToken, err := auth.CreateAccessToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create the access token"})
		return
	}
	log.Println(accessToken)

	c.JSON(http.StatusCreated, gin.H{})
}

// List lists all active user sessions for a project
func List() {

}

// Revoke deletes a session from the project collection
func Revoke() {

}
