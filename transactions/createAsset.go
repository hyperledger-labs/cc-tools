package transactions

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// CreateAsset is the transaction which creates a generic asset
var CreateAsset = Transaction{
	Tag:         "createAsset",
	Label:       "Create Asset",
	Description: "",
	Method:      "POST",

	MetaTx: true,
	Args: ArgList{
		{
			Tag:         "asset",
			Description: "List of assets to be created.",
			DataType:    "[]@asset",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		// This is safe to do because validation is done before calling routine
		assetList := req["asset"].([]interface{})

		responses := []map[string]interface{}{}
		for _, assetInterface := range assetList {
			// This is safe to do because validation is done before calling routine
			asset := assetInterface.(assets.Asset)

			// Marshal asset back to JSON format
			res, err := asset.PutNew(stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to write asset to ledger")
			}

			responses = append(responses, res)
		}

		resBytes, err := json.Marshal(responses)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}
