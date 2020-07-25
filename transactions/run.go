package transactions

import (
	"fmt"
	"regexp"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
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

	// Verify callers permissions
	if tx.Callers != nil {
		// Get tx caller MSP ID
		txCaller, err := cid.GetMSPID(stub)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error getting tx caller", 500)
		}

		// Check if caller is allowed
		writePermission := false
		for _, w := range tx.Callers {
			match, err := regexp.MatchString(w, txCaller)
			if err != nil {
				return nil, errors.NewCCError("failed to check if caller matches regexp", 500)
			}
			if match {
				writePermission = true
			}
		}
		if !writePermission {
			return nil, errors.NewCCError(fmt.Sprintf("%s cannot call this transaction", txCaller), 403)
		}
	}

	return tx.Routine(stub, reqMap)
}
