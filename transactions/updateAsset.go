package transactions

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// UpdateAsset is the function which updates a generic asset
var UpdateAsset = Transaction{
	Tag:         "updateAsset",
	Label:       "Update Asset",
	Description: "",
	Method:      "PUT",

	MetaTx: true,
	Args: []Argument{
		{
			Tag:         "update",
			Description: "Asset key and fields to be updated.",
			DataType:    "@update",
			Required:    true,
		},
	},
	Routine: func(stub shim.ChaincodeStubInterface, req map[string]interface{}) ([]byte, errors.ICCError) {
		var err error
		request := req["update"].(map[string]interface{})
		key, err := assets.NewKey(request)
		if err != nil {
			return nil, errors.WrapError(err, "argument 'update' must be valid key")
		}

		// Check if asset exists
		asset, err := key.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to fetch asset from blockchain")
		}

		// Update asset
		response, err := asset.Update(stub, request)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update asset")
		}

		resBytes, err := json.Marshal(response)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}
