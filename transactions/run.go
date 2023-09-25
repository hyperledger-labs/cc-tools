package transactions

import (
	"fmt"
	"regexp"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// Run defines the rules of transaction execution for the chaincode.
func Run(stub shim.ChaincodeStubInterface) ([]byte, errors.ICCError) {
	var err error

	// Extract the function and args from the transaction proposal
	txName, _ := stub.GetFunctionAndParameters()

	// Check if function exists
	tx := FetchTx(txName)
	if tx == nil {
		return nil, errors.NewCCError(fmt.Sprintf("tx named %s does not exist", txName), 400)
	}

	reqMap, err := tx.GetArgs(stub)
	if err != nil {
		return nil, errors.WrapError(err, "unable to get args")
	}

	sw := &sw.StubWrapper{
		Stub: stub,
	}

	if assets.GetEnabledDynamicAssetType() {
		err := assets.RestoreAssetList(sw, false)
		if err != nil {
			return nil, errors.WrapError(err, "failed to restore asset list")
		}
	}

	// Verify callers permissions
	if tx.Callers != nil {
		// Get tx caller MSP ID
		txCaller, err := sw.GetMSPID()
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error getting tx caller", 500)
		}

		// Check if caller is allowed
		callPermission := false
		for _, c := range tx.Callers {
			if len(c) <= 1 {
				continue
			}
			if c[0] == '$' { // if caller is regexp
				match, err := regexp.MatchString(c[1:], txCaller)
				if err != nil {
					return nil, errors.NewCCError("failed to check if caller matches regexp", 500)
				}
				if match {
					callPermission = true
					break
				}
			} else { // if caller is not regexp
				if c == txCaller {
					callPermission = true
					break
				}
			}
		}
		if !callPermission {
			return nil, errors.NewCCError(fmt.Sprintf("%s cannot call this transaction", txCaller), 403)
		}
	}

	return tx.Routine(sw, reqMap)
}
