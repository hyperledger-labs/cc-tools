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
			DataType:    "[]@object",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		assetTypes := req["assetType"].([]interface{})
		list := make([]assets.AssetType, 0)

		for _, assetType := range assetTypes {
			assetTypeMap := assetType.(map[string]interface{})

			props := make([]assets.AssetProp, len(assetTypeMap["props"].([]interface{})))
			for i, prop := range assetTypeMap["props"].([]interface{}) {
				propMap := prop.(map[string]interface{})

				writersArray := propMap["writers"].([]interface{})
				writers := make([]string, len(writersArray))
				for j, writer := range writersArray {
					writers[j] = writer.(string)
				}

				props[i] = assets.AssetProp{
					Tag:      propMap["tag"].(string),
					Label:    propMap["label"].(string),
					Required: propMap["required"].(bool),
					DataType: propMap["dataType"].(string),
					IsKey:    propMap["isKey"].(bool),
					Writers:  writers,
				}
			}

			var newAssetType = assets.AssetType{
				Tag:         assetTypeMap["tag"].(string),
				Label:       assetTypeMap["label"].(string),
				Description: assetTypeMap["description"].(string),
				Props:       props,
			}

			list = append(list, newAssetType)
		}

		// TODO: Check if asset type already exists

		assets.UpdateAssetList(list)

		resBytes, err := json.Marshal(list)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}
