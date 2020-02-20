package auth

import (
	"fmt"
	"time"

	"github.com/anothrnick/machinable/config"
	"github.com/dgrijalva/jwt-go"
)

const (
	// AccessTokenExpiry is how long access tokens last.
	AccessTokenExpiry = 15 // 15 minutes
	// RefreshTokenExpiry is the period of time that refresh tokens are valid.
	RefreshTokenExpiry = 60 * 24 * 3 // 3 days
)

// JWT wraps functions need to create and parse JWTs
type JWT struct {
	config *config.AppConfig
}

// NewJWT creates and returns a pointer to a new `JWT`
func NewJWT(config *config.AppConfig) *JWT {
	return &JWT{
		config: config,
	}
}

// TokenLookup returns an error if the JWT is invalid.
func (j *JWT) TokenLookup(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(j.config.AppSecret), nil
}

// CreateRefreshToken creates a new JWT to be used as the refresh token. The refresh token can be used to retrieve a new access token.
func (j *JWT) CreateRefreshToken(sessionID, userID string) (string, error) {
	// create refresh jwt
	claims := jwt.MapClaims{
		"session_id": sessionID,
		"user_id":    userID,
	}
	expiry := time.Now().Add(time.Minute * RefreshTokenExpiry).Unix()

	return j.CreateJWT(claims, expiry)
}

// CreateAccessToken creates a new JWT to be used as an access token.
func (j *JWT) CreateAccessToken(claims jwt.MapClaims) (string, error) {
	expiry := time.Now().Add(time.Minute * AccessTokenExpiry).Unix()
	return j.CreateJWT(claims, expiry)
}

// CreateJWT creates a new JWT. This function will add the `exp` key to the claims based on the expiry time.
func (j *JWT) CreateJWT(claims jwt.MapClaims, expiry int64) (string, error) {
	claims["exp"] = expiry
	// create access jwt
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(j.config.AppSecret))

	return tokenString, err
}
