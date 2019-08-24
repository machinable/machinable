package models

import (
	"time"
)

// Session is a user session model for either the mgmt application or a project
type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Location     string    `json:"location"`
	Mobile       bool      `json:"mobile"`
	IP           string    `json:"ip"`
	LastAccessed time.Time `json:"last_accessed"`
	Browser      string    `json:"browser"`
	OS           string    `json:"os"`
}
