package assets

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

func get(stub *sw.StubWrapper, pvtCollection, key string, committed bool) (*Asset, errors.ICCError) {
	var assetBytes []byte
	var err error

	if key == "" {
		return nil, errors.NewCCError("key cannot be empty", 500)
	}

	if committed {
		if pvtCollection != "" {
			assetBytes, err = stub.GetCommittedPrivateData(pvtCollection, key)
		} else {
			assetBytes, err = stub.GetCommittedState(key)
		}
	} else {
		if pvtCollection != "" {
			assetBytes, err = stub.GetPrivateData(pvtCollection, key)
		} else {
			assetBytes, err = stub.GetState(key)
		}
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

// Get fetches asset entry from write set or ledger.
func (a *Asset) Get(stub *sw.StubWrapper) (*Asset, errors.ICCError) {
	var pvtCollection string
	if a.IsPrivate() {
		pvtCollection = a.TypeTag()
	}

	return get(stub, pvtCollection, a.Key(), false)
}

// Get fetches asset entry from write set or ledger.
func (k *Key) Get(stub *sw.StubWrapper) (*Asset, errors.ICCError) {
	var pvtCollection string
	if k.IsPrivate() {
		pvtCollection = k.TypeTag()
	}

	return get(stub, pvtCollection, k.Key(), false)
}

// GetMany fetches assets entries from write set or ledger.
func GetMany(stub *sw.StubWrapper, keys []Key) ([]*Asset, errors.ICCError) {
	var assets []*Asset

	for _, k := range keys {
		var pvtCollection string
		if k.IsPrivate() {
			pvtCollection = k.TypeTag()
		}

		asset, err := get(stub, pvtCollection, k.Key(), false)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

// GetCommitted fetches asset entry from ledger.
func (a *Asset) GetCommitted(stub *sw.StubWrapper) (*Asset, errors.ICCError) {
	var pvtCollection string
	if a.IsPrivate() {
		pvtCollection = a.TypeTag()
	}

	return get(stub, pvtCollection, a.Key(), true)
}

// GetCommitted fetches asset entry from ledger.
func (k *Key) GetCommitted(stub *sw.StubWrapper) (*Asset, errors.ICCError) {
	var pvtCollection string
	if k.IsPrivate() {
		pvtCollection = k.TypeTag()
	}

	return get(stub, pvtCollection, k.Key(), true)
}

// GetBytes reads the asset as bytes from ledger
func (k *Key) GetBytes(stub *sw.StubWrapper) ([]byte, errors.ICCError) {
	var assetBytes []byte
	var err error
	if k.IsPrivate() {
		assetBytes, err = stub.GetPrivateData(k.TypeTag(), k.Key())
	} else {
		assetBytes, err = stub.GetState(k.Key())
	}
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to get asset bytes", 400)
	}
	if assetBytes == nil {
		return nil, errors.NewCCError("asset not found", 404)
	}

	return assetBytes, nil
}

// GetMap reads the asset as map from ledger
func (k *Key) GetMap(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {
	var err error
	assetBytes, err := k.GetBytes(stub)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to get asset bytes", 400)
	}

	var ret map[string]interface{}
	err = json.Unmarshal(assetBytes, &ret)
	if err != nil {
		return nil, errors.WrapError(err, "failed to unmarshal asset")
	}

	return ret, nil
}

/* GetRecursive-related code */

func getRecursive(stub *sw.StubWrapper, pvtCollection, key string, keysChecked []string) (map[string]interface{}, errors.ICCError) {
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
			response := map[string]interface{}{
				"@key":       key,
				"@assetType": pvtCollection,
				"@hash":      hash,
			}
			return response, nil
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

	var response map[string]interface{}
	err = json.Unmarshal(assetBytes, &response)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal asset from ledger", 500)
	}

	keysCheckedInScope := make([]string, 0)

	for k, v := range response {
		switch prop := v.(type) {
		case map[string]interface{}:

			assetType, ok := prop["@assetType"].(string)
			if !ok || assetType == "@object" {
				continue
			}

			propKey, err := NewKey(prop)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed to resolve asset references", 500)
			}

			keyIsFetchedInScope := false
			for _, key := range keysCheckedInScope {
				if key == propKey.Key() {
					keyIsFetchedInScope = true
					break
				}
			}

			keyIsFetched := false
			for _, key := range keysChecked {
				if key == propKey.Key() {
					keyIsFetched = true
					break
				}
			}
			if keyIsFetched && !keyIsFetchedInScope {
				continue
			}
			keysChecked = append(keysChecked, propKey.Key())
			keysCheckedInScope = append(keysCheckedInScope, propKey.Key())

			var subAsset map[string]interface{}
			if propKey.IsPrivate() {
				subAsset, err = getRecursive(stub, propKey.TypeTag(), propKey.Key(), keysChecked)
			} else {
				subAsset, err = getRecursive(stub, "", propKey.Key(), keysChecked)
			}
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed to get subasset", 500)
			}

			response[k] = subAsset

		case []interface{}:
			for idx, elem := range prop {
				if elemMap, ok := elem.(map[string]interface{}); ok {

					assetType, ok := elemMap["@assetType"].(string)
					if !ok || assetType == "@object" {
						continue
					}

					elemKey, err := NewKey(elemMap)
					if err != nil {
						return nil, errors.WrapErrorWithStatus(err, "failed to resolve asset references", 500)
					}

					keyIsFetchedInScope := false
					for _, key := range keysCheckedInScope {
						if key == elemKey.Key() {
							keyIsFetchedInScope = true
							break
						}
					}

					keyIsFetched := false
					for _, key := range keysChecked {
						if key == elemKey.Key() {
							keyIsFetched = true
							break
						}
					}
					if keyIsFetched && !keyIsFetchedInScope {
						continue
					}
					keysChecked = append(keysChecked, elemKey.Key())
					keysCheckedInScope = append(keysCheckedInScope, elemKey.Key())

					var subAsset map[string]interface{}
					if elemKey.IsPrivate() {
						subAsset, err = getRecursive(stub, elemKey.TypeTag(), elemKey.Key(), keysChecked)
					} else {
						subAsset, err = getRecursive(stub, "", elemKey.Key(), keysChecked)
					}
					if err != nil {
						return nil, errors.WrapErrorWithStatus(err, "failed to get subasset", 500)
					}

					prop[idx] = subAsset
				}
			}
			response[k] = prop
		}
	}

	return response, nil
}

// GetRecursive reads asset from ledger and resolves all references
func (a *Asset) GetRecursive(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {
	var pvtCollection string
	if a.IsPrivate() {
		pvtCollection = a.TypeTag()
	}

	return getRecursive(stub, pvtCollection, a.Key(), []string{})
}

// GetRecursive reads asset from ledger and resolves all references
func (k *Key) GetRecursive(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {
	var pvtCollection string
	if k.IsPrivate() {
		pvtCollection = k.TypeTag()
	}

	return getRecursive(stub, pvtCollection, k.Key(), []string{})
}
