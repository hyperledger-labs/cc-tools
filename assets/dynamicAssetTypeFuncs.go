package assets

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyperledger-labs/cc-tools/errors"
)

// BuildAssetProp builds an AssetProp from an object with the required fields
func BuildAssetProp(propMap map[string]interface{}, newTypesList []interface{}) (AssetProp, errors.ICCError) {
	// Tag
	tagValue, err := CheckValue(propMap["tag"], true, "string", "tag")
	if err != nil {
		return AssetProp{}, errors.WrapError(err, "invalid tag value")
	}

	// Label
	labelValue, err := CheckValue(propMap["label"], true, "string", "label")
	if err != nil {
		return AssetProp{}, errors.WrapError(err, "invalid label value")
	}

	// Description
	descriptionValue, err := CheckValue(propMap["description"], false, "string", "description")
	if err != nil {
		return AssetProp{}, errors.WrapError(err, "invalid description value")
	}

	// Required
	requiredValue, err := CheckValue(propMap["required"], false, "boolean", "required")
	if err != nil {
		return AssetProp{}, errors.WrapError(err, "invalid required value")
	}

	// IsKey
	isKeyValue, err := CheckValue(propMap["isKey"], false, "boolean", "isKey")
	if err != nil {
		return AssetProp{}, errors.WrapError(err, "invalid isKey value")
	}

	// ReadOnly
	readOnlyValue, err := CheckValue(propMap["readOnly"], false, "boolean", "readOnly")
	if err != nil {
		return AssetProp{}, errors.WrapError(err, "invalid readOnly value")
	}

	// DataType
	dataTypeValue, err := CheckValue(propMap["dataType"], true, "string", "dataType")
	if err != nil {
		return AssetProp{}, errors.WrapError(err, "invalid dataType value")
	}
	err = CheckDataType(dataTypeValue.(string), newTypesList)
	if err != nil {
		return AssetProp{}, errors.WrapError(err, "failed checking data type")
	}

	assetProp := AssetProp{
		Tag:         tagValue.(string),
		Label:       labelValue.(string),
		Description: descriptionValue.(string),
		Required:    requiredValue.(bool),
		IsKey:       isKeyValue.(bool),
		ReadOnly:    readOnlyValue.(bool),
		DataType:    dataTypeValue.(string),
	}

	// Writers
	writers := make([]string, 0)
	writersArr, ok := propMap["writers"].([]interface{})
	if ok {
		for _, writer := range writersArr {
			writerValue, err := CheckValue(writer, false, "string", "writer")
			if err != nil {
				return AssetProp{}, errors.WrapError(err, "invalid writer value")
			}

			writers = append(writers, writerValue.(string))
		}
	}
	if len(writers) > 0 {
		assetProp.Writers = writers
	}

	// Validate Default Value
	if propMap["defaultValue"] != nil {
		defaultValue, err := validateProp(propMap["defaultValue"], assetProp)
		if err != nil {
			return AssetProp{}, errors.WrapError(err, "invalid Default Value")
		}

		assetProp.DefaultValue = defaultValue
	}

	return assetProp, nil
}

// HandlePropUpdate updates an AssetProp with the values of the propMap
func HandlePropUpdate(assetProps AssetProp, propMap map[string]interface{}) (AssetProp, errors.ICCError) {
	handleDefaultValue := false
	for k, v := range propMap {
		switch k {
		case "defaultValue":
			handleDefaultValue = true
		case "label":
			labelValue, err := CheckValue(v, true, "string", "label")
			if err != nil {
				return assetProps, errors.WrapError(err, "invalid label value")
			}
			assetProps.Label = labelValue.(string)
		case "description":
			descriptionValue, err := CheckValue(v, true, "string", "description")
			if err != nil {
				return assetProps, errors.WrapError(err, "invalid description value")
			}
			assetProps.Description = descriptionValue.(string)
		case "required":
			requiredValue, err := CheckValue(v, true, "boolean", "required")
			if err != nil {
				return assetProps, errors.WrapError(err, "invalid required value")
			}
			assetProps.Required = requiredValue.(bool)
		case "readOnly":
			readOnlyValue, err := CheckValue(v, true, "boolean", "readOnly")
			if err != nil {
				return assetProps, errors.WrapError(err, "invalid readOnly value")
			}
			assetProps.ReadOnly = readOnlyValue.(bool)
		case "writers":
			writers := make([]string, 0)
			writersArr, ok := v.([]interface{})
			if ok {
				for _, writer := range writersArr {
					writerValue, err := CheckValue(writer, false, "string", "writer")
					if err != nil {
						return AssetProp{}, errors.WrapError(err, "invalid writer value")
					}

					writers = append(writers, writerValue.(string))
				}
			}
			assetProps.Writers = writers
		default:
			continue
		}
	}

	if handleDefaultValue {
		defaultValue, err := validateProp(propMap["defaultValue"], assetProps)
		if err != nil {
			return AssetProp{}, errors.WrapError(err, "invalid Default Value")
		}
		assetProps.DefaultValue = defaultValue
	}

	return assetProps, nil
}

// CheckDataType verifies if dataType is valid among the ones availiable in the chaincode
func CheckDataType(dataType string, newTypesList []interface{}) errors.ICCError {
	trimDataType := strings.TrimPrefix(dataType, "[]")

	if strings.HasPrefix(trimDataType, "->") {
		trimDataType = strings.TrimPrefix(trimDataType, "->")

		assetType := FetchAssetType(trimDataType)
		if assetType == nil {
			foundDataType := false
			for _, newTypeInterface := range newTypesList {
				newType := newTypeInterface.(map[string]interface{})
				if newType["tag"] == trimDataType {
					foundDataType = true
					break
				}
			}
			if !foundDataType {
				return errors.NewCCError(fmt.Sprintf("invalid dataType value '%s'", dataType), http.StatusBadRequest)
			}
		}
	} else {
		dataTypeObj := FetchDataType(trimDataType)
		if dataTypeObj == nil {
			return errors.NewCCError(fmt.Sprintf("invalid dataType value '%s'", dataType), http.StatusBadRequest)
		}
	}

	return nil
}

// CheckValue verifies if parameter value is of the expected type
func CheckValue(value interface{}, required bool, expectedType, fieldName string) (interface{}, errors.ICCError) {
	if value == nil {
		if required {
			return nil, errors.NewCCError(fmt.Sprintf("required value %s missing", fieldName), http.StatusBadRequest)
		}
		switch expectedType {
		case "string":
			return "", nil
		case "number":
			return 0, nil
		case "boolean":
			return false, nil
		}
	}

	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return nil, errors.NewCCError(fmt.Sprintf("value %s is not a string", fieldName), http.StatusBadRequest)
		}
	case "number":
		if _, ok := value.(float64); !ok {
			return nil, errors.NewCCError(fmt.Sprintf("value %s is not a number", fieldName), http.StatusBadRequest)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return nil, errors.NewCCError(fmt.Sprintf("value %s is not a boolean", fieldName), http.StatusBadRequest)
		}
	}

	return value, nil
}
