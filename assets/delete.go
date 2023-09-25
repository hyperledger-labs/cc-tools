package assets

import (
	"encoding/json"
	"strings"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

func (a *Asset) delete(stub *sw.StubWrapper) ([]byte, errors.ICCError) {
	var err error

	// Check if org has write permission
	err = a.CheckWriters(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed write permission check")
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

// Delete erases asset from world state and checks for all necessary permissions.
// An asset cannot be deleted if any other asset references it.
func (a *Asset) Delete(stub *sw.StubWrapper) ([]byte, errors.ICCError) {
	// Check if asset is referenced in other assets to avoid data inconsistency
	isReferenced, err := a.IsReferenced(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed to check if asset if being referenced")
	}
	if isReferenced {
		return nil, errors.NewCCError("another asset holds a reference to this one", 400)
	}

	return a.delete(stub)
}

// Delete erases asset from world state and checks for all necessary permissions.
// An asset cannot be deleted if any other asset references it.
func (k *Key) Delete(stub *sw.StubWrapper) ([]byte, errors.ICCError) {
	a, err := k.Get(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed to fetch asset from ledger")
	}

	return a.Delete(stub)
}

// DeleteCascade erases asset and recursively erases those which reference it.
// This method is experimental and might not work as intended. Use with caution.
func (k *Key) DeleteCascade(stub *sw.StubWrapper) ([]byte, errors.ICCError) {
	a, err := k.Get(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed to fetch asset from ledger")
	}

	return a.DeleteCascade(stub)
}

// DeleteCascade erases asset and recursively erases those which reference it.
// This method is experimental and might not work as intended. Use with caution.
func (a *Asset) DeleteCascade(stub *sw.StubWrapper) ([]byte, errors.ICCError) {
	var err error

	deletedKeys := []string{a.Key()}
	err = deleteRecursive(stub, a.Key(), &deletedKeys)
	if err != nil {
		return nil, errors.WrapError(err, "error deleting asset recursively")
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

func deleteRecursive(stub *sw.StubWrapper, key string, deletedKeys *[]string) errors.ICCError {
	states, err := stub.GetStateByPartialCompositeKey(key, []string{})
	if err != nil {
		return errors.WrapErrorWithStatus(err, "failed to check reference index", 500)
	}
	defer states.Close()

	for states.HasNext() {
		next, _ := states.Next()
		referrerKeyString := strings.ReplaceAll(next.Key, key, "")
		referrerKeyString = strings.ReplaceAll(referrerKeyString, "\x00", "")
		var isDeleted bool = false

		for _, deletedKey := range *deletedKeys {

			if deletedKey == referrerKeyString {
				isDeleted = true
				break
			}
		}
		*deletedKeys = append(*deletedKeys, referrerKeyString)

		if !isDeleted {
			err = deleteRecursive(stub, referrerKeyString, deletedKeys)
			if err != nil {
				return errors.WrapError(err, "error deleting referrer asset")
			}
		}

	}

	keyMap := make(map[string]interface{})
	keyMap["@key"] = key
	assetKey, err := NewKey(keyMap)
	if err != nil {
		return errors.WrapError(err, "failed to construct key")
	}

	asset, err := assetKey.Get(stub)
	if err != nil {
		return errors.WrapError(err, "failed to read asset from blockchain")
	}

	_, err = asset.delete(stub)
	if err != nil {
		return errors.WrapError(err, "failed to delete asset")
	}

	return nil
}
