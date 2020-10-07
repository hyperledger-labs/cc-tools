package assets

import (
	"fmt"
	"regexp"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Update receives a map[string]interface{} with key/vals to update in asset
func (a *Asset) Update(stub shim.ChaincodeStubInterface, update map[string]interface{}) (map[string]interface{}, error) {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return nil, errors.NewCCError(fmt.Sprintf("asset type named %s does not exist", a.TypeTag()), 400)
	}

	// Get tx creator MSP ID
	txCreator, err := cid.GetMSPID(stub)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "error getting tx creator", 500)
	}

	// Delete current reference indexes
	err = a.delRefs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed erasing old reference indexes from ledger")
	}

	// Delete current unique markers
	err = a.delUniqueMarkers(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed erasing old unique markers from ledger")
	}

	// Validate new asset properties
	for _, prop := range assetTypeDef.Props {
		// If prop is key, it cannot be updated
		if prop.IsKey {
			continue
		}

		// Check if property is included in the update map
		propInterface, propIncluded := update[prop.Tag]
		if !propIncluded {
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

	err = a.checkUniqueMarkers(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed unique props validation")
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
