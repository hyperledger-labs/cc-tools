package assets

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Get reads asset from ledger
func (a *Asset) Get(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var assetBytes []byte
	var err error
	if a.IsPrivate() {
		assetBytes, err = stub.GetPrivateData(a.TypeTag(), a.Key())
	} else {
		assetBytes, err = stub.GetState(a.Key())
	}
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "unable to get asset", 400)
	}
	if assetBytes == nil {
		return nil, errors.NewCCError("asset not found", 404)
	}

	var response Asset
	err = json.Unmarshal(assetBytes, &response)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal asset from ledger", 500)
	}

	return &response, nil
}
