package mongo

import (
	"errors"

	"github.com/mongodb/mongo-go-driver/bson"
)

const (
	// DocumentIDKey is the key of ids in mongodb
	DocumentIDKey = "_id"
	// MaxRecursion is the maximum amount of levels allowed in a JSON object (array and objects)
	MaxRecursion = 8
)

// parseUnknownArrayToInterfaces parses the bson.Array to a []interface{}, recursively
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
