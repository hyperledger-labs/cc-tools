package assets

import (
	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ExistsInLedger checks if asset already exists
func (a *Asset) ExistsInLedger(stub shim.ChaincodeStubInterface) (bool, errors.ICCError) {
	var assetBytes []byte
	var err error
	if a.IsPrivate() {
		assetBytes, err = stub.GetPrivateData(a.TypeTag(), a.Key())
	} else {
		assetBytes, err = stub.GetState(a.Key())
	}
	if err != nil {
		return false, errors.WrapErrorWithStatus(err, "unable to check asset existence", 400)
	}
	if assetBytes != nil {
		return true, nil
	}

	return false, nil
}
