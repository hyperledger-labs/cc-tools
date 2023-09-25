package transactions

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
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
			Tag:         "assetTypes",
			Description: "Asset Types to be created.",
			DataType:    "[]@object",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		assetTypes := req["assetTypes"].([]interface{})
		list := make([]assets.AssetType, 0)

		for _, assetType := range assetTypes {
			assetTypeMap := assetType.(map[string]interface{})

			newAssetType, err := buildAssetType(assetTypeMap, assetTypes)
			if err != nil {
				return nil, errors.WrapError(err, "failed to build asset type")
			}

			// Verify Asset Type existance
			assetTypeCheck := assets.FetchAssetType(newAssetType.Tag)
			if assetTypeCheck == nil {
				list = append(list, newAssetType)
			}
		}

		if len(list) > 0 {
			assets.UpdateAssetList(list)

			err := assets.StoreAssetList(stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to store asset list")
			}
		}

		resBytes, nerr := json.Marshal(list)
		if nerr != nil {
			return nil, errors.WrapError(nerr, "failed to marshal response")
		}

		return resBytes, nil
	},
}

func buildAssetType(typeMap map[string]interface{}, newTypesList []interface{}) (assets.AssetType, errors.ICCError) {
	// Build Props Array
	propsArr, ok := typeMap["props"].([]interface{})
	if !ok {
		return assets.AssetType{}, errors.NewCCError("invalid props array", http.StatusBadRequest)
	}

	hasKey := false
	props := make([]assets.AssetProp, len(propsArr))
	for i, prop := range propsArr {
		propMap := prop.(map[string]interface{})
		assetProp, err := assets.BuildAssetProp(propMap, newTypesList)
		if err != nil {
			return assets.AssetType{}, errors.WrapError(err, "failed to build asset prop")
		}
		if assetProp.IsKey {
			hasKey = true
		}
		props[i] = assetProp
	}

	if !hasKey {
		return assets.AssetType{}, errors.NewCCError("asset type must have a key", http.StatusBadRequest)
	}

	// Tag
	tagValue, err := assets.CheckValue(typeMap["tag"], true, "string", "tag")
	if err != nil {
		return assets.AssetType{}, errors.WrapError(err, "invalid tag value")
	}

	// Label
	labelValue, err := assets.CheckValue(typeMap["label"], true, "string", "label")
	if err != nil {
		return assets.AssetType{}, errors.WrapError(err, "invalid label value")
	}

	// Description
	descriptionValue, err := assets.CheckValue(typeMap["description"], false, "string", "description")
	if err != nil {
		return assets.AssetType{}, errors.WrapError(err, "invalid description value")
	}

	assetType := assets.AssetType{
		Tag:         tagValue.(string),
		Label:       labelValue.(string),
		Description: descriptionValue.(string),
		Props:       props,
		Dynamic:     true,
	}

	// Readers
	readers := make([]string, 0)
	readersArr, ok := typeMap["readers"].([]interface{})
	if ok {
		for _, reader := range readersArr {
			readerValue, err := assets.CheckValue(reader, false, "string", "reader")
			if err != nil {
				return assets.AssetType{}, errors.WrapError(err, "invalid reader value")
			}

			readers = append(readers, readerValue.(string))
		}
	}
	if len(readers) > 0 {
		assetType.Readers = readers
	}

	return assetType, nil
}
