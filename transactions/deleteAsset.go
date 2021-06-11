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
		{
			Tag:         "cascade",
			Description: "Delete all referrers on cascade",
			DataType:    "boolean",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		// This is safe to do because validation is done before calling routine
		key := req["key"].(assets.Key)
		cascade, ok := req["cascade"].(bool)
		if !ok {
			cascade = false
		}

		var err error

		// Fetch asset from blockchain
		asset, err := key.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to read asset from blockchain")
		}

		response := make([]byte, 0)
		if cascade {
			response, err = asset.DeleteCascade(stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to delete asset recursively")
			}
		} else {

			response, err = asset.Delete(stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to delete asset")
			}
		}

		return response, nil
	},
}
