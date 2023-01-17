package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// DeleteAssetType is the transaction which deletes a dynamic Asset Type
var DeleteAssetType = Transaction{
	Tag:         "deleteAssetType",
	Label:       "Delete Asset Type",
	Description: "",
	Method:      "POST",

	MetaTx: true,
	Args: ArgList{
		{
			Tag:         "assetTypes",
			Description: "Asset Types to be deleted.",
			DataType:    "[]@object",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		assetTypes := req["assetTypes"].([]interface{})

		assetTypeList := assets.AssetTypeList()

		resArr := make([]map[string]interface{}, 0)
		for _, assetType := range assetTypes {
			res := make(map[string]interface{})

			assetTypeMap := assetType.(map[string]interface{})

			tagValue, err := assets.CheckValue(assetTypeMap["tag"], true, "string", "tag")
			if err != nil {
				return nil, errors.WrapError(err, "no tag value in item")
			}

			forceValue, err := assets.CheckValue(assetTypeMap["force"], false, "boolean", "force")
			if err != nil {
				return nil, errors.WrapError(err, "error getting force value")
			}

			// Verify Asset Type existance
			assetTypeCheck := assets.FetchAssetType(tagValue.(string))
			if assetTypeCheck == nil {
				return nil, errors.WrapError(err, fmt.Sprintf("asset type '%s' not found", tagValue.(string)))
			}

			// Verify Asset Type usage
			res["assets"], err = handleRegisteredAssets(stub, tagValue.(string), forceValue.(bool))
			if err != nil {
				return nil, errors.WrapError(err, "error checking asset type usage")
			}

			// Verify Asset Type references
			err = checkAssetTypeReferences(tagValue.(string))
			if err != nil {
				return nil, errors.WrapError(err, "error checking asset type references")
			}

			// Delete Asset Type
			assetTypeList = assets.RemoveAssetType(tagValue.(string), assetTypeList)

			res["assetType"] = assetTypeCheck
			resArr = append(resArr, res)
		}

		assets.ReplaceAssetList(assetTypeList)

		resBytes, err := json.Marshal(resArr)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}

func handleRegisteredAssets(stub *sw.StubWrapper, tag string, force bool) ([]interface{}, errors.ICCError) {
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
		return nil, errors.WrapError(err, "failed to get query result")
	}

	if resultsIterator.HasNext() && !force {
		return nil, errors.NewCCError(fmt.Sprintf("asset type '%s' is in use", tag), http.StatusBadRequest)
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

		res = append(res, asset)

		_, err = asset.Delete(stub)
		if err != nil {
			return nil, errors.WrapError(err, "could not force delete asset")
		}
	}

	return res, nil
}

func checkAssetTypeReferences(tag string) errors.ICCError {
	assetTypeList := assets.AssetTypeList()

	for _, assetType := range assetTypeList {
		subAssets := assetType.SubAssets()
		for _, subAsset := range subAssets {
			if subAsset.Tag == tag {
				return errors.NewCCError(fmt.Sprintf("asset type '%s' is referenced by asset type '%s'", tag, assetType.Tag), http.StatusBadRequest)
			}
		}
	}

	return nil
}
