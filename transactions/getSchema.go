package transactions

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// GetSchema returns information about a specific AssetType or a list of every configured AssetType
var GetSchema = Transaction{
	Tag:         "getSchema",
	Label:       "Get Schema",
	Description: "",
	Method:      "GET",

	ReadOnly: true,
	MetaTx:   true,
	Args: ArgList{
		{
			Tag:         "assetType",
			DataType:    "string",
			Description: "The name of the asset type of which you want to fetch the definition. Leave empty to fetch a list of possible types.",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var assetTypeName string
		assetTypeInterface, ok := req["assetType"]
		if ok {
			assetTypeName, ok = assetTypeInterface.(string)
			if !ok {
				return nil, errors.NewCCError("argument 'assetType' must be a string", 400)
			}
		}

		// If user requested a specific asset type definition
		if assetTypeName != "" {
			assetTypeDef := assets.FetchAssetType(assetTypeName)
			if assetTypeDef == nil {
				errMsg := fmt.Sprintf("asset type named %s does not exist", assetTypeName)
				return nil, errors.NewCCError(errMsg, 404)
			}
			assetDefBytes, err := json.Marshal(assetTypeDef)
			if err != nil {
				errMsg := fmt.Sprintf("error marshaling asset definition: %s", err)
				return nil, errors.NewCCError(errMsg, 500)
			}
			return assetDefBytes, nil
		}

		assetTypeList := assets.AssetTypeList()
		// If user requested asset list
		type assetListElem struct {
			Tag         string   `json:"tag"`
			Label       string   `json:"label"`
			Description string   `json:"description"`
			Readers     []string `json:"readers,omitempty"`
			Writers     []string `json:"writers"`
			Dynamic     bool     `json:"dynamic"`
		}
		var assetList []assetListElem
		for _, assetTypeDef := range assetTypeList {
			assetList = append(assetList, assetListElem{
				Tag:         assetTypeDef.Tag,
				Label:       assetTypeDef.Label,
				Description: assetTypeDef.Description,
				Readers:     assetTypeDef.Readers,
				Dynamic:     assetTypeDef.Dynamic,
			})
		}

		assetListBytes, err := json.Marshal(assetList)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling asset list", 500)
		}
		return assetListBytes, nil
	},
}
