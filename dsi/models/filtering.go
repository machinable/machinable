package models

// Filters is a custom struct which contains query filters with the following structure:
// <field>: {
//    <operator>: <value>
// }
//
// This provides a datastore agnostic way to communicate filtering
const (
	GTE Op = "$gte"
	GT  Op = "$gt"
	LTE Op = "$lte"
	LT  Op = "$lt"
	EQ  Op = "$eq"
)

// Op represents the operation to filter with.
type Op string

// Filters is the map of query filters
type Filters map[string]Value

// AddFilter adds a filter
func (f *Filters) AddFilter(field string, value Value) {
	(*f)[field] = value
}

// Value is a map of `Op` to a value interface
type Value map[Op]interface{}
