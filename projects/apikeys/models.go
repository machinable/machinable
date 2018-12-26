package apikeys

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

// NewProjectKey is the JSON structure of a new api key request
type NewProjectKey struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Read        bool   `json:"read"`
	Write       bool   `json:"write"`
}

// Validate checks that the new key is not empty
func (u *NewProjectKey) Validate() error {
	if u.Key == "" {
		return errors.New("invalid key")
	}

	if _, err := uuid.FromString(u.Key); err != nil {
		return errors.New("invalid key")
	}

	return nil
}
