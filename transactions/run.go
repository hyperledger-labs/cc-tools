package transactions

import (
	"fmt"

	"github.com/hyperledger-labs/cc-tools/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// Run defines the rules of transaction execution for the chaincode.
// The success status set by a ResponseRoutine is discarded; use RunResponse to propagate it.
func Run(stub shim.ChaincodeStubInterface) ([]byte, errors.ICCError) {
	res, err := execute(stub)
	if err != nil {
		return nil, err
	}

	return res.Payload, nil
}

// RunResponse executes the transaction like Run and converts the result to a peer.Response,
// keeping the success status set by a ResponseRoutine and the status of any returned error.
func RunResponse(stub shim.ChaincodeStubInterface) pb.Response {
	res, err := execute(stub)
	if err != nil {
		return err.GetErrorResponse()
	}

	return pb.Response{
		Status:  res.Status,
		Payload: res.Payload,
	}
}

func execute(stub shim.ChaincodeStubInterface) (*Response, errors.ICCError) {
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

	if tx.ResponseRoutine == nil {
		payload, err := tx.Routine(sw, reqMap)
		if err != nil {
			return nil, err
		}
		return &Response{Status: 200, Payload: payload}, nil
	}

	res, ccErr := tx.ResponseRoutine(sw, reqMap)
	if ccErr != nil {
		return nil, ccErr
	}

	if res == nil {
		res = &Response{}
	}
	if res.Status == 0 {
		res.Status = 200
	}
	if res.Status < 200 || res.Status > 299 {
		return nil, errors.NewCCError(fmt.Sprintf("invalid success status %d returned by tx %s", res.Status, txName), 500)
	}

	return res, nil
}
