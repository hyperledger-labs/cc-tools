package assets

import (
	"fmt"
	"strings"

	"github.com/goledgerdev/cc-tools/errors"
)

func validateProp(prop interface{}, propDef AssetProp) (interface{}, error) {
	var isArray bool
	dataTypeName := propDef.DataType
	if strings.HasPrefix(dataTypeName, "[]") {
		dataTypeName = strings.TrimPrefix(dataTypeName, "[]")
		isArray = true
	}

	var retProp interface{}
	var ok bool
	var err error

	// Handle array-like properties
	var propAsArray []interface{}
	if !isArray {
		propAsArray = []interface{}{prop}
	} else {
		propAsArray, ok = prop.([]interface{})
		if !ok {
			return nil, errors.NewCCError(fmt.Sprintf("asset property '%s' must and array of type %s", propDef.Label, propDef.DataType), 400)
		}
		retProp = []interface{}{}
	}

	for _, prop := range propAsArray {
		var parsedProp interface{}

		// Validate data types
		dataType, dataTypeExists := dataTypeMap[dataTypeName]
		if dataTypeExists { // if data type is a primary data type
			parsedProp, err = dataType.Validate(prop)
			if err != nil {
				return nil, errors.WrapError(err, fmt.Sprintf("invalid '%s' (%s) asset property", propDef.Tag, propDef.Label))
			}
		} else {
			// If not a primary type, check if type is defined in assetMap
			subAssetType := FetchAssetType(dataTypeName)
			if subAssetType == nil {
				err := fmt.Errorf("invalid data type named '%s'", propDef.DataType)
				return nil, errors.WrapErrorWithStatus(err, "invalid property type", 400)
			}

			// Check if received subAsset is a map
			recvMap, isMap := prop.(map[string]interface{})
			if !isMap {
				return nil, errors.NewCCError("asset reference must be sent as a JSON object", 400)
			}

			// Add assetType to received object
			recvMap["@assetType"] = dataTypeName

			// Check if all key props are included
			_, err := GenerateKey(recvMap)
			if err != nil {
				return nil, errors.WrapError(err, "error validating subAsset reference")
			}

			parsedProp = recvMap
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
