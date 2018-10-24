package models

import "time"

// Collection is a mongo collection containing whatever data a user
// wants to save
type Collection struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Items   int64     `json:"items"`
}
