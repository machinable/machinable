package models

import "time"

// NewMetaData returns a pointer to a new MetaData object with the `Created` field set to now.
func NewMetaData(creator, creatorType string) *MetaData {
	return &MetaData{
		Creator:     creator,
		CreatorType: creatorType,
		Created:     time.Now().Unix(),
	}
}

// MetaData contains internal data about a collection/resource object.
type MetaData struct {
	Creator     string `json:"creator"`
	CreatorType string `json:"creator_type"`
	Created     int64  `json:"created"`
}

// Map returns the metadata object as a map[string]interface{}
func (md *MetaData) Map() map[string]interface{} {
	return map[string]interface{}{
		"creator":      md.Creator,
		"creator_type": md.CreatorType,
		"created":      md.Created,
	}
}
