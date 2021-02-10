package assets

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (a *Asset) put(stub shim.ChaincodeStubInterface) (map[string]interface{}, errors.ICCError) {
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

	err = a.validateRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed reference validation")
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

func putRecursive(stub shim.ChaincodeStubInterface, object map[string]interface{}, root bool) (map[string]interface{}, errors.ICCError) {
	var err error

	objAsAsset, err := NewAsset(object)
	if err != nil {
		return nil, errors.WrapError(err, "unable to create asset object")
	}

	if !root {
		exists, err := objAsAsset.ExistsInLedger(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed checking if asset exists")
		}
		if exists {
			asset, err := objAsAsset.GetRecursive(stub)
			return *asset, err
		}
	}

	putAsset, err := objAsAsset.put(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed to put asset")
	}

	subAssets := objAsAsset.Type().SubAssets()
	for _, subAsset := range subAssets {
		isArray := false
		dType := subAsset.DataType
		if strings.HasPrefix(dType, "[]") {
			isArray = true
			dType = strings.TrimPrefix(dType, "[]")
		}

		dType = strings.TrimPrefix(dType, "->")
		subAssetInterface, ok := object[subAsset.Tag]
		if !ok {
			// if subAsset is not included, continue onwards to the next possible subAsset
			continue
		}

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
			obj, ok := objInterface.(map[string]interface{})
			if !ok {
				return nil, errors.NewCCError(fmt.Sprintf("asset property %s must of type %s", subAsset.Label, subAsset.DataType), 400)
			}
			putSubAsset, err := putRecursive(stub, obj, false)
			if err != nil {
				return nil, errors.WrapError(err, "failed to put sub-asset recursively")
			}
			objArray[idx] = putSubAsset
		}

		if isArray {
			putAsset[subAsset.Tag] = objArray
		} else {
			putAsset[subAsset.Tag] = objArray[0]
		}
	}

	return putAsset, nil
}

// PutRecursive inserts asset and all it's subassets in blockchain
func PutRecursive(stub shim.ChaincodeStubInterface, object map[string]interface{}) (map[string]interface{}, errors.ICCError) {
	return putRecursive(stub, object, true)
}

// PutNewRecursive inserts asset and all it's subassets in blockchain and returns error if asset exists
func PutNewRecursive(stub shim.ChaincodeStubInterface, object map[string]interface{}) (map[string]interface{}, errors.ICCError) {
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
