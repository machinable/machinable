package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/dgrijalva/jwt-go"
	"github.com/mssola/user_agent"

	"bitbucket.org/nsjostrom/machinable/auth"
	"bitbucket.org/nsjostrom/machinable/projects/database"
	"bitbucket.org/nsjostrom/machinable/projects/models"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
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
		},
	}

	accessToken, err := auth.CreateAccessToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create the access token"})
		return
	}

	// create session in database (refresh token)
	sessionCollection := database.Collection(database.SessionDocs(projectSlug))
	session, err := createSession(user.ID.Hex(), c.ClientIP(), c.Request.UserAgent(), sessionCollection)
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
	sessions := make([]*models.ProjectSession, 0)

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
		var session models.ProjectSession
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

func getGeoIP(ip string) (string, error) {
	// ... this should be changed to get the access key from a config or environment variable
	accessKey := "85a38b87f3b696c7dcbf8f6f58c3c6a9"
	url := fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, accessKey)

	ipStackData := struct {
		City       string `json:"city"`
		RegionCode string `json:"region_code"`
	}{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New("error creating request")
	}

	// set client with 10 second timeout
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("error making request")
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&ipStackData); err != nil {
		return "", errors.New("error decoding response")
	}

	location := ""

	if ipStackData.City != "" && ipStackData.RegionCode != "" {
		location = ipStackData.City + ", " + ipStackData.RegionCode
	}

	return location, nil
}

func createSession(userID, ip, userAgent string, collection *mongo.Collection) (*models.ProjectSession, error) {
	location, _ := getGeoIP(ip)

	ua := user_agent.New(userAgent)

	bname, bversion := ua.Browser()
	session := &models.ProjectSession{
		ID:           objectid.New(),
		UserID:       userID,
		Location:     location,
		Mobile:       ua.Mobile(),
		IP:           ip,
		LastAccessed: time.Now(),
		Browser:      bname + " " + bversion,
		OS:           ua.OS(),
	}

	// save the user
	_, err := collection.InsertOne(
		context.Background(),
		session,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}
