package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/nsjostrom/machinable/database"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

func lookupFields(doc *bson.Document) (map[string]string, error) {
	// Lookup field definitions and parse into a map of fieldName, fieldType
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
	return fields, nil
}

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
			valueAssert, ok := value.(int64)
			if !ok {
				return nil, fmt.Errorf("invalid type on '%s'", key)
			}
			resourceElements = append(resourceElements, bson.EC.Int64(key, valueAssert))
		case "float":
			valueAssert, ok := value.(float64)
			if !ok {
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

func getFieldDefinitions(resourcePathName string) (map[string]string, error) {
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
	documentResult.Decode(&doc)
	// Lookup field definitions for this resource
	fieldDefinitions, err := lookupFields(&doc)
	if err != nil {
		return nil, err
	}

	return fieldDefinitions, nil
}

// AddObject creates a new document of the resource definition
func AddObject(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	fieldValues := make(map[string]interface{})

	c.BindJSON(&fieldValues)

	// Get field definitions for this resource
	fieldDefinitions, err := getFieldDefinitions(resourcePathName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create document for this resource based on the field definitions
	objectDocument, err := createFieldDocument(fieldValues, fieldDefinitions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the resources.{resourcePathName} collection
	rc := database.Collection(fmt.Sprintf(database.ResourceFormat, resourcePathName))
	result, err := rc.InsertOne(
		context.Background(),
		objectDocument,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the inserted ID for the response
	fieldValues["id"] = result.InsertedID.(objectid.ObjectID).Hex()

	c.JSON(http.StatusCreated, fieldValues)
}

// ListObjects returns the list of objects for a resource
func ListObjects(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	collection := database.Collection(fmt.Sprintf(database.ResourceFormat, resourcePathName))
	response := make([]map[string]interface{}, 0)

	// Find the resource definition for this object
	fieldDefinitions, err := getFieldDefinitions(resourcePathName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find all objects for this resource
	cursor, err := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	// Create response from documents
	doc := bson.NewDocument()
	for cursor.Next(context.Background()) {
		doc.Reset()
		err := cursor.Decode(doc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fields, err := parseDocumentToMap(doc, fieldDefinitions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response = append(response, fields)
	}
	c.JSON(http.StatusOK, gin.H{"items": response})
}

// GetObject returns a single object with the resourceID for this resource
func GetObject(c *gin.Context) {
	resourcePathName := c.Param("resourcePathName")
	resourceID := c.Param("resourceID")
	collection := database.Collection(fmt.Sprintf(database.ResourceFormat, resourcePathName))

	// Find the resource definition for this object
	fieldDefinitions, err := getFieldDefinitions(resourcePathName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create object ID from resource ID string
	objectID, err := objectid.FromHex(resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find object based on ID
	documentResult := collection.FindOne(
		nil,
		bson.NewDocument(
			bson.EC.ObjectID("_id", objectID),
		),
		nil,
	)

	if documentResult == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no documents for resource"})
	}

	// Decode result into document
	doc := bson.Document{}
	documentResult.Decode(&doc)
	// Lookup  definitions for this resource
	object, err := parseDocumentToMap(&doc, fieldDefinitions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, object)
}
