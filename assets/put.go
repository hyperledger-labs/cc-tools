package assets

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (a *Asset) put(stub shim.ChaincodeStubInterface) (map[string]interface{}, errors.ICCError) {
	// Clean asset of any nil entries
	a.clean()

	// Write index of references this asset points to
	err := a.putRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed writing reference index")
	}

	err = a.validateRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed reference validation")
	}

	// Marshal asset back to JSON format
	assetJSON, err := json.Marshal(a)
	if err != nil {
		return nil, errors.WrapError(err, "failed to encode asset to JSON format")
	}

	// Write asset to blockchain
	if a.IsPrivate() {
		err = stub.PutPrivateData(a.TypeTag(), a.Key(), assetJSON)
		if err != nil {
			return nil, errors.WrapError(err, "failed to write asset to ledger")
		}
		assetKeyOnly := map[string]interface{}{
			"@key":       a.Key(),
			"@assetType": a.TypeTag(),
		}
		return assetKeyOnly, nil
	}

	err = stub.PutState(a.Key(), assetJSON)
	if err != nil {
		return nil, errors.WrapError(err, "failed to write asset to ledger")
	}
	return *a, nil
}

// Put inserts asset in blockchain
func (a *Asset) Put(stub shim.ChaincodeStubInterface) (map[string]interface{}, errors.ICCError) {
	// Check if org has write permission
	err := a.CheckWriters(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed write permission check")
	}

	err = a.injectMetadata(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed injecting asset metadata")
	}

	return a.put(stub)
}

// PutNew inserts asset in blockchain and returns error if asset exists
func (a *Asset) PutNew(stub shim.ChaincodeStubInterface) (map[string]interface{}, errors.ICCError) {
	// Check if asset already exists
	exists, err := a.ExistsInLedger(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed to verify if asset already exists")
	}
	if exists {
		return nil, errors.NewCCError("asset already exists", 409)
	}

	// Marshal asset back to JSON format
	res, err := a.Put(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed to write asset to ledger")
	}

	return res, nil
}
