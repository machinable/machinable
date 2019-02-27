package users

import (
	"errors"

	"github.com/anothrnick/machinable/auth"
)

type newUserBody struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	ReCaptcha string `json:"recaptcha"`
}

// Validate checks that the new user has a username and password.
func (u *newUserBody) Validate() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("invalid username or password")
	}

	if u.ReCaptcha == "" {
		return errors.New("recaptcha must be submitted")
	}

	if err := auth.RecaptchaSiteVerify(u.ReCaptcha); err != nil {
		return err
	}

	return nil
}
