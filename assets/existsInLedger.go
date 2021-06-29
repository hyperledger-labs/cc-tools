package assets

import (
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func existsInLedger(stub *sw.StubWrapper, isPrivate bool, typeTag, key string) (bool, errors.ICCError) {
	var assetBytes []byte
	var err error
	if isPrivate {
		_, isMock := stub.Stub.(*shim.MockStub)
		if isMock {
			assetBytes, err = stub.GetPrivateData(typeTag, key)
		} else {
			assetBytes, err = stub.GetPrivateDataHash(typeTag, key)
		}
	} else {
		assetBytes, err = stub.GetState(key)
	}
	if err != nil {
		return false, errors.WrapErrorWithStatus(err, "unable to check asset existence", 400)
	}
	if assetBytes != nil {
		return true, nil
	}

	return false, nil
}

// ExistsInLedger checks if asset currently has a state on the ledger.
func (a *Asset) ExistsInLedger(stub *sw.StubWrapper) (bool, errors.ICCError) {
	return existsInLedger(stub, a.IsPrivate(), a.TypeTag(), a.Key())
}

// ExistsInLedger checks if asset referenced by a key object currently has a state on the ledger.
func (k *Key) ExistsInLedger(stub *sw.StubWrapper) (bool, errors.ICCError) {
	return existsInLedger(stub, k.IsPrivate(), k.TypeTag(), k.Key())
}
