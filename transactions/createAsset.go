package transactions

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// CreateAsset is the transaction which creates a generic asset
var CreateAsset = Transaction{
	Tag:         "createAsset",
	Label:       "Create Asset",
	Description: "",
	Method:      "POST",

	MetaTx: true,
	Args: []Argument{
		{
			Tag:         "asset",
			Description: "List of assets to be created.",
			DataType:    "[]@asset",
			Required:    true,
		},
	},
	Routine: func(stub shim.ChaincodeStubInterface, req map[string]interface{}) ([]byte, errors.ICCError) {
		// This is safe to do because validation is done before calling routine
		assetList := req["asset"].([]interface{})

		if len(assetList) == 0 {
			return nil, errors.NewCCError("argument 'asset' must be non-empty list", 400)
		}

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
