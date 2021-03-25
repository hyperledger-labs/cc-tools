package assets

import (
	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ExistsInLedger checks if asset currently has a state on the ledger.
func (a *Asset) ExistsInLedger(stub shim.ChaincodeStubInterface) (bool, errors.ICCError) {
	var assetBytes []byte
	var err error
	if a.IsPrivate() {
		assetBytes, err = stub.GetPrivateDataHash(a.TypeTag(), a.Key())
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

// ExistsInLedger checks if asset referenced by a key object currently has a state on the ledger.
func (k *Key) ExistsInLedger(stub shim.ChaincodeStubInterface) (bool, errors.ICCError) {
	var assetBytes []byte
	var err error
	if k.IsPrivate() {
		assetBytes, err = stub.GetPrivateDataHash(k.TypeTag(), k.Key())
	} else {
		assetBytes, err = stub.GetState(k.Key())
	}
	if err != nil {
		return false, errors.WrapErrorWithStatus(err, "unable to check asset existence", 400)
	}
	if assetBytes != nil {
		return true, nil
	}

	return false, nil
}
