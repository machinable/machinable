package models

// Collection is a mongo collection containing whatever data a user
// wants to save
type Collection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
