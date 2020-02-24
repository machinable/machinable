package users

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"

	"github.com/anothrnick/machinable/auth"
	"github.com/anothrnick/machinable/config"
	"github.com/anothrnick/machinable/dsi/interfaces"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/anothrnick/machinable/events"
	as "github.com/anothrnick/machinable/sessions"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// New returns a pointer to a new `Users`
func New(db interfaces.Datastore, cache redis.UniversalClient, config *config.AppConfig) *Users {
	jwt := auth.NewJWT(config)
	return &Users{
		store:  db,
		cache:  cache,
		config: config,
		jwt:    jwt,
	}
}

// Users contains the datastore and any HTTP handlers needed for application users
type Users struct {
	store  interfaces.Datastore
	cache  redis.UniversalClient
	config *config.AppConfig
	jwt    *auth.JWT
}

func (u *Users) createAccessToken(user *models.User) (string, error) {
	// TODO: add user project map to jwt
	claims := jwt.MapClaims{
		"projects": make(map[string]interface{}),
		"user": map[string]interface{}{
			"id":     user.ID,
			"name":   user.Username,
			"type":   "app",
			"active": true,
		},
	}

	accessToken, err := u.jwt.CreateAccessToken(claims)
	if err != nil {
		return "", errors.New("failed to create the access token")
	}

	return accessToken, nil
}

// createTokensAndSession returns an accessToken, refreshToken, error
func (u *Users) createTokensAndSession(user *models.User, c *gin.Context) (string, string, *models.Session, error) {
	// create access token
	accessToken, err := u.createAccessToken(user)
	if err != nil {
		return "", "", nil, err
	}

	// create session in database (refresh token)
	session := as.CreateSession(user.ID, c.ClientIP(), c.Request.UserAgent(), u.config)
	err = u.store.CreateAppSession(session)
	if err != nil {
		return "", "", nil, errors.New("failed to create session")
	}

	refreshToken, err := u.jwt.CreateRefreshToken(session.ID, user.ID)
	if err != nil {
		return "", "", nil, errors.New("failed to create refresh token")
	}

	return accessToken, refreshToken, session, nil
}

// verifyEmail creates the verification code, saves it to redis with a timeout, and queues the email
func (u *Users) verifyEmail(user *models.User) error {
	// create verification code (UUID)
	verificationCode := uuid.NewV4()

	// create redis key with 5 minute timeout
	vKey := fmt.Sprintf("verificationCode:%s", verificationCode)
	if err := u.cache.Set(vKey, user.ID, time.Minute*5).Err(); err != nil {
		log.Println(err)
		return errors.New("failed creating verification code")
	}

	// queue verification email
	vURL := fmt.Sprintf("%s/%s/%s", u.config.AppHost, "verify", verificationCode)
	plainText := fmt.Sprintf("Verify your Machinable email address by following this link: %s", vURL)
	n := &events.Notification{
		Template:         "default",
		Subject:          "Email Verification",
		ReceiverName:     user.Username,
		ReceiverEmail:    user.Email,
		PlainTextContent: plainText,
		Data: map[string]string{
			"Name":    user.Username,
			"Site":    "Machinable",
			"URL":     vURL,
			"Company": "Machinable",
		},
		Meta: map[string]string{
			"user_id": user.ID,
		},
	}

	b, err := json.Marshal(n)
	if err != nil {
		log.Println(err)
		return errors.New("failed creating verification code")
	}
	if err := u.cache.RPush(events.QueueEmailNotifications, b).Err(); err != nil {
		log.Println(err)
		return errors.New("failed creating verification code")
	}

	return nil
}

// RegisterUser creates a new valid user in the database. The user will receive an access and refresh
// token on register. The user can then login next time.
func (u *Users) RegisterUser(c *gin.Context) {
	var newUser newUserBody

	c.BindJSON(&newUser)

	// validate user
	err := newUser.Validate(u.config.ReCaptchaSecret)
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

	// check for duplicate email
	if _, err := u.store.GetAppUserByUsername(newUser.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is already in use"})
		return
	}

	// create project user object
	user := &models.User{
		Created:      time.Now(),
		PasswordHash: passwordHash,
		Email:        newUser.Email,
		Username:     newUser.Email,
	}

	// save the user
	err = u.store.CreateAppUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return success
	if err := u.verifyEmail(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success! Check your email to verify your account",
	})
}

// EmailVerification looks up the verification code in redis and activates the associated user
func (u *Users) EmailVerification(c *gin.Context) {
	// get verification code
	verificationCode := c.Param("verificationCode")

	// lookup key in redis
	vKey := fmt.Sprintf("verificationCode:%s", verificationCode)

	// get the request count key for the current window
	userID, rerr := u.cache.Get(vKey).Result()
	if rerr == redis.Nil {
		// key doesn't exist
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to verify email for user"})
		return
	} else if rerr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to verify email for user"})
		return
	}

	// find user
	user, err := u.store.GetAppUserByID(userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user does not exist"})
		return
	}

	// update user in db, set active
	if err := u.store.ActivateUser(userID, true); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to activate user"})
		return
	}

	// create session and return api tokens
	// TODO: refactor sessions
	accessToken, refreshToken, session, err := u.createTokensAndSession(user, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// delete key from redis
	u.cache.Del(vKey)

	// return session tokens
	c.JSON(http.StatusOK, gin.H{
		"message":       "Successfully activated user",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"session_id":    session.ID,
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
		c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	}

	// compare passwords
	if !auth.CompareHashAndPassword(user.PasswordHash, userPassword) {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid password"})
		return
	}

	if !user.Active {
		c.JSON(http.StatusNotFound, gin.H{"error": "user account is not active, check email for verification"})
		u.verifyEmail(user)
		return
	}

	accessToken, refreshToken, session, err := u.createTokensAndSession(user, c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"session_id":    session.ID,
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

	// load session to update last accessed

	// verify session exists
	_, err := u.store.GetAppSession(sessionID)

	if err != nil {
		log.Println(err)
		// no documents in result, user does not exist
		c.JSON(http.StatusNotFound, gin.H{"message": "error creating access token."})
		return
	}

	// verify user exists
	user, err := u.store.GetAppUserByID(userID)
	if err != nil {
		log.Println(err)
		// no documents in result, user does not exist
		c.JSON(http.StatusNotFound, gin.H{"message": "error creating access token."})
		return
	}

	// create new access jwt
	accessToken, err := u.createAccessToken(user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"message": "error creating access token."})
		return
	}

	// update session `last_accessed` time
	err = u.store.UpdateAppSessionLastAccessed(sessionID, time.Now())

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

	err := u.store.DeleteAppSession(sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// ResetPassword authenticates the user using the old password, then sets a new password for the application user.
func (u *Users) ResetPassword(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user"})
		return
	}

	var passwordUpdate updatePasswordBody

	c.BindJSON(&passwordUpdate)

	// validate user
	err := passwordUpdate.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate hashed password
	// passwordHash, err := auth.HashPassword(passwordUpdate.OldPW)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// passwordUpdate.OldPW = ""

	// look up the user
	user, err := u.store.GetAppUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// compare passwords
	if !auth.CompareHashAndPassword(user.PasswordHash, passwordUpdate.OldPW) {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid password"})
		return
	}

	// generate hashed NEW password
	newPasswordHash, err := auth.HashPassword(passwordUpdate.NewPW)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	passwordUpdate.OldPW = ""
	passwordUpdate.NewPW = ""

	// update password
	if err := u.store.UpdateUserPassword(userID, newPasswordHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// GetUser retrieves the user by ID
func (u *Users) GetUser(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user"})
		return
	}

	user, err := u.store.GetAppUserByID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// ListTiers retrieves the list available tiers for this account/user
func (u *Users) ListTiers(c *gin.Context) {
	tiers, err := u.store.ListTiers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tiers": tiers})
}

// GetUsage returns the current usage for the account/user
func (u *Users) GetUsage(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user"})
		return
	}

	hour := time.Now().Hour()
	currentKey := fmt.Sprintf("requestCount:%s:%d", userID, hour)

	// get the request count key for the current window
	val, err := u.cache.Get(currentKey).Int()

	if err == redis.Nil {
		// {currentKey} does not exist
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error retrieving usage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": val, "projects": 0, "storage": 0})
}

// GetSession retrieves the user's current session information.
func (u *Users) GetSession(c *gin.Context) {
	sessionID := c.Param("sessionID")

	session, err := u.store.GetAppSession(sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"session": session})
}

// ListUserSessions returns a list of the user's active sessions
func (u *Users) ListUserSessions(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user"})
		return
	}

	sessions, err := u.store.ListUserSessions(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}
