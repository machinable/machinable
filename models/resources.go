package models

// ResourceDefinition defines an API resource
type ResourceDefinition struct {
	ID       string            `json:"id"`        // ID is the unique identifier for this resource definition
	Name     string            `json:"name"`      // Name of this resource
	PathName string            `json:"path_name"` // PathName is the name that will appear in the URL path
	Fields   map[string]string `json:"fields"`    // Fields is a map of name, type for each field of this resource
}

/*
Supported Types:
int
float
date
bool
string

*/
