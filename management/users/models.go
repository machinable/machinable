package users

import (
	"errors"

	"github.com/anothrnick/machinable/auth"
)

type updatePasswordBody struct {
	OldPW string `json:"old_password"`
	NewPW string `json:"new_password"`
}

// Validate checks that the new user has a username and password.
func (u *updatePasswordBody) Validate() error {
	if u.OldPW == "" {
		return errors.New("invalid current password")
	}
	if u.NewPW == "" {
		return errors.New("invalid new password")
	}

	return nil
}

type newUserBody struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ReCaptcha string `json:"recaptcha"`
}

// Validate checks that the new user has a username and password.
func (u *newUserBody) Validate(reCaptchaSecret string) error {
	if u.Email == "" || u.Password == "" {
		return errors.New("invalid email or password")
	}

	if u.ReCaptcha == "" {
		return errors.New("recaptcha must be submitted")
	}

	if err := auth.RecaptchaSiteVerify(u.ReCaptcha, reCaptchaSecret); err != nil {
		return err
	}

	return nil
}
