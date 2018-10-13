package models

import (
	"encoding/json"
)

// ResourceDefinition defines an API resource
type ResourceDefinition struct {
	ID         string              `json:"id"`                 // ID is the unique identifier for this resource definition
	Title      string              `json:"title"`              // Title of this resource
	PathName   string              `json:"path_name"`          // PathName is the name that will appear in the URL path
	Required   []string            `json:"required,omitempty"` // Required is an array of required properties
	Properties map[string]Property `json:"properties"`         // Fields is a map of name, type for each field of this resource
}

// Property describes a resource property
type Property struct {
	Type        string              `json:"type"`                  // Type is the type of this property, see `handlers.supportedTypes`
	Description string              `json:"description,omitempty"` // Description is a human readable description of this property
	Format      string              `json:"format,omitempty"`      // Format it the format of the value
	Items       *Items              `json:"items,omitempty"`       // Items describes the type of the items of this property (if it is an array)
	Properties  map[string]Property `json:"properties,omitempty"`  // Fields is a map of name, type for each field of this resource
}

// MarshalJSON is a custom JSON marshaller for Property
// The main purpose of this is to omit Items if the Items.Type is empty
func (p Property) MarshalJSON() ([]byte, error) {
	prop := struct {
		Type        string              `json:"type"`
		Description string              `json:"description,omitempty"`
		Format      string              `json:"format,omitempty"`
		Items       *Items              `json:"items,omitempty"`
		Properties  map[string]Property `json:"properties,omitempty"`
	}{}

	prop.Type = p.Type
	prop.Description = p.Description
	prop.Format = p.Format
	prop.Properties = p.Properties
	if p.Items != nil && p.Items.Type != "" {
		prop.Items = p.Items
	}

	return json.Marshal(&prop)
}

// Items describes the items of an array property
type Items struct {
	Type string `json:"type,omitempty"`
}
