package assets

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Delete erases asset from world state
func (a *Asset) Delete(stub shim.ChaincodeStubInterface) ([]byte, error) {
	// Check if org has write permission
	err := a.CheckWriters(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed write permission check")
	}

	// Check if asset is referenced in other assets to avoid data inconsistency
	isReferenced, err := a.IsReferenced(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed to check if asset if being referenced")
	}
	if isReferenced {
		return nil, errors.NewCCError("another asset holds a reference to this one", 400)
	}

	// Clean up reference markers for this asset
	err = a.delRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed cleaning reference index")
	}

	var assetJSON []byte
	if !a.IsPrivate() {
		err = stub.DelState(a.Key())
		if err != nil {
			return nil, errors.WrapError(err, "failed to delete state from ledger")
		}
		assetJSON, err = json.Marshal(a)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal asset")
		}
	} else {
		err = stub.DelPrivateData(a.TypeTag(), a.Key())
		if err != nil {
			return nil, errors.WrapError(err, "failed to delete state from private collection")
		}
		assetKeyOnly := map[string]interface{}{
			"@key":       a.Key(),
			"@assetType": a.TypeTag(),
		}
		assetJSON, err = json.Marshal(assetKeyOnly)
		if err != nil {
			return nil, errors.WrapError(err, "failed to marshal private asset key")
		}
	}

	return assetJSON, nil
}
