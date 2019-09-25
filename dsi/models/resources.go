package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anothrnick/machinable/dsi"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// ResourceObject is a custom type which wraps a map[string]interface
type ResourceObject map[string]interface{}

// Validate validates that the object matches the schema
func (obj *ResourceObject) Validate(definition *ResourceDefinition) error {
	if err := dsi.ContainsReservedField(*obj); err != nil {
		return err
	}

	schema := new(spec.Schema)

	err := json.Unmarshal([]byte(definition.Schema), schema)
	if err != nil {
		return err
	}

	data := map[string]interface{}{}
	for key, val := range *obj {
		data[key] = val
	}

	// validate data against schema
	res := validate.NewSchemaValidator(schema, nil, "", strfmt.Default).Validate(data)
	if res.HasErrors() {
		errs := []string{}
		for _, e := range res.Errors {
			errs = append(errs, e.Error())
		}
		return errors.New(strings.Join(errs, ","))
	}
	return nil
}

// JSONSchemaObject is a simplified representation of the root schema
type JSONSchemaObject struct {
	Type       string                            `json:"type"`
	Properties map[string]map[string]interface{} `json:"properties"`
	Required   []string                          `json:"required"`
	// AdditionalProperties bool                              `json:"additionalProperties"`
}

// Property is a simplified representation of a JSON Schema property
// type Property struct {
// 	Type string `json:"type"`
// }

// ResourceDefinition defines an API resource
type ResourceDefinition struct {
	ID            string    `json:"id"` // ID is the unique identifier for this resource definition
	ProjectID     string    `json:"project_id"`
	Title         string    `json:"title"`     // Title of this resource
	PathName      string    `json:"path_name"` // PathName is the name that will appear in the URL path
	ParallelRead  bool      `json:"parallel_read"`
	ParallelWrite bool      `json:"parallel_write"`
	Create        bool      `json:"create"`
	Read          bool      `json:"read"`
	Update        bool      `json:"update"`
	Delete        bool      `json:"delete"`
	Created       time.Time `json:"created"` // Created is the timestamp the resource was created
	Schema        string    `json:"schema"`  // Properties is the string representation of the JSON schema properties
}

// GetSchema returns the schema as a `Schema` object
func (def *ResourceDefinition) GetSchema() (*JSONSchemaObject, error) {
	schema := &JSONSchemaObject{}
	err := json.Unmarshal([]byte(def.Schema), schema)

	return schema, err
}

// MarshalJSON custom marshaller to marshall properties to json
func (def *ResourceDefinition) MarshalJSON() ([]byte, error) {
	schema := JSONSchemaObject{}
	err := json.Unmarshal([]byte(def.Schema), &schema)
	if err != nil {
		panic(err)
	}

	return json.Marshal(&struct {
		ID            string           `json:"id"` // ID is the unique identifier for this resource definition
		ProjectID     string           `json:"project_id"`
		Title         string           `json:"title"`     // Title of this resource
		PathName      string           `json:"path_name"` // PathName is the name that will appear in the URL path
		ParallelRead  bool             `json:"parallel_read"`
		ParallelWrite bool             `json:"parallel_write"`
		Create        bool             `json:"create"`
		Read          bool             `json:"read"`
		Update        bool             `json:"update"`
		Delete        bool             `json:"delete"`
		Created       time.Time        `json:"created"` // Created is the timestamp the resource was created
		Schema        JSONSchemaObject `json:"schema"`  // Properties is the string representation of the JSON schema properties
	}{
		ID:            def.ID,
		ProjectID:     def.ProjectID,
		Title:         def.Title,
		PathName:      def.PathName,
		ParallelRead:  def.ParallelRead,
		ParallelWrite: def.ParallelWrite,
		Create:        def.Create,
		Read:          def.Read,
		Update:        def.Update,
		Delete:        def.Delete,
		Created:       def.Created,
		Schema:        schema,
	})
}

// UnmarshalJSON is a custom unmarshaller
func (def *ResourceDefinition) UnmarshalJSON(b []byte) error {
	payload := struct {
		Title         string          `json:"title"` // Title of this resource
		ProjectID     string          `json:"project_id"`
		PathName      string          `json:"path_name"` // PathName is the name that will appear in the URL path
		Schema        json.RawMessage `json:"schema"`    // Schema is the string representation of the JSON schema
		ParallelRead  bool            `json:"parallel_read"`
		ParallelWrite bool            `json:"parallel_write"`
		Create        bool            `json:"create"`
		Read          bool            `json:"read"`
		Update        bool            `json:"update"`
		Delete        bool            `json:"delete"`
	}{}

	err := json.Unmarshal(b, &payload)

	if err != nil {
		panic(err)
	}

	def.ProjectID = payload.ProjectID
	def.Title = payload.Title
	def.PathName = payload.PathName
	def.Schema = string(payload.Schema)
	def.ParallelRead = payload.ParallelRead
	def.ParallelWrite = payload.ParallelWrite
	def.Create = payload.Create
	def.Read = payload.Read
	def.Update = payload.Update
	def.Delete = payload.Delete

	return nil
}

// Validate validates the fields of a resource definition.
func (def *ResourceDefinition) Validate() error {
	if def.Title == "" {
		return errors.New("resource title cannot be empty")
	} else if len(def.Title) > dsi.MaxLengthOfCollectionInfo {
		return fmt.Errorf("resource title cannot be longer than %d characters", dsi.MaxLengthOfCollectionInfo)
	} else if def.PathName == "" {
		return errors.New("resource path_name cannot be empty")
	} else if len(def.PathName) > dsi.MaxLengthOfCollectionInfo {
		return fmt.Errorf("resource path_name cannot be longer than %d characters", dsi.MaxLengthOfCollectionInfo)
	} else if def.Schema == "" {
		return errors.New("resource schema cannot be empty")
	} else if !dsi.ValidPathFormat.MatchString(def.PathName) {
		return errors.New("invalid path name: only alphanumeric, dashes, and underscores allowed")
	}

	objectSchema := JSONSchemaObject{}
	err := json.Unmarshal([]byte(def.Schema), &objectSchema)

	m := map[string]interface{}{}
	rec, _ := json.Marshal(objectSchema.Properties)
	json.Unmarshal(rec, &m)

	if err := dsi.ContainsReservedField(m); err != nil {
		return err
	}

	schema := new(spec.Schema)

	err = json.Unmarshal([]byte(def.Schema), schema)

	return err
}
