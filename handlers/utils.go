package handlers

import (
	"errors"
	"fmt"

	"bitbucket.org/nsjostrom/machinable/database"
	"bitbucket.org/nsjostrom/machinable/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// supportedTypes is the list of supported resource field types, this will include any other
// resource definitions that have been created ("foreign key")
var supportedTypes = []string{"integer", "number", "boolean", "string"}
var supportedFormats = []string{"date-time", "email", "hostname", "ipv4", "ipv6"}
var reservedFieldKeys = []string{"id", "_id", "ID"}

// supportedType returns true if the string is a supported type, false otherwise.
func supportedType(a string) bool {
	for _, b := range supportedTypes {
		if b == a {
			return true
		}
	}
	return false
}

// reservedField returns true if the string is a reserved field key
func reservedField(a string) bool {
	for _, b := range reservedFieldKeys {
		if b == a {
			return true
		}
	}
	return false
}

// validateResourceDefinition validates the fields of a resource definition.
func validateResourceDefinition(def *models.ResourceDefinition) error {
	if def.Title == "" {
		return errors.New("resource title cannot be empty")
	} else if def.PathName == "" {
		return errors.New("resource path_name cannot be empty")
	} else if len(def.Properties) == 0 {
		return errors.New("resource properties cannot be empty")
	}

	return nil
}

// processProperties processes a map[string]string to a slice of *bson.Element for storing in mongo
func processProperties(properties map[string]models.Property) ([]*bson.Element, error) {
	propertyElements := make([]*bson.Element, 0)
	for key, prop := range properties {
		if !supportedType(prop.Type) {
			return nil, fmt.Errorf("'%s' is not a supported property type", prop.Type)
		}
		if reservedField(key) {
			return nil, fmt.Errorf("'%s' is a reserved property name", key)
		}
		itemsType := ""
		if prop.Items != nil {
			itemsType = prop.Items.Type
			if !supportedType(itemsType) {
				return nil, fmt.Errorf("'%s' is not a supported property items.type", itemsType)
			}
		}

		// TODO: `object` type, call processProperties on `properties`

		propertyElements = append(
			propertyElements,
			bson.EC.SubDocument(key, bson.NewDocument(
				bson.EC.String("type", prop.Type),
				bson.EC.String("description", prop.Description),
				bson.EC.String("format", prop.Format),
				bson.EC.SubDocument("items", bson.NewDocument(
					bson.EC.String("type", itemsType),
				)),
			)))
	}

	return propertyElements, nil
}

// definitionExists returns true if a definition already exists with path_name or name
func definitionExists(definition *models.ResourceDefinition) bool {
	collection := database.Collection(database.ResourceDefinitions)
	// Find the resource definition for this object
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ArrayFromElements("$or",
				bson.VC.DocumentFromElements(
					bson.EC.String("path_name", definition.PathName),
				),
				bson.VC.DocumentFromElements(
					bson.EC.String("title", definition.Title),
				),
			),
		),
		nil,
	)

	// Decode result into document
	doc := bson.Document{}
	err := documentResult.Decode(&doc)
	if err != nil {
		return false
	}

	return true
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

func getMutableDocument(key string, doc *bson.Document) (*bson.Document, error) {
	val, err := doc.LookupErr(key)
	if err != nil {
		return nil, fmt.Errorf("could not get '%s': %s", key, err.Error())
	}
	doc, ok := val.MutableDocumentOK()
	if !ok {
		return nil, fmt.Errorf("could not get mutable document from '%s' value", key)
	}

	return doc, nil
}

func propertyDocumentToModel(doc *bson.Document) (*models.Property, error) {
	prop := models.Property{}
	prop.Description = doc.Lookup("description").StringValue()
	prop.Format = doc.Lookup("format").StringValue()
	prop.Type = doc.Lookup("type").StringValue()
	items, err := getMutableDocument("items", doc)
	if err == nil {
		prop.Items = &models.Items{Type: items.Lookup("type").StringValue()}
	}

	return &prop, nil
}

// parseDefinition parses the *bson.Document of the resource definition to a *models.ResourceDefinition struct
func parseDefinition(doc *bson.Document) (*models.ResourceDefinition, error) {
	def := models.ResourceDefinition{
		Properties: make(map[string]models.Property),
	}
	def.ID = doc.Lookup("_id").ObjectID().Hex()
	def.Title = doc.Lookup("title").StringValue()
	def.PathName = doc.Lookup("path_name").StringValue()
	properties := make(map[string]models.Property)
	propertiesDoc, err := getMutableDocument("properties", doc)
	if err != nil {
		return nil, err
	}

	propertiesKeys, _ := propertiesDoc.Keys(false)
	for _, key := range propertiesKeys {
		propDoc, err := getMutableDocument(key.String(), propertiesDoc)
		if err != nil {
			return nil, err
		}
		property, err := propertyDocumentToModel(propDoc)
		if err != nil {
			return nil, err
		}
		properties[key.String()] = *property
	}
	def.Properties = properties

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

// createPropertyDocument creates a *bson.Document of the object fields based on their defined type
func createPropertyDocument(fields map[string]interface{}, types map[string]models.Property) (*bson.Document, error) {
	resourceElements := make([]*bson.Element, 0)

	// Iterate types and parse fields into document
	for key, property := range types {
		value, ok := fields[key]
		if !ok {
			return nil, fmt.Errorf("resource field not found in body '%s'", key)
		}
		switch property.Type {
		case "integer":
			valueAssert, err := Int64(value)
			if err != nil {
				return nil, fmt.Errorf("invalid type on '%s': %s", key, err.Error())
			}
			resourceElements = append(resourceElements, bson.EC.Int64(key, valueAssert))
		case "number":
			valueAssert, err := Float64(value)
			if err != nil {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			resourceElements = append(resourceElements, bson.EC.Double(key, valueAssert))
		// case "date":
		// 	rfc3339, ok := value.(string)
		// 	if !ok {
		// 		return nil, fmt.Errorf("invalid type on '%s'", key)
		// 	}
		// 	valueAssert, err := time.Parse(time.RFC3339, rfc3339)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("invalid type on '%s', cannot parse date", key)
		// 	}
		// 	resourceElements = append(resourceElements, bson.EC.Time(key, valueAssert))
		case "boolean":
			valueAssert, ok := value.(bool)
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			resourceElements = append(resourceElements, bson.EC.Boolean(key, valueAssert))
		case "string":
			// TODO: check based on `format` definition
			valueAssert, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			resourceElements = append(resourceElements, bson.EC.String(key, valueAssert))
		// case "array":
		// 	valueAssert, ok := value.([]interface{})
		// 	if !ok {
		// 		return nil, fmt.Errorf("invalid type on '%s'", key)
		// 	}
		// 	resourceElements = append(resourceElements, bson.EC.Array(key, bson.Array(valueAssert)))
		// case "object":
		// TODO: createPropertyDocument on object
		default:
			return nil, fmt.Errorf("unsupported type '%s'", property.Type)
		}
	}

	return bson.NewDocument(resourceElements...), nil
}

// parseDocumentToMap parses the object *bson.Document to a map for JSON marshalling
func parseDocumentToMap(doc *bson.Document, types map[string]models.Property) (map[string]interface{}, error) {
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
	for key, property := range types {
		switch property.Type {
		case "integer":
			value, err := doc.LookupErr(key)
			if err != nil {
				return nil, fmt.Errorf("error looking up field '%s', '%s'", key, err.Error())
			}
			typedValue, ok := value.Int64OK()
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			fields[key] = typedValue
		case "number":
			value, err := doc.LookupErr(key)
			if err != nil {
				return nil, fmt.Errorf("error looking up field '%s', '%s'", key, err.Error())
			}
			typedValue, ok := value.DoubleOK()
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			fields[key] = typedValue
		// case "date":
		// 	value, err := doc.LookupErr(key)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("error looking up field '%s', '%s'", key, err.Error())
		// 	}
		// 	typedValue, ok := value.TimeOK()
		// 	if !ok {
		// 		return nil, fmt.Errorf("invalid type on '%s'", key)
		// 	}
		// 	fields[key] = typedValue
		case "boolean":
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
			return nil, fmt.Errorf("unsupported type '%s'", property.Type)
		}
	}
	return fields, nil
}
