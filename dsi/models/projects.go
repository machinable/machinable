package models

import (
	"time"
)

// Project is an application project created and managed by a `User`
type Project struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Slug             string    `json:"slug"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Icon             string    `json:"icon"`
	Created          time.Time `json:"created"`
	Authn            bool      `json:"authn"`
	UserRegistration bool      `json:"user_registration"`
}
