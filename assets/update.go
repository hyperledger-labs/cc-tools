package assets

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// Update receives a map[string]interface{} with key/vals to update the asset value in the world state.
func (a *Asset) Update(stub *sw.StubWrapper, update map[string]interface{}) (map[string]interface{}, errors.ICCError) {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return nil, errors.NewCCError(fmt.Sprintf("asset type named %s does not exist", a.TypeTag()), 400)
	}

	// Get tx creator MSP ID
	txCreator, err := stub.GetMSPID()
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "error getting tx creator", 500)
	}

	// Delete current reference indexes
	err = a.delRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed erasing old reference indexes from blockchain")
	}

	// Validate new asset properties
	for _, prop := range assetTypeDef.Props {
		// If prop is key, it cannot be updated
		if prop.IsKey {
			continue
		}

		// Check if property is included in the update map
		propInterface, propIncluded := update[prop.Tag]
		if !propIncluded || propInterface == nil {
			continue
		}

		if prop.ReadOnly {
			return nil, errors.NewCCError(fmt.Sprintf("cannot update asset property %s", prop.Label), 403)
		}

		// Check if tx creator is allowed to update this attribute
		if prop.Writers != nil {
			writePermission := false
			for _, w := range prop.Writers {
				if len(w) <= 1 {
					continue
				}
				if w[0] == '$' { // if writer is regexp
					match, err := regexp.MatchString(w[1:], txCreator)
					if err != nil {
						return nil, errors.NewCCError("failed to check if writer matches regexp", 500)
					}
					if match {
						writePermission = true
						break
					}
				} else { // if writer is not regexp
					if w == txCreator {
						writePermission = true
						break
					}
				}
			}
			if !writePermission {
				return nil, errors.NewCCError(fmt.Sprintf("%s cannot write to the '%s' (%s) asset property", txCreator, prop.Tag, prop.Label), 403)
			}
		}

		// Validate data types
		propInterface, err := validateProp(propInterface, prop)
		if err != nil {
			return nil, errors.WrapError(err, "error validating asset property")
		}

		(*a)[prop.Tag] = propInterface
	}

	err = a.validateRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed reference validation")
	}

	err = a.injectMetadata(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed injecting asset metadata")
	}

	ret, err := a.put(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed putting asset in ledger")
	}

	return ret, nil
}

// Update receives a map[string]interface{} with key/vals to update the asset value in the world state.
func (k *Key) Update(stub *sw.StubWrapper, update map[string]interface{}) (map[string]interface{}, errors.ICCError) {
	// Fetch asset properties
	assetTypeDef := k.Type()
	if assetTypeDef == nil {
		return nil, errors.NewCCError(fmt.Sprintf("asset type named %s does not exist", k.TypeTag()), 400)
	}

	// Get tx creator MSP ID
	txCreator, err := stub.GetMSPID()
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "error getting tx creator", 500)
	}

	// Delete current reference indexes
	err = k.delRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed erasing old reference indexes from blockchain")
	}

	assetMap, err := k.GetMap(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get asset current state")
	}

	// Validate new asset properties
	for _, prop := range assetTypeDef.Props {
		// If prop is key, it cannot be updated
		if prop.IsKey {
			continue
		}

		// Check if property is included in the update map
		propInterface, propIncluded := update[prop.Tag]
		if !propIncluded || propInterface == nil {
			continue
		}

		if prop.ReadOnly {
			return nil, errors.NewCCError(fmt.Sprintf("cannot update asset property %s", prop.Label), 403)
		}

		// Check if tx creator is allowed to update this attribute
		if prop.Writers != nil {
			writePermission := false
			for _, w := range prop.Writers {
				if len(w) <= 1 {
					continue
				}
				if w[0] == '$' { // if writer is regexp
					match, err := regexp.MatchString(w[1:], txCreator)
					if err != nil {
						return nil, errors.NewCCError("failed to check if writer matches regexp", 500)
					}
					if match {
						writePermission = true
						break
					}
				} else { // if writer is not regexp
					if w == txCreator {
						writePermission = true
						break
					}
				}
			}
			if !writePermission {
				return nil, errors.NewCCError(fmt.Sprintf("%s cannot write to the '%s' (%s) asset property", txCreator, prop.Tag, prop.Label), 403)
			}
		}

		// Validate data types
		propInterface, err := validateProp(propInterface, prop)
		if err != nil {
			return nil, errors.WrapError(err, "error validating asset property")
		}

		assetMap[prop.Tag] = propInterface
	}

	newAsset, err := NewAsset(assetMap)
	if err != nil {
		return nil, errors.WrapError(err, "could not construct asset object after update")
	}

	err = newAsset.validateRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed reference validation")
	}

	err = newAsset.injectMetadata(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed injecting asset metadata")
	}

	ret, err := newAsset.put(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed putting asset in ledger")
	}

	return ret, nil
}

// UpdateRecursive updates asset and all its subassets in blockchain.
// It checks if root asset and subassets exist, if not, it returns error.
// This method is experimental and might not work as intended. Use with caution.
func UpdateRecursive(stub *sw.StubWrapper, object map[string]interface{}) (map[string]interface{}, errors.ICCError) {
	objAsAsset, err := NewAsset(object)
	if err != nil {
		return nil, errors.WrapError(err, "unable to create asset object")
	}

	exists, err := objAsAsset.ExistsInLedger(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed checking if asset exists")
	}
	if !exists {
		return nil, errors.NewCCError("root asset not found", 404)
	}

	// Check if all assets and subassets exist in blockchain
	err = checkUpdateRecursive(stub, objAsAsset, true)
	if err != nil {
		return nil, errors.WrapError(err, "failed checking update recursive")
	}

	return PutRecursive(stub, object)
}

func checkUpdateRecursive(stub *sw.StubWrapper, object map[string]interface{}, root bool) errors.ICCError {
	var err error

	objAsKey, err := NewKey(object)
	if err != nil {
		return errors.WrapError(err, "unable to create asset object")
	}

	if !root {
		exists, err := objAsKey.ExistsInLedger(stub)
		if err != nil {
			return errors.WrapError(err, "failed checking if asset exists")
		}
		if !exists {
			return errors.NewCCError("subasset not found", 404)
		}

		object, err = objAsKey.GetMap(stub)
		if err != nil {
			return errors.WrapError(err, "failed getting subasset")
		}
	}

	objAsAsset, err := NewAsset(object)
	if err != nil {
		return errors.WrapError(err, "unable to create asset object")
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
				return errors.NewCCError(fmt.Sprintf("asset property %s must an array of type %s", subAsset.Label, subAsset.DataType), 400)
			}
		}

		for _, objInterface := range objArray {
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
				return errors.NewCCError(fmt.Sprintf("asset reference property '%s' must be an object", subAsset.Tag), 400)
			}
			obj["@assetType"] = dType
			err := checkUpdateRecursive(stub, obj, false)
			if err != nil {
				return errors.WrapError(err, fmt.Sprintf("failed to check sub-asset %s recursively", subAsset.Tag))
			}
		}

	}

	return nil
}
