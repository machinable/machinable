package mongo

import (
	"errors"
	"fmt"

	"github.com/anothrnick/machinable/dsi"
	"github.com/anothrnick/machinable/dsi/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func appendMetadata(doc *bson.Document, metadata *models.MetaData) {
	md := bson.NewDocument(
		bson.EC.String("creator", metadata.Creator),
		bson.EC.String("creator_type", metadata.CreatorType),
		bson.EC.Int64("created", metadata.Created),
	)

	doc.Append(bson.EC.SubDocument(
		"_metadata",
		md,
	))
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

// parseDefinition parses the *bson.Document of the resource definition to a *models.ResourceDefinition struct
func parseDefinition(doc *bson.Document) (*models.ResourceDefinition, error) {
	def := models.ResourceDefinition{}
	def.ID = doc.Lookup(dsi.DocumentIDKey).ObjectID().Hex()
	def.Title = doc.Lookup("title").StringValue()
	def.Created, _ = doc.Lookup("created").TimeOK()
	def.PathName = doc.Lookup("path_name").StringValue()
	def.ParallelRead, _ = doc.Lookup("parallel_read").BooleanOK()
	def.ParallelWrite, _ = doc.Lookup("parallel_write").BooleanOK()
	def.Properties = doc.Lookup("properties").StringValue()

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
