package transactions

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// CreateAssetType is the transaction which creates a dynamic Asset Type
var CreateAssetType = Transaction{
	Tag:         "createAssetType",
	Label:       "Create Asset Type",
	Description: "",
	Method:      "POST",

	MetaTx: true,
	Args: ArgList{
		{
			Tag:         "assetType",
			Description: "Asset Type to be created.",
			DataType:    "@object",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		// This is safe to do because validation is done before calling routine
		assetType := req["assetType"].(map[string]interface{})

		props := make([]assets.AssetProp, len(assetType["props"].([]interface{})))
		for i, prop := range assetType["props"].([]interface{}) {
			props[i] = assets.AssetProp{
				Tag:      prop.(map[string]interface{})["tag"].(string),
				Label:    prop.(map[string]interface{})["label"].(string),
				Required: prop.(map[string]interface{})["required"].(bool),
				DataType: prop.(map[string]interface{})["dataType"].(string),
				IsKey:    prop.(map[string]interface{})["isKey"].(bool),
				Writers:  prop.(map[string]interface{})["writers"].([]string),
			}
		}

		var newAssetType = assets.AssetType{
			Tag:         assetType["tag"].(string),
			Label:       assetType["label"].(string),
			Description: assetType["description"].(string),
			Props:       props,
		}

		// TODO: Check if asset type already exists

		// Add asset type to assetTypeList
		list := make([]assets.AssetType, 0)
		list = append(list, newAssetType)

		assets.UpdateAssetList(list)

		resBytes, err := json.Marshal(list)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}
