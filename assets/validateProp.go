package assets

import (
	"fmt"
	"strings"
	"time"

	eh "github.com/goledgerdev/template-cc/chaincode/src/errorhandler"
)

func validateProp(prop interface{}, propDef AssetProp) (interface{}, error) {
	var isArray bool
	dataType := propDef.DataType
	if strings.HasPrefix(dataType, "[]") {
		dataType = strings.TrimPrefix(dataType, "[]")
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
			return nil, eh.NewCCError(400, fmt.Sprintf("asset property '%s' must and array of type %s", propDef.Label, propDef.DataType))
		}
		retProp = []interface{}{}
	}

	for _, prop := range propAsArray {
		var parsedProp interface{}

		// Validate data types
		switch dataType {
		case "string":
			parsedProp, ok = prop.(string)
			if !ok {
				err := fmt.Errorf("property '%s' (%s) must be a string", propDef.Tag, propDef.Label)
				return nil, eh.WrapErrorWithStatus(err, "invalid property type", 400)
			}
		case "number":
			parsedProp, ok = prop.(float64)
			if !ok {
				err := fmt.Errorf("property '%s' (%s) must be a number", propDef.Tag, propDef.Label)
				return nil, eh.WrapErrorWithStatus(err, "invalid property type", 400)
			}
		case "boolean":
			parsedProp, ok = prop.(bool)
			if !ok {
				err := fmt.Errorf("property '%s' (%s) must be a boolean", propDef.Tag, propDef.Label)
				return nil, eh.WrapErrorWithStatus(err, "invalid property type", 400)
			}
		case "datetime":
			propVal, ok := prop.(string)
			if !ok {
				return nil, eh.NewCCError(400, fmt.Sprintf("asset property %s should be an RFC3339 string", propDef.Label))
			}
			parsedProp, err = time.Parse(time.RFC3339, propVal)
			if err != nil {
				return nil, eh.WrapErrorWithStatus(err, fmt.Sprintf("invalid asset property %s RFC3339 format", propDef.Label), 400)
			}
		default:
			// If not a primary type, check if type is defined in assetMap
			subAssetType := FetchAssetType(dataType)
			if subAssetType == nil {
				err := fmt.Errorf("invalid data type named '%s'", propDef.DataType)
				return nil, eh.WrapErrorWithStatus(err, "invalid property type", 400)
			}

			// Check if received subAsset is a map
			recvMap, isMap := prop.(map[string]interface{})
			if !isMap {
				return nil, eh.NewCCError(400, "asset reference must be sent as a JSON object")
			}

			// Add assetType to received object
			recvMap["@assetType"] = dataType

			// Check if all key props are included
			_, err := GenerateKey(recvMap)
			if err != nil {
				return nil, eh.WrapError(err, "error validating subAsset reference")
			}

			parsedProp = recvMap
		}
		// If prop has specific validation method, call it
		if propDef.Validate != nil {
			err := propDef.Validate(prop)
			if err != nil {
				errMsg := fmt.Sprintf("failed validating '%s' (%s)", propDef.Tag, propDef.Label)
				return nil, eh.WrapErrorWithStatus(err, errMsg, 400)
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
