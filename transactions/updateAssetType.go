package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

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
		{
			Tag:         "skipAssetEmptyValidation",
			Description: "Do not validate existing assets on the update. Its use should be avoided.",
			DataType:    "boolean",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		assetTypes := req["assetTypes"].([]interface{})
		skipAssetsValidation, ok := req["skipAssetEmptyValidation"].(bool)
		if !ok {
			skipAssetsValidation = false
		}

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
				return nil, errors.NewCCError(fmt.Sprintf("asset type '%s' not found", tagValue.(string)), http.StatusBadRequest)
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
					var emptyAssets bool
					if skipAssetsValidation {
						emptyAssets = true
					} else {
						emptyAssets, err = checkEmptyAssets(stub, tagValue.(string))
						if err != nil {
							return nil, errors.WrapError(err, "failed to check if there assets for tag")
						}
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
			if len(requiredValuesMap) > 0 && !skipAssetsValidation {
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

	for _, p := range propMap {
		p, ok := p.(map[string]interface{})
		if !ok {
			return assetType, nil, errors.NewCCError("invalid prop object", http.StatusBadRequest)
		}

		tag, err := assets.CheckValue(p["tag"], false, "string", "tag")
		if err != nil {
			return assetType, nil, errors.WrapError(err, "invalid tag value")
		}
		tagValue := tag.(string)

		delete, err := assets.CheckValue(p["delete"], false, "boolean", "delete")
		if err != nil {
			return assetType, nil, errors.WrapError(err, "invalid delete info")
		}
		deleteVal := delete.(bool)

		hasProp := assetType.HasProp(tagValue)

		if deleteVal && !hasProp {
			// Handle inexistant prop deletion
			return assetType, nil, errors.NewCCError("attempt to delete inexistent prop", http.StatusBadRequest)
		} else if deleteVal && hasProp {
			// Delete prop
			for i, prop := range propObj {
				if prop.Tag == tagValue {
					if prop.IsKey {
						return assetType, nil, errors.NewCCError("cannot delete key prop", http.StatusBadRequest)
					}
					propObj = append(propObj[:i], propObj[i+1:]...)
				}
			}
		} else if !hasProp && !deleteVal {
			// Create new prop
			required, err := assets.CheckValue(p["required"], false, "boolean", "required")
			if err != nil {
				return assetType, nil, errors.WrapError(err, "invalid required info")
			}
			requiredVal := required.(bool)

			if requiredVal {
				defaultValue, ok := p["defaultValue"]
				if !ok && !emptyAssets {
					return assetType, nil, errors.NewCCError("required prop must have a default value in case of existing assets", http.StatusBadRequest)
				}

				requiredValue := map[string]interface{}{
					"tag":          tagValue,
					"defaultValue": defaultValue,
				}
				requiredValues = append(requiredValues, requiredValue)
			}

			newProp, err := assets.BuildAssetProp(p, nil)
			if err != nil {
				return assetType, nil, errors.WrapError(err, "failed to build prop")
			}

			if newProp.IsKey {
				return assetType, nil, errors.NewCCError("cannot create key prop", http.StatusBadRequest)
			}

			propObj = append(propObj, newProp)
		} else {
			// Update prop
			for i, prop := range propObj {
				if prop.Tag == tagValue {
					required, err := assets.CheckValue(p["required"], false, "boolean", "required")
					if err != nil {
						return assetType, nil, errors.WrapError(err, "invalid required info")
					}
					requiredVal := required.(bool)

					defaultValue, ok := p["defaultValue"]
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

					updatedProp, err := assets.HandlePropUpdate(prop, p)
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

	return !resultsIterator.HasNext(), nil
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

		for _, propMap := range defaultValuesMap {
			propTag := propMap["tag"].(string)
			if _, ok := data[propTag]; !ok {
				data[propTag] = propMap["defaultValue"]
			}
		}

		asset, err := assets.NewAsset(data)
		if err != nil {
			return nil, errors.WrapError(err, "could not assemble asset type")
		}
		assetMap := (map[string]interface{})(asset)

		assetMap, err = asset.Update(stub, assetMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update asset")
		}
		res = append(res, assetMap)
	}

	return res, nil
}
