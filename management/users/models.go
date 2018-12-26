package users

import "errors"

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
