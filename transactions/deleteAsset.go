package transactions

import (
	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
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
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
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

// DeleteRecursive deletes an asset and it's recursive referrers from the blockchain
var DeleteRecursive = Transaction{
	Tag:         "deleteAssetRecursive",
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
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		// This is safe to do because validation is done before calling routine
		key := req["key"].(assets.Key)

		var err error

		// Fetch asset from blockchain
		asset, err := key.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to read asset from blockchain")
		}

		response, err := asset.DeleteRecursive(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to delete asset recursively")
		}

		return response, nil
	},
}
