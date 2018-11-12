package dsi

const (
	// DocumentIDKey is the key of ids in mongodb
	DocumentIDKey = "_id"
	// MaxRecursion is the maximum amount of levels allowed in a JSON object (array and objects)
	MaxRecursion = 8
)

// reservedFieldKeys is the list of keys that cannot be used, as they are reserved for machinable use
var reservedFieldKeys = []string{"id", DocumentIDKey, "ID", "Id", "iD"}

// ReservedField returns true if the string is a reserved field key
func ReservedField(a string) bool {
	for _, b := range reservedFieldKeys {
		if b == a {
			return true
		}
	}
	return false
}
