package assets

import (
	"fmt"
	"strings"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Refs returns an array of Keys containing the reference keys for all present subAssets.
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

// validateRefs checks if subAsset references exist in blockchain.
func (a Asset) validateRefs(stub shim.ChaincodeStubInterface) errors.ICCError {
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

// delRefs deletes all the reference index for this asset from blockchain.
func (a Asset) delRefs(stub shim.ChaincodeStubInterface) errors.ICCError {
	// Fetch references contained in asset
	refKeys, err := a.Refs()
	if err != nil {
		return errors.WrapErrorWithStatus(err, "failed to fetch references", 400)
	}

	assetKey := a.Key()

	// Delete reference indexes
	for _, referencedKey := range refKeys {
		// Construct reference key
		indexKey, err := stub.CreateCompositeKey(referencedKey.Key(), []string{assetKey})
		// Check if asset exists in blockchain
		err = stub.DelState(indexKey)
		if err != nil {
			return errors.WrapErrorWithStatus(err, "failed to read asset from blockchain", 400)
		}
	}

	return nil
}

// putRefs writes the asset's reference index to the blockchain.
func (a Asset) putRefs(stub shim.ChaincodeStubInterface) errors.ICCError {
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
func (a Asset) IsReferenced(stub shim.ChaincodeStubInterface) (bool, errors.ICCError) {
	// Get asset key
	assetKey := a.Key()
	queryIt, err := stub.GetStateByPartialCompositeKey(assetKey, []string{})
	if err != nil {
		return false, errors.WrapErrorWithStatus(err, "failed to check reference index", 500)
	}
	defer queryIt.Close()

	if queryIt.HasNext() {
		return true, nil
	}
	return false, nil
}
