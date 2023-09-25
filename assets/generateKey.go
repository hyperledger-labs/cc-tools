package assets

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hyperledger-labs/cc-tools/errors"
)

// GenerateKey implements the logic to generate an asset's unique key. It validates
// the assets properties and generates a hash with this values. Based around SHA1 hash function
func GenerateKey(asset map[string]interface{}) (string, errors.ICCError) {
	if key, keyExists := asset["@key"]; keyExists {
		keyStr, ok := key.(string)
		if ok {
			return keyStr, nil
		}
	}

	// Perform validation of the @assetType field
	assetType, exists := asset["@assetType"]
	if !exists {
		return "", errors.NewCCError("property @assetType is required", 400)
	}
	assetTypeString, ok := assetType.(string)
	if !ok {
		return "", errors.NewCCError("property @assetType must be a string", 400)
	}

	// Fetch asset properties
	assetProps := FetchAssetType(assetTypeString)
	if assetProps == nil {
		errMsg := fmt.Sprintf("assetType named '%s' does not exist", assetTypeString)
		return "", errors.NewCCError(errMsg, 400)
	}

	keyProps := assetProps.Keys()

	keySeed := ""
	for _, prop := range keyProps {
		// Check if required key property is included
		propInterface, propIncluded := asset[prop.Tag]
		if !propIncluded {
			errMsg := fmt.Sprintf("primary key %s (%s) is required", prop.Tag, prop.Label)
			return "", errors.NewCCError(errMsg, 400)
		}

		var isArray bool
		var isSubAsset bool
		dataTypeName := prop.DataType

		if strings.HasPrefix(dataTypeName, "[]") {
			dataTypeName = strings.TrimPrefix(dataTypeName, "[]")
			isArray = true
		}

		if strings.HasPrefix(dataTypeName, "->") {
			dataTypeName = strings.TrimPrefix(dataTypeName, "->")
			isSubAsset = true
		}

		// Handle array-like asset property types
		var propAsArray []interface{}
		if !isArray {
			propAsArray = []interface{}{propInterface}
		} else {
			propAsArray, ok = propInterface.([]interface{})
			if !ok {
				return "", errors.NewCCError(fmt.Sprintf("asset property %s must an array of type %s", prop.Label, prop.DataType), 400)
			}
		}

		// Iterate asset properties to form keySeed
		for _, propInterface := range propAsArray {
			if !isSubAsset {
				// If key is a primitive data type, append its String value to seed
				dataType, dataTypeExists := dataTypeMap[dataTypeName]
				if !dataTypeExists {
					return "", errors.NewCCError(fmt.Sprintf("internal error: invalid prop data type %s", prop.DataType), 500)
				}
				var seed string
				var err error

				seed, _, err = dataType.Parse(propInterface)
				if err != nil {
					return "", errors.WrapError(err, fmt.Sprintf("failed to generate key for asset property '%s'", prop.Label))
				}

				keySeed += seed
			} else {
				// If key is a subAsset, generate subAsset's key to append to seed
				assetTypeDef := FetchAssetType(dataTypeName)
				if assetTypeDef == nil {
					return "", errors.NewCCError(fmt.Sprintf("internal error: invalid sub asset type %s", prop.DataType), 500)
				}

				var propMap map[string]interface{}
				switch t := propInterface.(type) {
				case map[string]interface{}:
					propMap = t
				case Key:
					propMap = t
				case Asset:
					propMap = t
				default:
					errMsg := fmt.Sprintf("subAsset key %s must be sent as map[string]interface{} (JSON object)", prop.Tag)
					return "", errors.NewCCError(errMsg, 400)
				}

				propMap["@assetType"] = dataTypeName
				subAssetKey, err := GenerateKey(propMap)
				if err != nil {
					errMsg := fmt.Sprintf("error generating key for subAsset key '%s'", prop.Tag)
					return "", errors.WrapError(err, errMsg)
				}

				keySeed += subAssetKey
			}
		}
	}

	key := assetTypeString + ":" + uuid.NewSHA1(uuid.NameSpaceOID, []byte(keySeed)).String()

	return key, nil
}
