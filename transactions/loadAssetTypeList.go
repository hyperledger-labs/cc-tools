package transactions

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// LoadAssetTypeList is the transaction which loads the asset Type list from the blockchain
var LoadAssetTypeList = Transaction{
	Tag:         "loadAssetTypeList",
	Label:       "Load Asset Type List from blockchain",
	Description: "",
	Method:      "POST",

	MetaTx: true,
	Args:   ArgList{},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		err := assets.RestoreAssetList(stub, false)
		if err != nil {
			return nil, errors.WrapError(err, "failed to restore asset list")
		}

		resBytes, nerr := json.Marshal("Asset Type List loaded successfully")
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return resBytes, nil
	},
}
