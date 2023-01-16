package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// ? Update tag name?
// TODO: Handle not required -> required

// UpdateAssetType is the transaction which updates a dynamic Asset Type
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

		assetTypeList := assets.AssetTypeList()

		resArr := make([]map[string]interface{}, 0)
		requiredValues := make(map[string]interface{}, 0)

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
					newAssetType, newRequiredValues, err := handleProps(assetTypeObj, propsArr)
					if err != nil {
						return nil, errors.WrapError(err, "invalid props array")
					}
					requiredValues[tagValue.(string)] = newRequiredValues
					assetTypeObj = newAssetType
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

		for k, v := range requiredValues {
			requiredValuesMap := v.([]map[string]interface{})
			if len(requiredValuesMap) > 0 {
				InitilizeDefaultValues(stub, k, requiredValuesMap)
			}
		}

		resBytes, err := json.Marshal(resArr)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}

func handleProps(assetType assets.AssetType, propMap []interface{}) (assets.AssetType, []map[string]interface{}, errors.ICCError) {
	propObj := assetType.Props
	requiredValues := make([]map[string]interface{}, 0)

	for _, v := range propMap {
		v, ok := v.(map[string]interface{})
		if !ok {
			return assetType, nil, errors.NewCCError("invalid prop object", http.StatusBadRequest)
		}

		tag, err := CheckValue(v["tag"], false, "string", "tag")
		if err != nil {
			return assetType, nil, errors.WrapError(err, "invalid tag value")
		}
		tagValue := tag.(string)

		hasProp := assetType.HasProp(tagValue)

		delete, err := CheckValue(v["delete"], false, "boolean", "delete")
		if err != nil {
			return assetType, nil, errors.WrapError(err, "invalid delete info")
		}
		deleteVal := delete.(bool)

		if deleteVal && !hasProp {
			return assetType, nil, errors.WrapError(err, "attempt to delete inexistent prop")
		} else if deleteVal && hasProp {
			// ? Should you be able to delete a required prop?
			for i, prop := range propObj {
				if prop.Tag == tagValue {
					if prop.IsKey {
						return assetType, nil, errors.NewCCError("cannot delete key prop", http.StatusBadRequest)
					}
					propObj = append(propObj[:i], propObj[i+1:]...)
				}
			}
		} else if !hasProp && !deleteVal {
			// ? Should you be able to create a isKey prop?
			// TODO: Handle verification if assets exists on require
			required, err := CheckValue(v, false, "boolean", "required")
			if err != nil {
				return assetType, nil, errors.WrapError(err, "invalid required info")
			}
			requiredVal := required.(bool)
			if requiredVal {
				defaultValue, ok := v["defaultValue"]
				if !ok {
					return assetType, nil, errors.NewCCError("required prop must have a default value in case of existing assets", http.StatusBadRequest)
				}

				requiredValue := map[string]interface{}{
					"tag":          tagValue,
					"defaultValue": defaultValue,
				}
				requiredValues = append(requiredValues, requiredValue)
			}
			newProp, err := BuildAssetProp(v)
			if err != nil {
				return assetType, nil, errors.WrapError(err, "failed to build prop")
			}
			propObj = append(propObj, newProp)
		} else {
			// TODO: Handle required/isKey prop
			for i, prop := range propObj {
				if prop.Tag == tagValue {
					updatedProp, err := handlePropUpdate(prop, v)
					if err != nil {
						return assetType, nil, errors.WrapError(err, "failed to update prop")
					}
					propObj[i] = updatedProp
				}
			}
		}
	}

	assetType.Props = propObj
	return assetType, requiredValues, nil
}

func handlePropUpdate(assetProps assets.AssetProp, propMap map[string]interface{}) (assets.AssetProp, errors.ICCError) {
	// ? Update tag?
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

func CheckExistingAssets(stub *sw.StubWrapper, tag string) (bool, errors.ICCError) {
	query := fmt.Sprintf(
		`{
			"selector": {
			   "@assetType": "%s"
			}
		}`,
		tag,
	)

	resultsIterator, err := stub.GetQueryResult(query)
	if err != nil {
		return false, errors.WrapError(err, "failed to get query result")
	}

	if resultsIterator.HasNext() {
		return false, nil
	}

	return true, nil
}

func InitilizeDefaultValues(stub *sw.StubWrapper, assetTag string, defaultValuesMap []map[string]interface{}) ([]interface{}, errors.ICCError) {
	query := fmt.Sprintf(
		`{
			"selector": {
			   "@assetType": "%s"
			}
		}`,
		assetTag,
	)

	resultsIterator, err := stub.GetQueryResult(query)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get query result")
	}

	res := make([]interface{}, 0)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error iterating response", http.StatusInternalServerError)
		}

		var data map[string]interface{}

		err = json.Unmarshal(queryResponse.Value, &data)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal queryResponse values", http.StatusInternalServerError)
		}

		asset, err := assets.NewAsset(data)
		if err != nil {
			return nil, errors.WrapError(err, "could not assemble asset type")
		}
		assetMap := (map[string]interface{})(asset)

		for _, propMap := range defaultValuesMap {
			propTag := propMap["tag"].(string)
			if _, ok := assetMap[propTag]; !ok {
				assetMap[propTag] = propMap["defaultValue"]
			}

		}

		assetMap, err = asset.Update(stub, assetMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update asset")
		}
		res = append(res, assetMap)
	}

	return res, nil
}
