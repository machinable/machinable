package models

// Collection is a mongo collection containing whatever data a user
// wants to save
type Collection struct {
	Name string `json:"name"`
}
