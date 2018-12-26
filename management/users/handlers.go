package users

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/nsjostrom/machinable/auth"
	"bitbucket.org/nsjostrom/machinable/dsi/interfaces"
	"bitbucket.org/nsjostrom/machinable/dsi/models"
	"bitbucket.org/nsjostrom/machinable/management/database"
	as "bitbucket.org/nsjostrom/machinable/sessions"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// New returns a pointer to a new `Users`
func New(db interfaces.Datastore) *Users {
	return &Users{
		store: db,
	}
}

// Users contains the datastore and any HTTP handlers needed for application users
type Users struct {
	store interfaces.Datastore
}

type newUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate checks that the new user has a username and password.
func (u *newUserBody) Validate() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("invalid username or password")
	}
	return nil
}

func createAccessToken(user *models.User) (string, error) {
	// TODO: add user project map to jwt
	claims := jwt.MapClaims{
		"projects": make(map[string]interface{}),
		"user": map[string]interface{}{
			"id":     user.ID.Hex(),
			"name":   user.Username,
			"type":   "app",
			"active": true,
		},
	}

	accessToken, err := auth.CreateAccessToken(claims)
	if err != nil {
		return "", errors.New("failed to create the access token")
	}

	return accessToken, nil
}

// createTokensAndSession returns an accessToken, refreshToken, error
func createTokensAndSession(user *models.User, c *gin.Context) (string, string, *models.Session, error) {
	// create access token
	accessToken, err := createAccessToken(user)
	if err != nil {
		return "", "", nil, err
	}

	// create session in database (refresh token)
	sessionCollection := database.Connect().Collection(database.Sessions)
	session, err := as.CreateSession(user.ID.Hex(), c.ClientIP(), c.Request.UserAgent(), sessionCollection)
	if err != nil {
		return "", "", nil, errors.New("failed to create session")
	}

	refreshToken, err := auth.CreateRefreshToken(session.ID.Hex(), user.ID.Hex())
	if err != nil {
		return "", "", nil, errors.New("failed to create refresh token")
	}

	return accessToken, refreshToken, session, nil
}

// RegisterUser creates a new valid user in the database. The user will receive an access and refresh
// token on register. The user can then login next time.
func (u *Users) RegisterUser(c *gin.Context) {
	var newUser newUserBody

	c.BindJSON(&newUser)

	// validate user
	err := newUser.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate hashed password
	passwordHash, err := auth.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newUser.Password = ""

	// check for duplicate username
	existingUser, _ := u.store.GetAppUserByUsername(newUser.Username)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	// create project user object
	user := &models.User{
		ID:           objectid.New(), // I don't like this
		Created:      time.Now(),
		PasswordHash: passwordHash,
		Username:     newUser.Username,
	}

	// save the user
	err = u.store.CreateAppUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: refactor sessions
	accessToken, refreshToken, session, err := createTokensAndSession(user, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Successfully registered",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"session_id":    session.ID.Hex(),
	})
}

// LoginUser creates a session for an existing management application user.
func (u *Users) LoginUser(c *gin.Context) {
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

	// look up the user
	user, err := u.store.GetAppUserByUsername(userName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "username not found"})
		return
	}

	// compare passwords
	if !auth.CompareHashAndPassword(user.PasswordHash, userPassword) {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid password"})
		return
	}

	accessToken, refreshToken, session, err := createTokensAndSession(user, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"session_id":    session.ID.Hex(),
	})
}

// RefreshToken exchanges a valid refresh token for a new access token.
func (u *Users) RefreshToken(c *gin.Context) {
	// get session and user id from context, should have been injected by ValidateRefreshToken
	sessionID, ok := c.MustGet("session_id").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no session"})
		return
	}

	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user"})
		return
	}

	var session models.Session
	var user models.User
	// load session to update last accessed

	// verify session exists
	sessionCollection := database.Connect().Collection(database.Sessions)
	sessionObjectID, _ := objectid.FromHex(sessionID)
	documentResult := sessionCollection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", sessionObjectID),
		),
		nil,
	)

	// decode session document
	err := documentResult.Decode(&session)
	if err != nil {
		// no documents in result, user does not exist
		c.JSON(http.StatusNotFound, gin.H{"message": "error creating access token."})
		return
	}

	// verify user exists
	_, err = u.store.GetAppUserByID(userID)
	if err != nil {
		// no documents in result, user does not exist
		c.JSON(http.StatusNotFound, gin.H{"message": "error creating access token."})
		return
	}

	// create new access jwt
	accessToken, err := createAccessToken(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "error creating access token."})
		return
	}

	// update session `last_accessed` time
	_, err = sessionCollection.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.ObjectID("_id", sessionObjectID),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Time("last_accessed", time.Now()),
			),
		),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

// RevokeSession deletes a user's session.
func (u *Users) RevokeSession(c *gin.Context) {
	sessionID := c.Param("sessionID")

	sessionCollection := database.Connect().Collection(database.Sessions)
	// Get the object id
	objectID, err := objectid.FromHex(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the definition
	_, err = sessionCollection.DeleteOne(
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

// ResetPassword authenticates the user using the old password, then sets a new password for the application user.
func (u *Users) ResetPassword(c *gin.Context) {

}
