package assets

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/google/uuid"
)

// GenerateKey implements the logic to generate an asset's unique key
func GenerateKey(asset map[string]interface{}) (string, errors.ICCError) {
	if key, keyExists := asset["@key"]; keyExists {
		keyStr, ok := key.(string)
		if ok {
			return keyStr, nil
		}
	}

	var err error

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
		dataType := prop.DataType

		if strings.HasPrefix(dataType, "[]") {
			dataType = strings.TrimPrefix(dataType, "[]")
			isArray = true
		}

		// Handle array-like asset property types
		var propAsArray []interface{}
		if !isArray {
			propAsArray = []interface{}{propInterface}
		} else {
			propAsArray, ok = propInterface.([]interface{})
			if !ok {
				return "", errors.NewCCError(fmt.Sprintf("asset property %s must and array of type %s", prop.Label, prop.DataType), 400)
			}
		}

		// Iterate asset properties to form keySeed
		for _, propInterface := range propAsArray {
			// If key is a subAsset, generate subAsset's key to append to seed
			assetTypeDef := FetchAssetType(dataType)
			if assetTypeDef != nil {
				propMap, ok := propInterface.(map[string]interface{})
				if !ok {
					errMsg := fmt.Sprintf("subAsset key %s must be sent as JSON object", prop.Tag)
					return "", errors.NewCCError(errMsg, 400)
				}
				propMap["@assetType"] = dataType
				subAssetKey, err := GenerateKey(propMap)
				if err != nil {
					errMsg := fmt.Sprintf("error generating key for subAsset key '%s'", prop.Tag)
					return "", errors.WrapError(err, errMsg)
				}
				keySeed += subAssetKey
			} else {
				// If key is a primitive data type, append its raw value to seed
				switch dataType {
				case "string":
					propVal, ok := propInterface.(string)
					if !ok {
						return "", errors.NewCCError(fmt.Sprintf("asset property %s should be of type %s", prop.Label, prop.DataType), 400)
					}
					keySeed += propVal
				case "number":
					propVal, ok := propInterface.(float64)
					if !ok {
						propValStr, okStr := propInterface.(string)
						if !okStr {
							return "", errors.NewCCError(fmt.Sprintf("asset property %s should be of type %s", prop.Label, prop.DataType), 400)
						}
						propVal, err = strconv.ParseFloat(propValStr, 64)
						if err != nil {
							return "", errors.NewCCError(fmt.Sprintf("asset property %s should be of type %s", prop.Label, prop.DataType), 400)
						}
					}
					keySeed += strconv.FormatUint(math.Float64bits(propVal), 16) // Float IEEE 754 hexadecimal representation
				case "boolean":
					propVal, ok := propInterface.(bool)
					if !ok {
						propValStr, okStr := propInterface.(string)
						if !okStr {
							return "", errors.NewCCError(fmt.Sprintf("asset property %s should be of type %s", prop.Label, prop.DataType), 400)
						}
						if propValStr != "true" && propValStr != "false" {
							return "", errors.NewCCError(fmt.Sprintf("asset property %s should be of type %s", prop.Label, prop.DataType), 400)
						}
						if propValStr == "true" {
							propVal = true
						}
					}
					if propVal {
						keySeed += "t"
					} else {
						keySeed += "f"
					}
				case "datetime":
					propVal, ok := propInterface.(string)
					if !ok {
						return "", errors.NewCCError(fmt.Sprintf("asset property %s should be a RFC3339 string", prop.Label), 400)
					}
					propTime, err := time.Parse(time.RFC3339, propVal)
					if err != nil {
						return "", errors.WrapErrorWithStatus(err, fmt.Sprintf("invalid asset property %s RFC3339 format", prop.Label), 400)
					}
					keySeed += propTime.Format(time.RFC3339)
				default:
					return "", errors.NewCCError(fmt.Sprintf("internal error: invalid prop data type %s", prop.DataType), 500)
				}
			}
		}
	}

	key := assetTypeString + ":" + uuid.NewSHA1(uuid.NameSpaceOID, []byte(keySeed)).String()

	return key, nil
}
