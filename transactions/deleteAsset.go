package transactions

import (
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// DeleteAsset deletes an asset from the blockchain
var DeleteAsset = Transaction{
	Tag:         "deleteAsset",
	Label:       "Delete Asset",
	Method:      "DELETE",
	Description: "",

	MetaTx: true,
	Args: ArgList{
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
		var response []byte
		if cascade {
			response, err = key.DeleteCascade(stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to delete asset recursively")
			}
		} else {
			response, err = key.Delete(stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to delete asset")
			}
		}

		return response, nil
	},
}
