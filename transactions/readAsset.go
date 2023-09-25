package transactions

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// ReadAsset fetches an asset from the blockchain
var ReadAsset = Transaction{
	Tag:         "readAsset",
	Label:       "Read Asset",
	Description: "",
	Method:      "GET",

	MetaTx: true,
	Args: ArgList{
		{
			Tag:         "key",
			Description: "Key of the asset to be read.",
			DataType:    "@key",
			Required:    true,
		},
		{
			Tag:         "resolve",
			Description: "Resolve references recursively.",
			DataType:    "boolean",
		},
	},
	ReadOnly: true,
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var assetJSON []byte
		var err error

		// This is safe to do because validation is done before calling routine
		key := req["key"].(assets.Key)

		resolve, ok := req["resolve"].(bool)

		if ok && resolve {
			var asset map[string]interface{}
			asset, err = key.GetRecursive(stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to read asset from blockchain")
			}

			assetJSON, err = json.Marshal(asset)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed to serialize asset", 500)
			}
		} else {
			assetJSON, err = key.GetBytes(stub)
			if err != nil {
				return nil, errors.WrapError(err, "failed to get asset state")
			}
		}

		return assetJSON, nil
	},
}
