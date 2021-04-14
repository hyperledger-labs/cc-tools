package assets

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func get(stub shim.ChaincodeStubInterface, pvtCollection, key string) (*Asset, errors.ICCError) {
	var assetBytes []byte
	var err error
	if pvtCollection != "" {
		assetBytes, err = stub.GetPrivateData(pvtCollection, key)
	} else {
		assetBytes, err = stub.GetState(key)
	}
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "unable to get asset", 400)
	}
	if assetBytes == nil {
		return nil, errors.NewCCError("asset not found", 404)
	}

	var response Asset
	err = json.Unmarshal(assetBytes, &response)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal asset from ledger", 500)
	}

	return &response, nil
}

// Get fetches asset entry from ledger.
func (a *Asset) Get(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var pvtCollection string
	if a.IsPrivate() {
		pvtCollection = a.TypeTag()
	}

	return get(stub, pvtCollection, a.Key())
}

// Get fetches asset entry from ledger.
func (k *Key) Get(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var pvtCollection string
	if k.IsPrivate() {
		pvtCollection = k.TypeTag()
	}

	return get(stub, pvtCollection, k.Key())
}

/* GetRecursive-related code */

func getRecursive(stub shim.ChaincodeStubInterface, pvtCollection, key string, keysChecked []string) (*Asset, errors.ICCError) {
	var assetBytes []byte
	var err error
	if pvtCollection != "" {
		assetBytes, err = stub.GetPrivateData(pvtCollection, key)
		// If org cannot get private data it might be because it has no permission, so we fetch the data hash
		if err != nil {
			hash, err := stub.GetPrivateDataHash(pvtCollection, key)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "unable to get asset", 400)
			}
			if hash == nil {
				return nil, errors.NewCCError("asset not found", 404)
			}
			response := Asset{
				"@key":       key,
				"@assetType": pvtCollection,
				"@hash":      hash,
			}
			return &response, nil
		}
	} else {
		assetBytes, err = stub.GetState(key)
	}
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "unable to get asset", 400)
	}
	if assetBytes == nil {
		return nil, errors.NewCCError("asset not found", 404)
	}

	var response Asset
	err = json.Unmarshal(assetBytes, &response)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal asset from ledger", 500)
	}

	for k, v := range response {
		switch prop := v.(type) {
		case map[string]interface{}:
			propKey, err := NewKey(prop)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed to resolve asset references", 500)
			}

			keyIsFetched := false
			for _, key := range keysChecked {
				if key == propKey.Key() {
					keyIsFetched = true
					break
				}
			}
			if keyIsFetched {
				continue
			}
			keysChecked = append(keysChecked, propKey.Key())

			var subAsset *Asset
			if propKey.IsPrivate() {
				subAsset, err = getRecursive(stub, propKey.TypeTag(), propKey.Key(), keysChecked)
			} else {
				subAsset, err = getRecursive(stub, "", propKey.Key(), keysChecked)
			}
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed to get subasset", 500)
			}

			response[k] = *subAsset

		case []interface{}:
			for idx, elem := range prop {
				if elemMap, ok := elem.(map[string]interface{}); ok {
					elemKey, err := NewKey(elemMap)
					if err != nil {
						return nil, errors.WrapErrorWithStatus(err, "failed to resolve asset references", 500)
					}

					keyIsFetched := false
					for _, key := range keysChecked {
						if key == elemKey.Key() {
							keyIsFetched = true
							break
						}
					}
					if keyIsFetched {
						continue
					}
					keysChecked = append(keysChecked, elemKey.Key())

					var subAsset *Asset
					if elemKey.IsPrivate() {
						subAsset, err = getRecursive(stub, elemKey.TypeTag(), elemKey.Key(), keysChecked)
					} else {
						subAsset, err = getRecursive(stub, "", elemKey.Key(), keysChecked)
					}
					if err != nil {
						return nil, errors.WrapErrorWithStatus(err, "failed to get subasset", 500)
					}

					prop[idx] = *subAsset
				}
			}
			response[k] = prop
		}
	}

	return &response, nil
}

// GetRecursive reads asset from ledger and resolves all references
func (a *Asset) GetRecursive(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var pvtCollection string
	if a.IsPrivate() {
		pvtCollection = a.TypeTag()
	}

	return getRecursive(stub, pvtCollection, a.Key(), []string{})
}

// GetRecursive reads asset from ledger and resolves all references
func (k *Key) GetRecursive(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var pvtCollection string
	if k.IsPrivate() {
		pvtCollection = k.TypeTag()
	}

	return getRecursive(stub, pvtCollection, k.Key(), []string{})
}
