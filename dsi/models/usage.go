package models

// Stats is a struct that contains the size and count of objects in a Resource or Collection
type Stats struct {
	Size  int64 `json:"size"`
	Count int64 `json:"count"`
}
