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
// ? Allow isKey prop to be updated/created/removed?

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
		assetTypeListFallback := assets.AssetTypeList()

		resAssetArr := make([]assets.AssetType, 0)
		requiredValues := make(map[string]interface{}, 0)

		for _, assetType := range assetTypes {
			assetTypeMap := assetType.(map[string]interface{})

			tagValue, err := assets.CheckValue(assetTypeMap["tag"], true, "string", "tag")
			if err != nil {
				return nil, errors.WrapError(err, "no tag value in item")
			}

			// Verify Asset Type existance
			assetTypeCheck := assets.FetchAssetType(tagValue.(string))
			if assetTypeCheck == nil {
				return nil, errors.WrapError(err, fmt.Sprintf("asset type '%s' not found", tagValue.(string)))
			}
			assetTypeObj := *assetTypeCheck

			// Verify if Asset Type allows dynamic modifications
			if !assetTypeObj.Dynamic {
				return nil, errors.WrapError(err, fmt.Sprintf("asset type '%s' does not allows dynamic modifications", tagValue.(string)))
			}

			for key, value := range assetTypeMap {
				switch key {
				case "label":
					labelValue, err := assets.CheckValue(value, true, "string", "label")
					if err != nil {
						return nil, errors.WrapError(err, "invalid label value")
					}
					assetTypeObj.Label = labelValue.(string)
				case "description":
					descriptionValue, err := assets.CheckValue(value, true, "string", "description")
					if err != nil {
						return nil, errors.WrapError(err, "invalid description value")
					}
					assetTypeObj.Description = descriptionValue.(string)
				case "readers":
					readers := make([]string, 0)
					readersArr, ok := value.([]interface{})
					if ok {
						for _, reader := range readersArr {
							readerValue, err := assets.CheckValue(reader, false, "string", "reader")
							if err != nil {
								return nil, errors.WrapError(err, "invalid reader value")
							}

							readers = append(readers, readerValue.(string))
						}
						assetTypeObj.Readers = readers
					}
				case "props":
					propsArr, ok := value.([]interface{})
					if !ok {
						return nil, errors.NewCCError("invalid props array", http.StatusBadRequest)
					}
					emptyAssets, err := checkEmptyAssets(stub, tagValue.(string))
					if err != nil {
						return nil, errors.WrapError(err, "failed to check if there assets for tag")
					}
					newAssetType, newRequiredValues, err := handleProps(assetTypeObj, propsArr, emptyAssets)
					if err != nil {
						return nil, errors.WrapError(err, "invalid props array")
					}
					requiredValues[tagValue.(string)] = newRequiredValues
					assetTypeObj = newAssetType
				default:
					continue
				}
			}

			// Update Asset Type
			assets.ReplaceAssetType(assetTypeObj, assetTypeList)
			resAssetArr = append(resAssetArr, assetTypeObj)
		}

		response := map[string]interface{}{
			"assetTypes": resAssetArr,
		}

		assets.ReplaceAssetList(assetTypeList)

		for k, v := range requiredValues {
			requiredValuesMap := v.([]map[string]interface{})
			if len(requiredValuesMap) > 0 {
				updatedAssets, err := initilizeDefaultValues(stub, k, requiredValuesMap)
				if err != nil {
					// Rollback Asset Type List
					assets.ReplaceAssetList(assetTypeListFallback)
					return nil, errors.WrapError(err, "failed to initialize default values")
				}
				response["assets"] = updatedAssets
			}
		}

		err := assets.StoreAssetList(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to store asset list")
		}

		err = assets.SetEventForList(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to set event for asset list")
		}

		resBytes, nerr := json.Marshal(response)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}

func handleProps(assetType assets.AssetType, propMap []interface{}, emptyAssets bool) (assets.AssetType, []map[string]interface{}, errors.ICCError) {
	propObj := assetType.Props
	requiredValues := make([]map[string]interface{}, 0)

	for _, v := range propMap {
		v, ok := v.(map[string]interface{})
		if !ok {
			return assetType, nil, errors.NewCCError("invalid prop object", http.StatusBadRequest)
		}

		// Get "tag" and "delete" values
		tag, err := assets.CheckValue(v["tag"], false, "string", "tag")
		if err != nil {
			return assetType, nil, errors.WrapError(err, "invalid tag value")
		}
		tagValue := tag.(string)

		delete, err := assets.CheckValue(v["delete"], false, "boolean", "delete")
		if err != nil {
			return assetType, nil, errors.WrapError(err, "invalid delete info")
		}
		deleteVal := delete.(bool)

		hasProp := assetType.HasProp(tagValue)

		if deleteVal && !hasProp {
			// Deleting nexistant prop
			return assetType, nil, errors.WrapError(err, "attempt to delete inexistent prop")
		} else if deleteVal && hasProp {
			// Prop deletion
			for i, prop := range propObj {
				if prop.Tag == tagValue {
					if prop.IsKey {
						return assetType, nil, errors.NewCCError("cannot delete key prop", http.StatusBadRequest)
					}
					propObj = append(propObj[:i], propObj[i+1:]...)
				}
			}
		} else if !hasProp && !deleteVal {
			// Prop creation
			required, err := assets.CheckValue(v["required"], false, "boolean", "required")
			if err != nil {
				return assetType, nil, errors.WrapError(err, "invalid required info")
			}
			requiredVal := required.(bool)

			if requiredVal {
				defaultValue, ok := v["defaultValue"]
				if !ok && !emptyAssets {
					return assetType, nil, errors.NewCCError("required prop must have a default value in case of existing assets", http.StatusBadRequest)
				}

				requiredValue := map[string]interface{}{
					"tag":          tagValue,
					"defaultValue": defaultValue,
				}
				requiredValues = append(requiredValues, requiredValue)
			}

			newProp, err := assets.BuildAssetProp(v)
			if err != nil {
				return assetType, nil, errors.WrapError(err, "failed to build prop")
			}

			propObj = append(propObj, newProp)
		} else {
			// Prop update
			for i, prop := range propObj {
				if prop.Tag == tagValue {
					required, err := assets.CheckValue(v["required"], false, "boolean", "required")
					if err != nil {
						return assetType, nil, errors.WrapError(err, "invalid required info")
					}
					requiredVal := required.(bool)

					defaultValue, ok := v["defaultValue"]
					if !ok {
						defaultValue = prop.DefaultValue
					}

					if !prop.Required && requiredVal {
						if defaultValue == nil && !emptyAssets {
							return assetType, nil, errors.NewCCError("required prop must have a default value in case of existing assets", http.StatusBadRequest)
						}

						requiredValue := map[string]interface{}{
							"tag":          tagValue,
							"defaultValue": defaultValue,
						}
						requiredValues = append(requiredValues, requiredValue)
					}

					updatedProp, err := assets.HandlePropUpdate(prop, v)
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

func checkEmptyAssets(stub *sw.StubWrapper, tag string) (bool, errors.ICCError) {
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
		return true, nil
	}

	return false, nil
}

func initilizeDefaultValues(stub *sw.StubWrapper, assetTag string, defaultValuesMap []map[string]interface{}) ([]interface{}, errors.ICCError) {
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
