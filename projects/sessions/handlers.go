package sessions

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mssola/user_agent"

	"github.com/anothrnick/machinable/auth"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	as "github.com/anothrnick/machinable/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// New returns a pointer to a new `Users` struct
func New(db interfaces.Datastore) *Sessions {
	return &Sessions{
		store: db,
	}
}

// Sessions wraps the datastore and any HTTP handlers for project user sessions
type Sessions struct {
	store interfaces.Datastore
}

func (s *Sessions) generateSession(userID, ip, userAgent string) *models.Session {
	location, _ := as.GetGeoIP(ip)

	ua := user_agent.New(userAgent)

	bname, bversion := ua.Browser()
	session := &models.Session{
		ID:           objectid.New(), // no
		UserID:       userID,
		Location:     location,
		Mobile:       ua.Mobile(),
		IP:           ip,
		LastAccessed: time.Now(),
		Browser:      bname + " " + bversion,
		OS:           ua.OS(),
	}

	return session
}

// CreateSession creates a new project user session
func (s *Sessions) CreateSession(c *gin.Context) {
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

	user, err := s.store.GetUserByUsername(projectSlug, userName)
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
		"projects": map[string]interface{}{
			projectSlug: true,
		},
		"user": map[string]interface{}{
			"id":     user.ID.Hex(),
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
	session := s.generateSession(user.ID.Hex(), c.ClientIP(), c.Request.UserAgent())
	err = s.store.CreateSession(projectSlug, session)
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
func (s *Sessions) ListSessions(c *gin.Context) {
	projectSlug := c.MustGet("project").(string)

	sessions, err := s.store.ListSessions(projectSlug)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": sessions})
}

// RevokeSession deletes a session from the project collection
func (s *Sessions) RevokeSession(c *gin.Context) {
	sessionID := c.Param("sessionID")
	projectSlug := c.MustGet("project").(string)

	err := s.store.DeleteSession(projectSlug, sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// RefreshSession uses the refresh token to generate a new access token
func (s *Sessions) RefreshSession(c *gin.Context) {

}
