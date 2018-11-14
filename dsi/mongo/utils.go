package mongo

import (
	"errors"
	"fmt"

	"bitbucket.org/nsjostrom/machinable/dsi"
	"bitbucket.org/nsjostrom/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
)

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
	def.ID = doc.Lookup(dsi.DocumentIDKey).ObjectID().Hex()
	def.Title = doc.Lookup("title").StringValue()
	c, _ := doc.Lookup("created").TimeOK()
	def.Created = c
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

// parseUnknownArrayToInterfaces parses the bson.Array to a []interface{}, recursively
func parseUnknownArrayToInterfaces(arrValue *bson.Array, layer int) ([]interface{}, error) {
	interfaceArr := make([]interface{}, 0)

	// this object goes deeper than supported
	if layer > dsi.MaxRecursion {
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
	if layer > dsi.MaxRecursion {
		return keyVals, nil
	}

	keys, err := doc.Keys(false)

	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		docVal := doc.Lookup(key.String())

		if key.String() == dsi.DocumentIDKey {
			newID, ok := docVal.ObjectIDOK()
			if ok {
				keyVals["id"] = newID.Hex()
			} else {
				keyVals["id"] = docVal.Interface()
			}
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
