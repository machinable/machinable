package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// SecretKey used for JWT signing. This should be a config or environment variable.
	SecretKey = "YmQzYTNkYzExMTVmZTc1YzY0NzY0NGU4"
	// AccessTokenExpiry is how long access tokens last.
	AccessTokenExpiry = 15 // 15 minutes
	// RefreshTokenExpiry is the period of time that refresh tokens are valid.
	RefreshTokenExpiry = 60 * 24 * 3 // 3 days
)

// TokenLookup returns an error if the JWT is invalid.
func TokenLookup(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(SecretKey), nil
}

// CreateRefreshToken creates a new JWT to be used as the refresh token. The refresh token can be used to retrieve a new access token.
func CreateRefreshToken(sessionID, userID string) (string, error) {
	// create refresh jwt
	claims := jwt.MapClaims{
		"session_id": sessionID,
		"user_id":    userID,
	}
	expiry := time.Now().Add(time.Minute * RefreshTokenExpiry).Unix()

	return CreateJWT(claims, expiry)
}

// CreateAccessToken creates a new JWT to be used as an access token.
func CreateAccessToken(claims jwt.MapClaims) (string, error) {
	expiry := time.Now().Add(time.Minute * AccessTokenExpiry).Unix()
	return CreateJWT(jwt.MapClaims{}, expiry)
}

// CreateJWT creates a new JWT. This function will add the `exp` key to the claims based on the expiry time.
func CreateJWT(claims jwt.MapClaims, expiry int64) (string, error) {
	claims["exp"] = expiry
	// create access jwt
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(SecretKey))

	return tokenString, err
}
