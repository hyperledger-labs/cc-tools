package transactions

import (
	"fmt"

	"github.com/hyperledger-labs/cc-tools/v2/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/v2/assets"
	"github.com/hyperledger-labs/cc-tools/v2/errors"
	sw "github.com/hyperledger-labs/cc-tools/v2/stubwrapper"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	pb "github.com/hyperledger/fabric-protos-go-apiv2/peer"
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
func RunResponse(stub shim.ChaincodeStubInterface) *pb.Response {
	res, err := execute(stub)
	if err != nil {
		// Other ICCError implementations may keep statuses below 400, which peers would endorse
		errRes := err.GetErrorResponse()
		if errRes == nil {
			errRes = &pb.Response{Message: err.Error()}
		}
		if errRes.Status < 400 {
			errRes.Status = 500
		}
		return errRes
	}

	return &pb.Response{
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

	// Copy the response so the value returned by the routine is not modified
	out := Response{}
	if res != nil {
		out = *res
	}
	if out.Status == 0 {
		out.Status = 200
	}
	if out.Status < 200 || out.Status > 299 {
		return nil, errors.NewCCError(fmt.Sprintf("invalid success status %d returned by tx %s, expected a 2xx status", out.Status, txName), 500)
	}

	return &out, nil
}
