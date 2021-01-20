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

// Get reads asset from ledger
func (a *Asset) Get(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var pvtCollection string
	if a.IsPrivate() {
		pvtCollection = a.TypeTag()
	}

	return get(stub, pvtCollection, a.Key())
}

// Get reads asset from ledger
func (k *Key) Get(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var pvtCollection string
	if k.IsPrivate() {
		pvtCollection = k.TypeTag()
	}

	return get(stub, pvtCollection, k.Key())
}

/* GetRecursive-related code */

func getRecursive(stub shim.ChaincodeStubInterface, pvtCollection, key string) (*Asset, errors.ICCError) {
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
		var fullValue interface{}
		switch prop := v.(type) {
		case map[string]interface{}:
			propKey, err := NewKey(prop)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed to resolve asset references", 500)
			}

			var subAsset *Asset
			if propKey.IsPrivate() {
				subAsset, err = getRecursive(stub, propKey.TypeTag(), propKey.Key())
			} else {
				subAsset, err = getRecursive(stub, "", propKey.Key())
			}
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed to get subasset", 500)
			}

			fullValue = *subAsset
		case []interface{}:
			for idx, elem := range prop {
				if elemMap, ok := elem.(map[string]interface{}); ok {
					elemKey, err := NewKey(elemMap)
					if err != nil {
						return nil, errors.WrapErrorWithStatus(err, "failed to resolve asset references", 500)
					}

					var subAsset *Asset
					if elemKey.IsPrivate() {
						subAsset, err = getRecursive(stub, elemKey.TypeTag(), elemKey.Key())
					} else {
						subAsset, err = getRecursive(stub, "", elemKey.Key())
					}
					if err != nil {
						return nil, errors.WrapErrorWithStatus(err, "failed to get subasset", 500)
					}

					prop[idx] = *subAsset
				}
			}
			fullValue = prop
		default:
			continue
		}
		response[k] = fullValue
	}

	return &response, nil
}

// GetRecursive reads asset from ledger and resolves all references
func (a *Asset) GetRecursive(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var pvtCollection string
	if a.IsPrivate() {
		pvtCollection = a.TypeTag()
	}

	return getRecursive(stub, pvtCollection, a.Key())
}

// GetRecursive reads asset from ledger and resolves all references
func (k *Key) GetRecursive(stub shim.ChaincodeStubInterface) (*Asset, errors.ICCError) {
	var pvtCollection string
	if k.IsPrivate() {
		pvtCollection = k.TypeTag()
	}

	return getRecursive(stub, pvtCollection, k.Key())
}
