package dsi

import "regexp"

const (
	// JSONIDKey is the key of the id returned in the JSON response
	JSONIDKey = "id"
	// DocumentIDKey is the key of ids in mongodb
	DocumentIDKey = "_id"
	// LimitKey is used for paginating HTTP requests
	LimitKey = "_limit"
	// OffsetKey is used for paginating HTTP requests
	OffsetKey = "_offset"
	// MetadataKey is the key used to store internal metadata for an object
	MetadataKey         = "_metadata"
	MetadataCreated     = "_metadata.created"
	MetadataCreator     = "_metadata.creator"
	MetadataCreatorType = "_metadata.creator_type"

	// MaxRecursion is the maximum amount of levels allowed in a JSON object (array and objects)
	MaxRecursion = 8
)

// ValidPathFormat is the regular expression used to validate resource path names, collection names, and project slugs
var ValidPathFormat = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)

// reservedFieldKeys is the list of keys that cannot be used, as they are reserved for machinable use
var reservedFieldKeys = []string{JSONIDKey, DocumentIDKey, LimitKey, OffsetKey, MetadataKey, MetadataCreated, MetadataCreator, MetadataCreatorType}

// ReservedField returns true if the string is a reserved field key
func ReservedField(a string) bool {
	for _, b := range reservedFieldKeys {
		if b == a {
			return true
		}
	}
	return false
}
