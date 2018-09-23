package handlers

import (
	"errors"
	"fmt"
	"time"

	"bitbucket.org/nsjostrom/machinable/database"
	"bitbucket.org/nsjostrom/machinable/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// supportedTypes is the list of supported resource field types, this will include any other
// resource definitions that have been created ("foreign key")
var supportedTypes = []string{"int", "float", "date", "bool", "string"}

// supportedType returns true if the string is a supported type, false otherwise.
func supportedType(a string) bool {
	for _, b := range supportedTypes {
		if b == a {
			return true
		}
	}
	return false
}

// validateResourceDefinition validates the fields of a resource definition.
func validateResourceDefinition(def *models.ResourceDefinition) error {
	if def.Name == "" {
		return errors.New("resource name cannot be empty")
	} else if def.PathName == "" {
		return errors.New("resource path_name cannot be empty")
	} else if len(def.Fields) == 0 {
		return errors.New("resource fields cannot be empty")
	}

	return nil
}

// processFields processes a map[string]string to a slice of *bson.Element for storing in mongo
func processFields(fields map[string]string) ([]*bson.Element, error) {
	fieldElements := make([]*bson.Element, 0)
	for key, value := range fields {
		if !supportedType(value) {
			return nil, fmt.Errorf("'%s' is not a supported field type", value)
		}
		fieldElements = append(fieldElements, bson.EC.String(key, value))
	}

	return fieldElements, nil
}

// getDefinition returns the *model.ResourceDefinition for a resource definition path name
func getDefinition(resourcePathName string) (*models.ResourceDefinition, error) {
	collection := database.Collection(database.ResourceDefinitions)

	// Find the resource definition for this object
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.String("path_name", resourcePathName),
		),
		nil,
	)

	if documentResult == nil {
		return nil, errors.New("no documents for resource")
	}

	// Decode result into document
	doc := bson.Document{}
	err := documentResult.Decode(&doc)
	if err != nil {
		return nil, errors.New("no documents for resource")
	}
	// Lookup field definitions for this resource
	resourceDefinition, err := parseDefinition(&doc)
	if err != nil {
		return nil, err
	}

	return resourceDefinition, nil
}

// getDefinitionByID returns the *model.ResourceDefinition by resource definition ID
func getDefinitionByID(resourceID string) (*models.ResourceDefinition, error) {
	objectID, err := objectid.FromHex(resourceID)
	if err != nil {
		return nil, err
	}
	collection := database.Collection(database.ResourceDefinitions)

	// Find the resource definition for this object
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		nil,
	)

	if documentResult == nil {
		return nil, errors.New("no documents for resource")
	}

	// Decode result into document
	doc := bson.Document{}
	err = documentResult.Decode(&doc)
	if err != nil {
		return nil, errors.New("no documents for resource")
	}
	// Lookup field definitions for this resource
	resourceDefinition, err := parseDefinition(&doc)
	if err != nil {
		return nil, err
	}

	return resourceDefinition, nil
}

// parseDefinition parses the *bson.Document of the resource definition to a *models.ResourceDefinition struct
func parseDefinition(doc *bson.Document) (*models.ResourceDefinition, error) {
	def := models.ResourceDefinition{
		Fields: make(map[string]string),
	}
	def.ID = doc.Lookup("_id").ObjectID().Hex()
	def.Name = doc.Lookup("name").StringValue()
	def.PathName = doc.Lookup("path_name").StringValue()
	fields := make(map[string]string)
	fieldsVal, err := doc.LookupErr("fields")
	if err != nil {
		return nil, fmt.Errorf("could not get fields: %s", err.Error())
	}
	fieldsDoc, ok := fieldsVal.MutableDocumentOK()
	if !ok {
		return nil, errors.New("could not get mutable document from fields value")
	}

	fieldKeys, _ := fieldsDoc.Keys(false)
	for _, key := range fieldKeys {
		fields[key.String()] = fieldsDoc.Lookup(key.String()).StringValue()
	}
	def.Fields = fields

	return &def, nil
}

// Int64 attempts to cast the interface to a int64
func Int64(unk interface{}) (int64, error) {
	switch unk.(type) {
	case int64:
		return unk.(int64), nil
	case int32:
		val := unk.(int32)
		return int64(val), nil
	case int:
		val := unk.(int)
		return int64(val), nil
	case uint:
		val := unk.(uint)
		return int64(val), nil
	case float64:
		val := unk.(float64)
		return int64(val), nil
	case float32:
		val := unk.(float32)
		return int64(val), nil
	default:
		return -1, errors.New("unknown value is of incompatible type, int64")
	}
}

// Float64 attempts to cast the interface to a float64
func Float64(unk interface{}) (float64, error) {
	switch unk.(type) {
	case float64:
		return unk.(float64), nil
	case int64:
		val := unk.(int32)
		return float64(val), nil
	case int32:
		val := unk.(int32)
		return float64(val), nil
	case int:
		val := unk.(int)
		return float64(val), nil
	case uint:
		val := unk.(uint)
		return float64(val), nil
	case float32:
		val := unk.(float32)
		return float64(val), nil
	default:
		return -1, errors.New("unknown value is of incompatible type, float64")
	}
}

// createFieldDocument creates a *bson.Document of the object fields based on their defined type
func createFieldDocument(fields map[string]interface{}, types map[string]string) (*bson.Document, error) {
	resourceElements := make([]*bson.Element, 0)

	// Iterate types and parse fields into document
	for key, ftype := range types {
		value, ok := fields[key]
		if !ok {
			return nil, fmt.Errorf("resource field not found in body '%s'", key)
		}
		switch ftype {
		case "int":
			valueAssert, err := Int64(value)
			if err != nil {
				return nil, fmt.Errorf("invalid type on '%s': %s", key, err.Error())
			}
			resourceElements = append(resourceElements, bson.EC.Int64(key, valueAssert))
		case "float":
			valueAssert, err := Float64(value)
			if err != nil {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			resourceElements = append(resourceElements, bson.EC.Double(key, valueAssert))
		case "date":
			rfc3339, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			valueAssert, err := time.Parse(time.RFC3339, rfc3339)
			if err != nil {
				return nil, fmt.Errorf("invalid type on '%s', cannot parse date", key)
			}
			resourceElements = append(resourceElements, bson.EC.Time(key, valueAssert))
		case "bool":
			valueAssert, ok := value.(bool)
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			resourceElements = append(resourceElements, bson.EC.Boolean(key, valueAssert))
		case "string":
			valueAssert, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			resourceElements = append(resourceElements, bson.EC.String(key, valueAssert))
		default:
			return nil, fmt.Errorf("unsupported type '%s'", ftype)
		}
	}

	return bson.NewDocument(resourceElements...), nil
}

// parseDocumentToMap parses the object *bson.Document to a map for JSON marshalling
func parseDocumentToMap(doc *bson.Document, types map[string]string) (map[string]interface{}, error) {
	// Create field map for this document
	fields := make(map[string]interface{})

	// Lookup ID and set field
	idValue, err := doc.LookupErr("_id")
	if err != nil {
		return nil, fmt.Errorf("error looking up field '_id', '%s'", err.Error())
	}
	fields["id"] = idValue.ObjectID().Hex()

	// Iterate types and parse fields
	// NOTE: this will ignore any fields that are not in the resource definition
	for key, ftype := range types {
		switch ftype {
		case "int":
			value, err := doc.LookupErr(key)
			if err != nil {
				return nil, fmt.Errorf("error looking up field '%s', '%s'", key, err.Error())
			}
			typedValue, ok := value.Int64OK()
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			fields[key] = typedValue
		case "float":
			value, err := doc.LookupErr(key)
			if err != nil {
				return nil, fmt.Errorf("error looking up field '%s', '%s'", key, err.Error())
			}
			typedValue, ok := value.DoubleOK()
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			fields[key] = typedValue
		case "date":
			value, err := doc.LookupErr(key)
			if err != nil {
				return nil, fmt.Errorf("error looking up field '%s', '%s'", key, err.Error())
			}
			typedValue, ok := value.TimeOK()
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			fields[key] = typedValue
		case "bool":
			value, err := doc.LookupErr(key)
			if err != nil {
				return nil, fmt.Errorf("error looking up field '%s', '%s'", key, err.Error())
			}
			typedValue, ok := value.BooleanOK()
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			fields[key] = typedValue
		case "string":
			value, err := doc.LookupErr(key)
			if err != nil {
				return nil, fmt.Errorf("error looking up field '%s', '%s'", key, err.Error())
			}
			typedValue, ok := value.StringValueOK()
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			fields[key] = typedValue
		default:
			return nil, fmt.Errorf("unsupported type '%s'", ftype)
		}
	}
	return fields, nil
}
