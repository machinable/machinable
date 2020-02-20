package auth

import (
	"crypto/sha1"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

// CompareHashAndPassword compares the hash value to the password.
func CompareHashAndPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SHA1 hashes using sha1 algorithm
// Used for api keys
func SHA1(text, appSecret string) string {
	algorithm := sha1.New()
	algorithm.Write([]byte(text + appSecret)) // salt
	return hex.EncodeToString(algorithm.Sum(nil))
}
