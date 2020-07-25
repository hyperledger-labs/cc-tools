package assets

import (
	"fmt"
	"regexp"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Update receives a map[string]interface{} with key/vals to update in asset
func (a *Asset) Update(stub shim.ChaincodeStubInterface, update map[string]interface{}) error {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return errors.NewCCError(fmt.Sprintf("asset type named %s does not exist", a.TypeTag()), 400)
	}

	// Get tx creator MSP ID
	txCreator, err := cid.GetMSPID(stub)
	if err != nil {
		return errors.WrapErrorWithStatus(err, "error getting tx creator", 500)
	}

	// Check full asset write permission
	err = a.CheckGlobalWriters(stub)
	if err != nil {
		return errors.WrapError(err, "failed writers check")
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
			return errors.NewCCError(fmt.Sprintf("cannot update asset property %s", prop.Label), 403)
		}

		// Check if tx creator is allowed to update this attribute
		if prop.Writers != nil {
			writePermission := false
			for _, w := range prop.Writers {
				match, err := regexp.MatchString(w, txCreator)
				if err != nil {
					return errors.NewCCError("failed to check if writer matches regexp", 500)
				}
				if match {
					writePermission = true
				}
			}
			if !writePermission {
				return errors.NewCCError(fmt.Sprintf("%s cannot write to this asset property", txCreator), 403)
			}
		}

		// Validate data types
		propInterface, err := validateProp(propInterface, prop)
		if err != nil {
			return errors.WrapError(err, "error validating asset property")
		}

		(*a)[prop.Tag] = propInterface
	}
	return nil
}
