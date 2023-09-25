package assets

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// Refs returns an array of Keys containing the reference keys for all present subAssets.
// The referenced keys are fetched based on current asset configuration.
func (a Asset) Refs() ([]Key, errors.ICCError) {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return nil, errors.NewCCError(fmt.Sprintf("asset type named '%s' does not exist", a.TypeTag()), 400)
	}
	assetSubAssets := assetTypeDef.SubAssets()
	var keys []Key
	for _, subAsset := range assetSubAssets {
		subAssetRefInterface, subAssetIncluded := a[subAsset.Tag]
		if !subAssetIncluded || subAssetRefInterface == nil {
			// If subAsset is not included, no need to append
			continue
		}

		var isArray bool
		subAssetDataType := subAsset.DataType
		if strings.HasPrefix(subAssetDataType, "[]") {
			subAssetDataType = strings.TrimPrefix(subAssetDataType, "[]")
			isArray = true
		}

		subAssetDataType = strings.TrimPrefix(subAssetDataType, "->")

		// Handle array-like sub-asset property types
		var subAssetAsArray []interface{}
		if !isArray {
			subAssetAsArray = []interface{}{subAssetRefInterface}
		} else {
			var ok bool
			subAssetAsArray, ok = subAssetRefInterface.([]interface{})
			if !ok {
				return nil, errors.NewCCError(fmt.Sprintf("asset property '%s' must be an array of type '%s'", subAsset.Label, subAsset.DataType), 400)
			}
		}

		for _, subAssetRefInterface := range subAssetAsArray {
			// This is here as a safety measure
			if subAssetRefInterface == nil {
				continue
			}

			var subAssetRefMap map[string]interface{}
			switch t := subAssetRefInterface.(type) {
			case map[string]interface{}:
				subAssetRefMap = t
			case Key:
				subAssetRefMap = t
			case Asset:
				subAssetRefMap = t
			default:
				// If subAsset is badly formatted, this method shouldn't have been called
				return nil, errors.NewCCError("asset reference must be an object", 400)
			}

			subAssetTypeName, ok := subAssetRefMap["@assetType"]
			if ok && subAssetTypeName != subAssetDataType {
				return nil, errors.NewCCError("sub-asset reference of wrong asset type", 400)
			}
			if !ok {
				subAssetRefMap["@assetType"] = subAssetDataType
			}

			// Generate key for subAsset
			key, err := NewKey(subAssetRefMap)
			if err != nil {
				return nil, errors.WrapError(err, "failed to generate unique identifier for asset")
			}

			// Append key to response
			keys = append(keys, key)
		}
	}
	return keys, nil
}

// Refs returns an array of Keys containing the reference keys for all present subAssets.
// The referenced keys are fetched from current asset state on ledger.
func (k Key) Refs(stub *sw.StubWrapper) ([]Key, errors.ICCError) {
	assetMap, err := k.GetMap(stub)
	if err != nil {
		return nil, errors.WrapError(err, "could not get asset map from ledger")
	}

	var keys []Key
	for _, prop := range assetMap {
		// Convert everything to interface slice to handle every sub-asset with the same code
		var subAssetAsInterfaceSlice []interface{}
		switch v := prop.(type) {
		case map[string]interface{}:
			subAssetAsInterfaceSlice = []interface{}{v}
		case []interface{}:
			subAssetAsInterfaceSlice = v
		default:
			// Whatever is not a map[string]interface{} or a []interface{} can never be a sub-asset
			continue
		}

		for _, elem := range subAssetAsInterfaceSlice {
			switch v := elem.(type) {
			case map[string]interface{}:
				newKey, err := NewKey(v)
				if err != nil {
					continue
				}
				keys = append(keys, newKey)
			default:
				continue
			}
		}
	}

	return keys, nil
}

// Referrers returns an array of Keys of all the assets pointing to asset.
// assetTypeFilter can be used to filter the results by asset type.
func (a Asset) Referrers(stub *sw.StubWrapper, assetTypeFilter ...string) ([]Key, errors.ICCError) {
	assetKey := a.Key()
	return referrers(stub, assetKey, assetTypeFilter)
}

// Referrers returns an array of Keys of all the assets pointing to key.
// assetTypeFilter can be used to filter the results by asset type.
func (k Key) Referrers(stub *sw.StubWrapper, assetTypeFilter ...string) ([]Key, errors.ICCError) {
	assetKey := k.Key()
	return referrers(stub, assetKey, assetTypeFilter)
}

func referrers(stub *sw.StubWrapper, assetKey string, assetTypeFilter []string) ([]Key, errors.ICCError) {
	queryIt, err := stub.GetStateByPartialCompositeKey(assetKey, []string{})
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to check reference index", 500)
	}
	defer queryIt.Close()

	var retKeys []string
	for queryIt.HasNext() {
		ref, err := queryIt.Next()
		if err != nil {
			return nil, errors.WrapError(err, "failed to iterate in reference index")
		}

		newIndexState, isWritten := stub.WriteSet[ref.GetKey()]
		if isWritten && newIndexState == nil {
			continue
		}

		referredKey, keyParts, err := stub.Stub.SplitCompositeKey(ref.GetKey())
		if err != nil {
			return nil, errors.WrapError(err, "failed to split composite key")
		}

		if referredKey != assetKey || len(keyParts) == 0 {
			return nil, errors.WrapError(err, fmt.Sprintf("invalid reference index %s", ref.GetKey()))
		}

		retKeys = append(retKeys, keyParts[0])
	}

	for key, val := range stub.WriteSet {
		if len(key) > 0 && key[0] == 0x00 {
			referredKey, keyParts, err := stub.Stub.SplitCompositeKey(key)
			if err != nil {
				return nil, errors.WrapError(err, "failed to split composite key")
			}

			if referredKey != assetKey || len(keyParts) == 0 || !bytes.Equal(val, []byte{0x00}) {
				continue
			}

			isCounted := false
			for _, countedKey := range retKeys {
				if countedKey == keyParts[0] {
					isCounted = true
				}
			}
			if !isCounted {
				retKeys = append(retKeys, keyParts[0])
			}
		}
	}

	var ret []Key
	for _, retKey := range retKeys {
		assetType := strings.Split(retKey, ":")[0]
		if len(assetTypeFilter) <= 0 || contains(assetTypeFilter, assetType) {
			ret = append(ret, Key{
				"@assetType": strings.Split(retKey, ":")[0],
				"@key":       retKey,
			})
		}
	}

	return ret, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// validateRefs checks if subAsset references exist in blockchain.
func (a Asset) validateRefs(stub *sw.StubWrapper) errors.ICCError {
	// Fetch references contained in asset
	refKeys, err := a.Refs()
	if err != nil {
		return errors.WrapError(err, "failed to fetch references")
	}

	// Check if references exist
	for _, referencedKey := range refKeys {
		// Check if asset exists in blockchain
		assetJSON, err := referencedKey.Get(stub)
		if err != nil {
			return errors.WrapErrorWithStatus(err, "failed to read asset from blockchain", 400)
		}
		if assetJSON == nil {
			return errors.NewCCError("referenced asset not found", 404)
		}
	}
	return nil
}

func delRefs(stub *sw.StubWrapper, assetKey string, refKeys []Key) errors.ICCError {
	// Delete reference indexes
	for _, referencedKey := range refKeys {
		// Construct reference key
		indexKey, err := stub.CreateCompositeKey(referencedKey.Key(), []string{assetKey})
		if err != nil {
			return errors.WrapErrorWithStatus(err, "could not create composite key", 400)
		}
		// Check if asset exists in blockchain
		err = stub.DelState(indexKey)
		if err != nil {
			return errors.WrapErrorWithStatus(err, "failed to read asset from blockchain", 400)
		}
	}

	return nil
}

// delRefs deletes all the reference index for this asset from blockchain.
func (a Asset) delRefs(stub *sw.StubWrapper) errors.ICCError {
	// Fetch references contained in asset
	refKeys, err := a.Refs()
	if err != nil {
		return errors.WrapErrorWithStatus(err, "failed to fetch references", 400)
	}

	assetKey := a.Key()

	return delRefs(stub, assetKey, refKeys)
}

// delRefs deletes all the reference index for this asset from blockchain.
func (k Key) delRefs(stub *sw.StubWrapper) errors.ICCError {
	// Fetch references contained in asset
	refKeys, err := k.Refs(stub)
	if err != nil {
		return errors.WrapErrorWithStatus(err, "failed to fetch references", 400)
	}

	assetKey := k.Key()

	return delRefs(stub, assetKey, refKeys)
}

// putRefs writes the asset's reference index to the blockchain.
func (a Asset) putRefs(stub *sw.StubWrapper) errors.ICCError {
	// Fetch references contained in asset
	refKeys, err := a.Refs()
	if err != nil {
		return errors.WrapError(err, "failed to fetch references")
	}

	assetKey := a.Key()

	// Delete reference indexes
	for _, referencedKey := range refKeys {
		// Construct reference key
		refKey, err := stub.CreateCompositeKey(referencedKey.Key(), []string{assetKey})
		if err != nil {
			return errors.WrapErrorWithStatus(err, "failed generating composite key for reference", 500)
		}
		err = stub.PutState(refKey, []byte{0x00})
		if err != nil {
			return errors.WrapErrorWithStatus(err, "failed to put sub asset reference into blockchain", 500)
		}
	}

	return nil
}

// IsReferenced checks if the asset is referenced by another asset.
func (a Asset) IsReferenced(stub *sw.StubWrapper) (bool, errors.ICCError) {
	assetKey := a.Key()
	return isReferenced(stub, assetKey)
}

// IsReferenced checks if the asset is referenced by another asset.
func (k Key) IsReferenced(stub *sw.StubWrapper) (bool, errors.ICCError) {
	assetKey := k.Key()
	return isReferenced(stub, assetKey)
}

func isReferenced(stub *sw.StubWrapper, assetKey string) (bool, errors.ICCError) {
	queryIt, err := stub.GetStateByPartialCompositeKey(assetKey, []string{})
	if err != nil {
		return false, errors.WrapErrorWithStatus(err, "failed to check reference index", 500)
	}
	defer queryIt.Close()

	for queryIt.HasNext() {
		ref, err := queryIt.Next()
		if err != nil {
			return false, errors.WrapError(err, "failed to iterate in reference index")
		}

		newIndexState, isWritten := stub.WriteSet[ref.GetKey()]
		if !isWritten || newIndexState != nil {
			return true, nil
		}
	}

	for key, val := range stub.WriteSet {
		if strings.HasPrefix(key, assetKey) && key != assetKey && bytes.Equal(val, []byte{0x00}) {
			return true, nil
		}
	}

	return false, nil
}
