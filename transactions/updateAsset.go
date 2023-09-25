package transactions

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// UpdateAsset is the function which updates a generic asset
var UpdateAsset = Transaction{
	Tag:         "updateAsset",
	Label:       "Update Asset",
	Description: "",
	Method:      "PUT",

	MetaTx: true,
	Args: ArgList{
		{
			Tag:         "update",
			Description: "Asset key and fields to be updated.",
			DataType:    "@update",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var err error
		request := req["update"].(map[string]interface{})
		key, _ := assets.NewKey(request)

		// Check if asset exists
		exists, err := key.ExistsInLedger(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to check asset existance in ledger")
		}
		if !exists {
			return nil, errors.NewCCError("asset does not exist", 404)
		}

		// Update asset
		response, err := key.Update(stub, request)
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
