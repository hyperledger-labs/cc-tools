package transactions

import (
	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// DeleteAsset deletes an asset from the blockchain
var DeleteAsset = Transaction{
	Tag:         "deleteAsset",
	Label:       "Delete Asset",
	Method:      "DELETE",
	Description: "",

	MetaTx: true,
	Args: []Argument{
		{
			Tag:         "key",
			Description: "Key of the asset to be deleted.",
			DataType:    "@key",
			Required:    true,
		},
	},
	Routine: func(stub shim.ChaincodeStubInterface, req map[string]interface{}) ([]byte, errors.ICCError) {
		// This is safe to do because validation is done before calling routine
		key := req["key"].(assets.Key)

		var err error

		// Fetch asset from blockchain
		asset, err := key.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to read asset from blockchain")
		}

		response, err := asset.Delete(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to delete asset")
		}

		return response, nil
	},
}

// DeleteAssetForced deletes an asset from the blockchain even if it is referenced
var DeleteAssetForced = Transaction{
	Tag:         "deleteAssetForced",
	Label:       "Delete Asset Forced",
	Method:      "DELETE",
	Description: "",

	MetaTx: true,
	Args: []Argument{
		{
			Tag:         "key",
			Description: "Key of the asset to be deleted.",
			DataType:    "@key",
			Required:    true,
		},
	},
	Routine: func(stub shim.ChaincodeStubInterface, req map[string]interface{}) ([]byte, errors.ICCError) {
		// This is safe to do because validation is done before calling routine
		key := req["key"].(assets.Key)

		var err error

		// Fetch asset from blockchain
		asset, err := key.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to read asset from blockchain")
		}

		response, err := asset.DeleteForced(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to delete asset")
		}

		return response, nil
	},
}
