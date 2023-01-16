package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// TODO: Update tag name?
// TODO: Handle not required -> required

// CreateAssetType is the transaction which creates a dynamic Asset Type
var UpdateAssetType = Transaction{
	Tag:         "updateAssetType",
	Label:       "Update Asset Type",
	Description: "",
	Method:      "POST",

	MetaTx: true,
	Args: ArgList{
		{
			Tag:         "assetTypes",
			Description: "Asset Types to be updated.",
			DataType:    "[]@object",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		assetTypes := req["assetTypes"].([]interface{})
		resArr := make([]map[string]interface{}, 0)
		assetTypeList := assets.AssetTypeList()

		for _, assetType := range assetTypes {
			assetTypeMap := assetType.(map[string]interface{})

			tagValue, err := CheckValue(assetTypeMap["tag"], true, "string", "tag")
			if err != nil {
				return nil, errors.WrapError(err, "no tag value in item")
			}

			// Verify Asset Type existance
			assetTypeCheck := assets.FetchAssetType(tagValue.(string))
			if assetTypeCheck == nil {
				return nil, errors.WrapError(err, fmt.Sprintf("asset type '%s' not found", tagValue.(string)))
			}
			assetTypeObj := *assetTypeCheck

			for key, value := range assetTypeMap {
				switch key {
				case "props":
					propsArr, ok := value.([]interface{})
					if !ok {
						return nil, errors.NewCCError("invalid props array", http.StatusBadRequest)
					}
					assetTypeObj, err = handleProps(assetTypeObj, propsArr)
					if err != nil {
						return nil, errors.WrapError(err, "invalid props array")
					}
				case "label":
					labelValue, err := CheckValue(value, true, "string", "label")
					if err != nil {
						return nil, errors.WrapError(err, "invalid label value")
					}
					assetTypeObj.Label = labelValue.(string)
				case "description":
					descriptionValue, err := CheckValue(value, true, "string", "description")
					if err != nil {
						return nil, errors.WrapError(err, "invalid description value")
					}
					assetTypeObj.Description = descriptionValue.(string)
				case "readers":
					readers := make([]string, 0)
					readersArr, ok := value.([]interface{})
					if ok {
						for _, reader := range readersArr {
							readerValue, err := CheckValue(reader, false, "string", "reader")
							if err != nil {
								return nil, errors.WrapError(err, "invalid reader value")
							}

							readers = append(readers, readerValue.(string))
						}
						assetTypeObj.Readers = readers
					}
				default:
					continue
				}
			}

			// Update Asset Type
			assets.ReplaceAssetType(assetTypeObj, assetTypeList)
		}

		assets.InitAssetList(assetTypeList)

		resBytes, err := json.Marshal(resArr)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}

func handleProps(assetType assets.AssetType, propMap []interface{}) (assets.AssetType, errors.ICCError) {
	propObj := assetType.Props

	for _, v := range propMap {
		v, ok := v.(map[string]interface{})
		if !ok {
			return assetType, errors.NewCCError("invalid prop object", http.StatusBadRequest)
		}

		hasProp := assetType.HasProp(v["tag"].(string))

		delete, err := CheckValue(v["delete"], false, "boolean", "delete")
		if err != nil {
			return assetType, errors.WrapError(err, "invalid delete info")
		}
		deleteVal := delete.(bool)

		if deleteVal && !hasProp {
			return assetType, errors.WrapError(err, "attempt to delete inexistent prop")
		} else if deleteVal && hasProp {
			// TODO: Handle required/isKey
			for i, prop := range propObj {
				if prop.Tag == v["tag"].(string) {
					propObj = append(propObj[:i], propObj[i+1:]...)
				}
			}
		} else if !hasProp && !deleteVal {
			// TODO: Handle required/isKey prop
			newProp, err := BuildAssetProp(v)
			if err != nil {
				return assetType, errors.WrapError(err, "failed to build prop")
			}
			propObj = append(propObj, newProp)
		} else {
			// TODO: Handle required/isKey prop
			for i, prop := range propObj {
				if prop.Tag == v["tag"].(string) {
					updatedProp, err := handlePropUpdate(prop, v)
					if err != nil {
						return assetType, errors.WrapError(err, "failed to update prop")
					}
					propObj[i] = updatedProp
				}
			}
		}
	}

	assetType.Props = propObj
	return assetType, nil
}

func handlePropUpdate(assetProps assets.AssetProp, propMap map[string]interface{}) (assets.AssetProp, errors.ICCError) {
	// TODO: Update tag?
	handleDefaultValue := false
	for k, v := range propMap {
		switch k {
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
		case "isKey":
			// TODO: Allow isKey to be updated?
			isKeyValue, err := CheckValue(v, true, "boolean", "isKey")
			if err != nil {
				return assetProps, errors.WrapError(err, "invalid isKey value")
			}
			assetProps.IsKey = isKeyValue.(bool)
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
		case "defaultValue":
			handleDefaultValue = true
		case "dataType":
			dataTypeValue, err := CheckValue(propMap["dataType"], true, "string", "dataType")
			if err != nil {
				return assets.AssetProp{}, errors.WrapError(err, "invalid dataType value")
			}

			err = CheckDataType(dataTypeValue.(string))
			if err != nil {
				return assets.AssetProp{}, errors.WrapError(err, "failed checking data type")
			}
			assetProps.DataType = dataTypeValue.(string)
		case "writeres":
			writers := make([]string, 0)
			writersArr, ok := v.([]interface{})
			if ok {
				for _, writer := range writersArr {
					writerValue, err := CheckValue(writer, false, "string", "writer")
					if err != nil {
						return assets.AssetProp{}, errors.WrapError(err, "invalid writer value")
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
		defaultValue, err := assets.ValidateProp(propMap["defaultValue"], assetProps)
		if err != nil {
			return assets.AssetProp{}, errors.WrapError(err, "invalid Default Value")
		}

		assetProps.DefaultValue = defaultValue
	}

	return assetProps, nil
}
