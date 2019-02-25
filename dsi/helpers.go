package dsi

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

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

// ContainsReservedField checks for reserved fields in a map[string]whatever
func ContainsReservedField(doc map[string]interface{}) error {
	for key := range doc {
		if ReservedField(key) {
			return fmt.Errorf("'%s' is a reserved field", key)
		}
	}
	return nil
}

// CastInterfaceToType returns the interface with a proper type
func CastInterfaceToType(typ string, value string) (interface{}, error) {
	switch typ {
	case "string":
		return value, nil
	case "integer":
		i, err := strconv.ParseInt(value, 10, 64)
		return i, cleanParseError(err)
	case "number":
		f, err := strconv.ParseFloat(value, 64)
		return f, cleanParseError(err)
	case "boolean":
		b, err := strconv.ParseBool(value)
		return b, cleanParseError(err)
	default:
		return value, errors.New("unable to filter on type")
	}
}

func cleanParseError(err error) error {
	if err == nil {
		return nil
	}

	return errors.New("error parsing value, invalid format")
}
