package models

import (
	"time"
)

// Tier is an application tier, describing the limitations of an application subscribed to a certain tier
type Tier struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Cost     string `json:"cost"`
	Requests int    `json:"requests"`
	Projects int    `json:"projects"`
	Storage  int    `json:"storage"`
}

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

// ProjectDetail is read from the app_project_limits view and contains app tier values
// based on the currently active account tier
type ProjectDetail struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Slug             string    `json:"slug"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Icon             string    `json:"icon"`
	Created          time.Time `json:"created"`
	Authn            bool      `json:"authn"`
	UserRegistration bool      `json:"user_registration"`
	Requests         int       `json:"requests"`
}
