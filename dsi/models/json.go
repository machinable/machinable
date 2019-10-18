package models

// RootKey defines the metadata of a project root JSON key
type RootKey struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	ProjectID string `json:"project_id"`
	Create    bool   `json:"create"`
	Read      bool   `json:"read"`
	Update    bool   `json:"update"`
	Delete    bool   `json:"delete"`
}
