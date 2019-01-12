package dsi

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
	MetadataKey = "_metadata"
	// MaxRecursion is the maximum amount of levels allowed in a JSON object (array and objects)
	MaxRecursion = 8
)

// supportedTypes is the list of supported resource field types
var supportedTypes = []string{"integer", "number", "boolean", "string", "array", "object"}

// supportedArrayItemTypes is the list of supported array types. This does not supported arrays of arrays
var supportedArrayItemTypes = []string{"integer", "number", "boolean", "string", "object"}

// supportedFormats is the list of supported String formats, which are used to validate the field value
var supportedFormats = []string{"date-time", "email", "hostname", "ipv4", "ipv6"}

// reservedFieldKeys is the list of keys that cannot be used, as they are reserved for machinable use
var reservedFieldKeys = []string{JSONIDKey, DocumentIDKey, LimitKey, OffsetKey, MetadataKey}

// SupportedType returns true if the string is a supported type, false otherwise.
func SupportedType(a string) bool {
	for _, b := range supportedTypes {
		if b == a {
			return true
		}
	}
	return false
}

// SupportedArrayType returns true if the string is a supported array type, false otherwise.
func SupportedArrayType(a string) bool {
	for _, b := range supportedArrayItemTypes {
		if b == a {
			return true
		}
	}
	return false
}

// ReservedField returns true if the string is a reserved field key
func ReservedField(a string) bool {
	for _, b := range reservedFieldKeys {
		if b == a {
			return true
		}
	}
	return false
}
