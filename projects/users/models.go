package users

import "errors"

var validRoles = []string{"user", "admin"}

// NewProjectUser is the JSON structure of a new user request
type NewProjectUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Read     bool   `json:"read"`
	Write    bool   `json:"write"`
	Role     string `json:"role"`
}

// SupportedRole verifies that the users role is valid and supported
func (u *NewProjectUser) SupportedRole(a string) bool {
	for _, b := range validRoles {
		if b == a {
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

	return nil
}
