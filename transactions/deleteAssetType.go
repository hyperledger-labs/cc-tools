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
		res := make([]map[string]interface{}, 0)
		assetTypeList := assets.AssetTypeList()

		for _, assetType := range assetTypes {
			assetTypeMap := assetType.(map[string]interface{})

			tagValue, err := CheckValue(assetTypeMap["tag"], true, "string", "tag")
			if err != nil {
				return nil, errors.WrapError(err, "no tag value in item")
			}

			forceValue, err := CheckValue(assetTypeMap["force"], false, "boolean", "force")
			if err != nil {
				return nil, errors.WrapError(err, "error getting force value")
			}

			forceCascadeValue, err := CheckValue(assetTypeMap["forceCascade"], false, "boolean", "forceCascade")
			if err != nil {
				return nil, errors.WrapError(err, "error getting force value")
			}

			// Verify Asset Type existance
			assetTypeCheck := assets.FetchAssetType(tagValue.(string))
			if assetTypeCheck == nil {
				return nil, errors.WrapError(err, fmt.Sprintf("asset type '%s' not found", tagValue.(string)))
			}

			// Verify Asset Type usage
			err = CheckUsedAssetTypes(stub, tagValue.(string), forceValue.(bool), forceCascadeValue.(bool))
			if err != nil {
				return nil, errors.WrapError(err, "error checking asset type usage")
			}

			// Delete Asset Type
			assetTypeList = assets.RemoveAssetType(tagValue.(string), assetTypeList)
		}

		assets.InitAssetList(assetTypeList)

		resBytes, err := json.Marshal(res)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}

func CheckUsedAssetTypes(stub *sw.StubWrapper, tag string, force, cascade bool) errors.ICCError {
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
		return errors.WrapError(err, "failed to get query result")
	}

	if resultsIterator.HasNext() && !force {
		return errors.NewCCError(fmt.Sprintf("asset type '%s' is in use", tag), http.StatusBadRequest)
	}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return errors.WrapErrorWithStatus(err, "error iterating response", http.StatusInternalServerError)
		}

		var data map[string]interface{}

		err = json.Unmarshal(queryResponse.Value, &data)
		if err != nil {
			return errors.WrapErrorWithStatus(err, "failed to unmarshal queryResponse values", http.StatusInternalServerError)
		}

		asset, err := assets.NewAsset(data)
		if err != nil {
			return errors.WrapError(err, "could not assemble asset type")
		}

		if cascade {
			_, err = asset.DeleteCascade(stub)
		} else {
			_, err = asset.Delete(stub)
		}
		if err != nil {
			return errors.WrapError(err, "could not force delete asset")
		}
	}

	return nil
}
