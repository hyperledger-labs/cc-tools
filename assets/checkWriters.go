package assets

import (
	"fmt"
	"regexp"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// CheckGlobalWriters checks if tx creator is allowed to write asset
func (a Asset) CheckGlobalWriters(stub shim.ChaincodeStubInterface) error {
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
	if assetTypeDef.Writers != nil {
		writePermission := false
		for _, w := range assetTypeDef.Writers {
			match, err := regexp.MatchString(w, txCreator)
			if err != nil {
				return errors.NewCCError("failed to check if writer matches regexp", 500)
			}
			if match {
				writePermission = true
			}
		}
		if !writePermission {
			return errors.NewCCError(fmt.Sprintf("%s cannot write to this asset", txCreator), 403)
		}
	}

	return nil
}

// CheckWriters checks if tx creator is allowed to write asset
func (a Asset) CheckWriters(stub shim.ChaincodeStubInterface) error {
	// Check full asset write permission
	err := a.CheckGlobalWriters(stub)
	if err != nil {
		return errors.WrapError(err, "failed writers check")
	}

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

	// Check attributes write permission
	for _, prop := range assetTypeDef.Props {
		if _, exists := a[prop.Tag]; exists && prop.Writers != nil {
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
	}

	return nil
}
