package users

import (
	"errors"

	"github.com/machinable/machinable/auth"
)

// NewProjectUser is the JSON structure of a new user request
type NewProjectUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Read     bool   `json:"read"`
	Write    bool   `json:"write"`
	Role     string `json:"role"`
}

// SupportedRole verifies that the users role is valid and supported
func (u *NewProjectUser) SupportedRole() bool {
	for _, b := range auth.ValidRoles {
		if b == u.Role {
			return true
		}
	}
	return false
}

// Validate checks that the new user has a username and password.
func (u *NewProjectUser) Validate() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("invalid username or password")
	}

	return u.ValidateAccessRole()
}

// ValidateAccessRole validates the user's access and role
func (u *NewProjectUser) ValidateAccessRole() error {
	// Set default role
	if u.Role == "" {
		u.Role = auth.RoleUser
	}

	// Validate role
	if ok := u.SupportedRole(); !ok {
		return errors.New("invalid role")
	}

	return nil
}
