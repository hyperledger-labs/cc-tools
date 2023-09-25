package assets

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger-labs/cc-tools/errors"
)

// validateProp checks if a given assetProp is valid according to the given property definition
func validateProp(prop interface{}, propDef AssetProp) (interface{}, error) {
	var isArray bool
	dataTypeName := propDef.DataType
	if strings.HasPrefix(dataTypeName, "[]") {
		dataTypeName = strings.TrimPrefix(dataTypeName, "[]")
		isArray = true
	}

	var isSubAsset bool
	if strings.HasPrefix(dataTypeName, "->") {
		dataTypeName = strings.TrimPrefix(dataTypeName, "->")
		isSubAsset = true
	}

	var retProp interface{}
	var err error

	// Handle array-like properties
	var propAsArray []interface{}
	if !isArray {
		propAsArray = []interface{}{prop}
	} else {
		propReflectValue := reflect.ValueOf(prop)
		if propReflectValue.Kind() != reflect.Slice || propReflectValue.IsNil() {
			return nil, errors.NewCCError(fmt.Sprintf("asset property '%s' must be a slice", propDef.Label), 400)
		}
		propAsArray = make([]interface{}, propReflectValue.Len())
		for i := 0; i < propReflectValue.Len(); i++ {
			propAsArray[i] = propReflectValue.Index(i).Interface()
		}
		retProp = []interface{}{}
	}

	for _, prop := range propAsArray {
		var parsedProp interface{}

		if prop == nil {
			continue
		}

		// Validate data types
		if !isSubAsset {
			dataType, dataTypeExists := dataTypeMap[dataTypeName]
			if !dataTypeExists {
				return nil, errors.NewCCError(fmt.Sprintf("invalid data type named '%s'", propDef.DataType), 400)
			}

			_, parsedProp, err = dataType.Parse(prop)
			if err != nil {
				return nil, errors.WrapError(err, fmt.Sprintf("invalid '%s' (%s) asset property", propDef.Tag, propDef.Label))
			}
		} else {
			// Check if type is defined in assetList
			subAssetType := FetchAssetType(dataTypeName)
			if subAssetType == nil {
				return nil, errors.NewCCError(fmt.Sprintf("invalid asset type named '%s'", propDef.DataType), 400)
			}

			// Check if received subAsset is a map
			var recvMap map[string]interface{}
			switch t := prop.(type) {
			case map[string]interface{}:
				recvMap = t
			case Key:
				recvMap = t
			case Asset:
				recvMap = t
			default:
				return nil, errors.NewCCError("asset reference must be an object", 400)
			}

			// Add assetType to received object
			recvMap["@assetType"] = dataTypeName

			// Check if all key props are included
			key, err := NewKey(recvMap)
			if err != nil {
				return nil, errors.WrapError(err, "error validating subAsset reference")
			}

			parsedProp = (map[string]interface{})(key)
		}

		// If prop has specific validation method, call it
		if propDef.Validate != nil {
			err := propDef.Validate(prop)
			if err != nil {
				errMsg := fmt.Sprintf("failed validating '%s' (%s)", propDef.Tag, propDef.Label)
				return nil, errors.WrapErrorWithStatus(err, errMsg, 400)
			}
		}

		if isArray {
			retProp = append(retProp.([]interface{}), parsedProp)
		} else {
			retProp = parsedProp
		}
	}

	return retProp, nil
}
