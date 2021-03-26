package assets

import (
	"encoding/json"
	"strings"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Delete erases asset from world state and checks for all necessary permissions.
// An asset cannot be deleted if any other asset references it.
func (a *Asset) Delete(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var err error

	// Check if org has write permission
	err = a.CheckWriters(stub)
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

// DeleteRecursive erases asset and recursively erases those which reference it
func (a *Asset) DeleteRecursive(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var err error

	// Check if org has write permission
	err = a.CheckWriters(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed write permission check")
	}

	err = deleteRecursive(stub, a.Key())
	if err != nil {
		return nil, errors.WrapError(err, "error deleting asset Recursively")
	}

	response := map[string]interface{}{
		"deletedKeys": deletedKeys,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return nil, errors.WrapError(err, "failed to marshal response")
	}

	return responseJSON, nil
}

var deletedKeys []string

func deleteRecursive(stub shim.ChaincodeStubInterface, key string) error {
	deletedKeys = append(deletedKeys, key)
	queryIt, err := stub.GetStateByPartialCompositeKey(key, []string{})
	if err != nil {
		return errors.WrapErrorWithStatus(err, "failed to check reference index", 500)
	}
	defer queryIt.Close()

	for queryIt.HasNext() {
		next, _ := queryIt.Next()
		referrerKeyString := strings.ReplaceAll(next.Key, key, "")
		referrerKeyString = strings.ReplaceAll(referrerKeyString, "\x00", "")
		var isDeleted bool

		for _, deletedKey := range deletedKeys {
			if deletedKey == referrerKeyString {
				isDeleted = true
				break
			}
		}
		if !isDeleted {
			err = deleteRecursive(stub, referrerKeyString)
			if err != nil {
				return errors.WrapError(err, "error deleting referrer asset:")
			}
		}
	}

	keyMap := make(map[string]interface{})
	keyMap["@key"] = key
	assetKey, err := NewKey(keyMap)

	asset, err := assetKey.Get(stub)
	if err != nil {
		return errors.WrapError(err, "failed to read asset from blockchain")
	}
	// Clean up reference markers for this asset
	err = asset.delRefs(stub)
	if err != nil {
		return errors.WrapError(err, "failed cleaning reference index")
	}

	if !asset.IsPrivate() {
		err = stub.DelState(asset.Key())
		if err != nil {
			return errors.WrapError(err, "failed to delete state from ledger")
		}
	} else {
		err = stub.DelPrivateData(asset.TypeTag(), asset.Key())
		if err != nil {
			return errors.WrapError(err, "failed to delete state from private collection")
		}
	}

	return nil
}
