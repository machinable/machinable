package handlers

import (
	"errors"
	"fmt"

	"bitbucket.org/nsjostrom/machinable/projects/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const (
	// DocumentIDKey is the key of ids in mongodb
	DocumentIDKey = "_id"
	// MaxRecursion is the maximum amount of levels allowed in a JSON object (array and objects)
	MaxRecursion = 8
)

// supportedTypes is the list of supported resource field types, this will include any other
// resource definitions that have been created ("foreign key")
var supportedTypes = []string{"integer", "number", "boolean", "string", "array", "object"}
var supportedArrayItemTypes = []string{"integer", "number", "boolean", "string"}
var supportedFormats = []string{"date-time", "email", "hostname", "ipv4", "ipv6"}
var reservedFieldKeys = []string{"id", DocumentIDKey, "ID"}

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
func processProperties(properties map[string]models.Property, layer int) ([]*bson.Element, error) {
	propertyElements := make([]*bson.Element, 0)

	// this object goes deeper than supported
	if layer > MaxRecursion {
		return propertyElements, nil
	}

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

		properties := make([]*bson.Element, 0)

		if prop.Properties != nil {
			// Process the resource fields into bson
			var err error
			properties, err = processProperties(prop.Properties, layer+1)
			if err != nil {
				return nil, errors.New("could not process property's properties")
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
				bson.EC.SubDocumentFromElements("properties", properties...),
			)))
	}

	return propertyElements, nil
}

// definitionExists returns true if a definition already exists with path_name or name
func definitionExists(definition *models.ResourceDefinition, collection *mongo.Collection) bool {
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
func getDefinition(resourcePathName string, collection *mongo.Collection) (*models.ResourceDefinition, error) {

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
func getDefinitionByID(resourceID string, collection *mongo.Collection) (*models.ResourceDefinition, error) {
	objectID, err := objectid.FromHex(resourceID)
	if err != nil {
		return nil, err
	}

	// Find the resource definition for this object
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID(DocumentIDKey, objectID),
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

// TODO: limit recursion?
func propertyDocumentToModel(doc *bson.Document) (*models.Property, error) {
	prop := models.Property{}
	prop.Description = doc.Lookup("description").StringValue()
	prop.Format = doc.Lookup("format").StringValue()
	prop.Type = doc.Lookup("type").StringValue()
	items, err := getMutableDocument("items", doc)
	if err == nil {
		prop.Items = &models.Items{Type: items.Lookup("type").StringValue()}
	}

	if prop.Type == "object" {
		propertiesDoc, err := getMutableDocument("properties", doc)
		if err != nil {
			return nil, err
		}

		// RECURSION
		properties, err := parseDefinitionProperties(propertiesDoc)
		if err != nil {
			return nil, err
		}

		prop.Properties = properties
	}

	return &prop, nil
}

func parseDefinitionProperties(propertiesDoc *bson.Document) (map[string]models.Property, error) {
	properties := make(map[string]models.Property)
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

	return properties, nil
}

// parseDefinition parses the *bson.Document of the resource definition to a *models.ResourceDefinition struct
func parseDefinition(doc *bson.Document) (*models.ResourceDefinition, error) {
	def := models.ResourceDefinition{
		Properties: make(map[string]models.Property),
	}
	def.ID = doc.Lookup(DocumentIDKey).ObjectID().Hex()
	def.Title = doc.Lookup("title").StringValue()
	def.PathName = doc.Lookup("path_name").StringValue()

	propertiesDoc, err := getMutableDocument("properties", doc)
	if err != nil {
		return nil, err
	}

	properties, err := parseDefinitionProperties(propertiesDoc)
	if err != nil {
		return nil, err
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

// createVCType creates a bson.Value from the interface based on the `propType`
func createVCType(propType, value interface{}) (*bson.Value, error) {
	switch propType {
	case "integer":
		valueAssert, err := Int64(value)
		if err != nil {
			return nil, fmt.Errorf("invalid value for type '%s': %s", propType, err.Error())
		}
		return bson.VC.Int64(valueAssert), nil
	case "number":
		valueAssert, err := Float64(value)
		if err != nil {
			return nil, fmt.Errorf("invalid value for type '%s'", propType)
		}
		return bson.VC.Double(valueAssert), nil
	case "boolean":
		valueAssert, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid value for type '%s'", propType)
		}
		return bson.VC.Boolean(valueAssert), nil
	case "string":
		// TODO: check based on `format` definition
		valueAssert, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid value for type '%s'", propType)
		}
		return bson.VC.String(valueAssert), nil
	default:
		return nil, fmt.Errorf("unsupported value for type '%s'", propType)
	}
}

// createPropertyBSONElement creates a bson.Element for the interface based on the property definition
func createPropertyBSONElement(property *models.Property, key string, value interface{}, layer int) (*bson.Element, error) {
	switch property.Type {
	case "integer":
		valueAssert, err := Int64(value)
		if err != nil {
			return nil, fmt.Errorf("invalid type on '%s': %s", key, err.Error())
		}
		return bson.EC.Int64(key, valueAssert), nil
	case "number":
		valueAssert, err := Float64(value)
		if err != nil {
			return nil, fmt.Errorf("invalid type on '%s'", key)
		}
		return bson.EC.Double(key, valueAssert), nil
	case "boolean":
		valueAssert, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid type on '%s'", key)
		}
		return bson.EC.Boolean(key, valueAssert), nil
	case "string":
		// TODO: validate value based on `format` definition
		valueAssert, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type on '%s'", key)
		}
		return bson.EC.String(key, valueAssert), nil
	case "array":
		valueAssert, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type on '%s'", key)
		}

		bsonArray := bson.NewArray()

		// populate bson array with correct type, based on property definition
		for _, arrValue := range valueAssert {
			bcValue, err := createVCType(property.Items.Type, arrValue)
			if err != nil {
				return nil, fmt.Errorf("invalid type on array items for '%s', %s required", key, property.Items.Type)
			}
			bsonArray.Append(bcValue)
		}

		return bson.EC.Array(key, bsonArray), nil
	case "object":
		// TODO: createPropertyDocument on object
		propertyProperties := property.Properties
		valueAssert, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type on '%s'", key)
		}

		propDoc, err := createPropertyDocument(valueAssert, propertyProperties, layer+1)
		if err != nil {
			return nil, fmt.Errorf("invalid type on object items for '%s'", key)
		}

		return bson.EC.SubDocument(key, propDoc), err
	default:
		return nil, fmt.Errorf("unsupported type '%s'", property.Type)
	}
}

// createPropertyDocument creates a *bson.Document of the object fields based on their defined type
func createPropertyDocument(fields map[string]interface{}, types map[string]models.Property, layer int) (*bson.Document, error) {
	resourceElements := make([]*bson.Element, 0)

	// this object goes deeper than supported
	if layer > MaxRecursion {
		return bson.NewDocument(resourceElements...), nil
	}

	// Iterate types and parse fields into document
	for key, property := range types {
		value, ok := fields[key]
		if !ok {
			return nil, fmt.Errorf("resource field not found in body '%s'", key)
		}

		// create element from property
		element, err := createPropertyBSONElement(&property, key, value, layer)

		if err != nil {
			return nil, err
		}

		// append to list of elements
		resourceElements = append(resourceElements, element)
	}

	// create new document from elements
	return bson.NewDocument(resourceElements...), nil
}

// parseUnknownDocumentToMap parses the bson.Array to a []interface{}, recursively
func parseUnknownArrayToInterfaces(arrValue *bson.Array, layer int) ([]interface{}, error) {
	interfaceArr := make([]interface{}, 0)

	// this object goes deeper than supported
	if layer > MaxRecursion {
		return interfaceArr, nil
	}

	itr, err := arrValue.Iterator()
	if err != nil {
		return nil, errors.New("could not parse array iterator")
	}

	// create array of interfaces, this will end up marshalling the proper type in the JSON response
	for itr.Next() {
		itrVal := itr.Value()
		if arrValue, ok := itrVal.MutableArrayOK(); ok {
			// this is an array and we need to get the interface of each element
			// recursively call parseUnknownArrayToInterfaces
			// recursion limited to `MaxRecursion` levels
			arrIntr, err := parseUnknownArrayToInterfaces(arrValue, layer+1)
			if err == nil {
				interfaceArr = append(interfaceArr, arrIntr)
			}
		} else if objValue, ok := itrVal.MutableDocumentOK(); ok {
			// this is an object and we need to get the interface of each element
			// recursively call parseUnknownDocumentToMap
			// recursion limited to `MaxRecursion` levels
			mObj, err := parseUnknownDocumentToMap(objValue, layer+1)
			if err == nil {
				interfaceArr = append(interfaceArr, mObj)
			}
		} else {
			// this is a primitive value, append the interface to the array
			interfaceArr = append(interfaceArr, itrVal.Interface())
		}
	}

	return interfaceArr, nil
}

// parseUnknownDocumentToMap parses the bson.Document to a map[string]interface{}, recursively
func parseUnknownDocumentToMap(doc *bson.Document, layer int) (map[string]interface{}, error) {
	keyVals := make(map[string]interface{})

	// this object goes deeper than supported
	if layer > MaxRecursion {
		return keyVals, nil
	}

	keys, err := doc.Keys(false)

	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		docVal := doc.Lookup(key.String())

		if key.String() == DocumentIDKey {
			keyVals["id"] = docVal.ObjectID().Hex()
		} else {
			if arrValue, ok := docVal.MutableArrayOK(); ok {
				// this is an array and we need to get the interface of each element
				// recursively call parseUnknownArrayToInterfaces
				// recursion limited to `MaxRecursion` levels
				arrIntr, err := parseUnknownArrayToInterfaces(arrValue, layer+1)
				if err == nil {
					keyVals[key.String()] = arrIntr
				}
			} else if objValue, ok := docVal.MutableDocumentOK(); ok {
				// this is an object and we need to get the interface of each element
				// recursively call parseUnknownDocumentToMap
				// recursion limited to `MaxRecursion` levels
				mObj, err := parseUnknownDocumentToMap(objValue, layer+1)
				if err == nil {
					keyVals[key.String()] = mObj
				}
			} else {
				keyVals[key.String()] = docVal.Interface()
			}
		}
	}

	return keyVals, nil
}
