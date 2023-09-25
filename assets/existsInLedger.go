package assets

import (
	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/cc-tools/mock"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

func existsInLedger(stub *sw.StubWrapper, isPrivate bool, typeTag, key string) (bool, errors.ICCError) {
	var assetBytes []byte
	var err error
	if isPrivate {
		_, isMock := stub.Stub.(*mock.MockStub)
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

// ExistsInLedger checks if asset currently has a state.
func (a *Asset) ExistsInLedger(stub *sw.StubWrapper) (bool, errors.ICCError) {
	if a.Key() == "" {
		return false, errors.NewCCError("asset key is empty", 500)
	}
	return existsInLedger(stub, a.IsPrivate(), a.TypeTag(), a.Key())
}

// ExistsInLedger checks if asset referenced by a key object currently has a state.
func (k *Key) ExistsInLedger(stub *sw.StubWrapper) (bool, errors.ICCError) {
	if k.Key() == "" {
		return false, errors.NewCCError("key is empty", 500)
	}
	return existsInLedger(stub, k.IsPrivate(), k.TypeTag(), k.Key())
}

// ----------------------------------------

func committedInLedger(stub *sw.StubWrapper, isPrivate bool, typeTag, key string) (bool, errors.ICCError) {
	var assetBytes []byte
	var err error
	if isPrivate {
		_, isMock := stub.Stub.(*mock.MockStub)
		if isMock {
			assetBytes, err = stub.Stub.GetPrivateData(typeTag, key)
		} else {
			assetBytes, err = stub.Stub.GetPrivateDataHash(typeTag, key)
		}
	} else {
		assetBytes, err = stub.Stub.GetState(key)
	}
	if err != nil {
		return false, errors.WrapErrorWithStatus(err, "unable to check asset existence", 400)
	}
	if assetBytes != nil {
		return true, nil
	}

	return false, nil
}

// CommittedInLedger checks if asset currently has a state committed in ledger.
func (a *Asset) CommittedInLedger(stub *sw.StubWrapper) (bool, errors.ICCError) {
	if a.Key() == "" {
		return false, errors.NewCCError("asset key is empty", 500)
	}
	return committedInLedger(stub, a.IsPrivate(), a.TypeTag(), a.Key())
}

// CommittedInLedger checks if asset referenced by a key object currently has a state committed in ledger.
func (k *Key) CommittedInLedger(stub *sw.StubWrapper) (bool, errors.ICCError) {
	if k.Key() == "" {
		return false, errors.NewCCError("key is empty", 500)
	}
	return committedInLedger(stub, k.IsPrivate(), k.TypeTag(), k.Key())
}
