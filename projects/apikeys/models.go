package apikeys

import (
	"errors"

	"github.com/anothrnick/machinable/auth"
	uuid "github.com/satori/go.uuid"
)

// NewProjectKey is the JSON structure of a new api key request
type NewProjectKey struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Read        bool   `json:"read"`
	Write       bool   `json:"write"`
	Role        string `json:"role"`
}

// SupportedRole verifies that the users role is valid and supported
func (u *NewProjectKey) SupportedRole() bool {
	for _, b := range auth.ValidRoles {
		if b == u.Role {
			return true
		}
	}
	return false
}

// Validate checks that the new key is not empty
func (u *NewProjectKey) Validate() error {
	if u.Key == "" {
		return errors.New("invalid key")
	}

	if _, err := uuid.FromString(u.Key); err != nil {
		return errors.New("invalid key")
	}

	return u.ValidateRoleAccess()
}

// ValidateRoleAccess validates the role and access fields
func (u *NewProjectKey) ValidateRoleAccess() error {
	// Set default role
	if u.Role == "" {
		u.Role = auth.RoleUser
	}

	if ok := u.SupportedRole(); !ok {
		return errors.New("invalid role")
	}

	return nil
}
