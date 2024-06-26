package assets

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// put writes the reference index to the ledger, then encodes the
// asset to JSON format and puts it into the ledger.
func (a *Asset) put(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {
	var err error

	// Clean asset of any nil entries
	a.clean()

	// Write index of references this asset points to
	err = a.putRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed writing reference index")
	}

	// Marshal asset back to JSON format
	assetJSON, err := json.Marshal(a)
	if err != nil {
		return nil, errors.WrapError(err, "failed to encode asset to JSON format")
	}

	// Write asset to blockchain
	if a.IsPrivate() {
		err = stub.PutPrivateData(a.CollectionName(), a.Key(), assetJSON)
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
func (a *Asset) Put(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {
	// Check if org has write permission
	err := a.CheckWriters(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed write permission check")
	}

	err = a.injectMetadata(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed injecting asset metadata")
	}

	err = a.validateRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed reference validation")
	}
	return a.put(stub)
}

// PutNew inserts asset in blockchain and returns error if asset exists.
func (a *Asset) PutNew(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {
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

func putRecursive(stub *sw.StubWrapper, object map[string]interface{}) (map[string]interface{}, errors.ICCError) {
	var err error

	objAsKey, err := NewKey(object)
	if err != nil {
		return nil, errors.WrapError(err, "unable to create asset object")
	}

	exists, err := objAsKey.ExistsInLedger(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed checking if asset exists")
	}

	propsToUpdate := map[string]bool{}
	if exists {
		asset, err := objAsKey.GetMap(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed fetching asset that already exists")
		}

		// If asset key is not in object, add asset value to object (so that properties are not erased)
		for k := range asset {
			if _, ok := object[k]; !ok {
				object[k] = asset[k]
			}
		}

		// Check props to update
		for k, v := range object {
			if !reflect.DeepEqual(v, asset[k]) {
				propsToUpdate[k] = true
			}
		}
	}

	subAssetsMap := map[string]interface{}{}
	subAssets := objAsKey.Type().SubAssets()
	for _, subAsset := range subAssets {
		subAssetInterface, ok := object[subAsset.Tag]
		if !ok {
			// if subAsset is not included, continue onwards to the next possible subAsset
			continue
		}

		if propsToUpdate[subAsset.Tag] {
			delete(propsToUpdate, subAsset.Tag)
		}

		// Extract asset type
		isArray := false
		dType := subAsset.DataType
		if strings.HasPrefix(dType, "[]") {
			isArray = true
			dType = strings.TrimPrefix(dType, "[]")
		}
		dType = strings.TrimPrefix(dType, "->")

		var objArray []interface{}
		if !isArray {
			objArray = []interface{}{subAssetInterface}
		} else {
			objArray, ok = subAssetInterface.([]interface{})
			if !ok {
				return nil, errors.NewCCError(fmt.Sprintf("asset property %s must an array of type %s", subAsset.Label, subAsset.DataType), 400)
			}
		}

		for idx, objInterface := range objArray {
			var obj map[string]interface{}
			switch t := objInterface.(type) {
			case map[string]interface{}:
				obj = t
			case Key:
				obj = t
			case Asset:
				obj = t
			default:
				// If subAsset is badly formatted, this method shouldn't have been called
				return nil, errors.NewCCError(fmt.Sprintf("asset reference property '%s' must be an object", subAsset.Tag), 400)
			}
			if dType != "@asset" {
				obj["@assetType"] = dType
			} else {
				_, ok := obj["@assetType"].(string)
				if !ok {
					return nil, errors.NewCCError(fmt.Sprintf("asset reference property '%s' must have an '@assetType' property", subAsset.Tag), 400)
				}
			}
			putSubAsset, err := putRecursive(stub, obj)
			if err != nil {
				return nil, errors.WrapError(err, fmt.Sprintf("failed to put sub-asset %s recursively", subAsset.Tag))
			}
			objArray[idx] = putSubAsset
		}

		if isArray {
			object[subAsset.Tag] = objArray
		} else {
			object[subAsset.Tag] = objArray[0]
		}
		subAssetsMap[subAsset.Tag] = object[subAsset.Tag]
	}

	if shouldUpdate := len(propsToUpdate) > 0; shouldUpdate || !exists {
		objAsAsset, err := NewAsset(object)
		if err != nil {
			return nil, errors.WrapError(err, "unable to create asset object")
		}

		object, err = objAsAsset.Put(stub)
		if err != nil {
			return nil, errors.WrapError(err, fmt.Sprintf("failed to put asset of type %s", objAsAsset.TypeTag()))
		}

		for tag, subAsset := range subAssetsMap {
			object[tag] = subAsset
		}
	}

	return object, nil
}

// PutRecursive inserts asset and all its subassets in blockchain.
// This method is experimental and might not work as intended. Use with caution.
func PutRecursive(stub *sw.StubWrapper, object map[string]interface{}) (map[string]interface{}, errors.ICCError) {
	return putRecursive(stub, object)
}

// PutNewRecursive inserts asset and all its subassets in blockchain
// This method is experimental and might not work as intended. Use with caution.
// It returns conflict error only if root asset exists.
// If one of the subassets already exist in ledger, it is not updated.
func PutNewRecursive(stub *sw.StubWrapper, object map[string]interface{}) (map[string]interface{}, errors.ICCError) {
	objAsAsset, err := NewAsset(object)
	if err != nil {
		return nil, errors.WrapError(err, "unable to create asset object")
	}

	exists, err := objAsAsset.ExistsInLedger(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed checking if asset exists")
	}
	if exists {
		return nil, errors.NewCCError("asset already exists", 409)
	}

	return PutRecursive(stub, object)
}
