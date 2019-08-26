package models

import (
	"time"
)

// ProjectAPIKey is a static key used to access resources and collections of the project
type ProjectAPIKey struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	KeyHash     string    `json:"-"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
	Read        bool      `json:"read"`
	Write       bool      `json:"write"`
	Role        string    `json:"role"`
}

// ProjectUser is a user of a project. A user can access resources and collections of the project.
type ProjectUser struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"project_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Created      time.Time `json:"created"`
	Read         bool      `json:"read"`
	Write        bool      `json:"write"`
	Role         string    `json:"role"`
}
