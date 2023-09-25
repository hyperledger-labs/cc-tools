package assets

import (
	"fmt"
	"regexp"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// CheckWriters checks if tx creator is allowed to write asset.
func (a Asset) CheckWriters(stub *sw.StubWrapper) errors.ICCError {
	// Get tx creator MSP ID
	txCreator, err := stub.GetMSPID()
	if err != nil {
		return errors.WrapErrorWithStatus(err, "error getting tx creator", 500)
	}

	return a.checkWriters(txCreator)
}

// checkWriters is an internal function that checks if tx creator is allowed to write asset.
func (a Asset) checkWriters(txCreator string) errors.ICCError {
	// Fetch asset properties
	assetTypeDef := a.Type()
	if assetTypeDef == nil {
		return errors.NewCCError(fmt.Sprintf("asset type named %s does not exist", a.TypeTag()), 400)
	}

	// Check attributes write permission
	for _, prop := range assetTypeDef.Props {
		if _, exists := a[prop.Tag]; exists && prop.Writers != nil {
			writePermission := false
			for _, w := range prop.Writers {
				if len(w) <= 1 {
					continue
				}
				if w[0] == '$' { // if writer is regexp
					match, err := regexp.MatchString(w[1:], txCreator)
					if err != nil {
						return errors.NewCCError("failed to check if writer matches regexp", 500)
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
				return errors.NewCCError(fmt.Sprintf("%s cannot write to the '%s' (%s) asset property", txCreator, prop.Tag, prop.Label), 403)
			}
		}
	}

	return nil
}
