package transactions

import (
	"fmt"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
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
	callPermission, err := accesscontrol.AllowCaller(stub, tx.Callers)
	if err != nil {
		return nil, errors.WrapError(err, "failed to check permissions")
	}

	if !callPermission {
		return nil, errors.NewCCError("current caller not allowed", 403)
	}

	return tx.Routine(sw, reqMap)
}
