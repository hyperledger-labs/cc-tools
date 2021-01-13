package assets

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// testCC implements the shim.Chaincode interface
type testCC struct{}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *testCC) Init(stub shim.ChaincodeStubInterface) (response pb.Response) {
	err := StartupCheck()
	if err != nil {
		response = err.GetErrorResponse()
		return
	}

	// Get the args from the transaction proposal
	args := stub.GetStringArgs()

	// Test if argument list is empty
	if len(args) != 1 {
		response = shim.Error("the Init method expects 1 argument")
		response.Status = 400
		return
	}

	// Test if argument is "init" or "upgrade". Fails otherwise.
	if args[0] != "init" && args[0] != "upgrade" {
		response = shim.Error("the argument should be init or upgrade (as sent by Node.js SDK)")
		response.Status = 400
		return
	}

	response = shim.Success(nil)
	return
}

// Invoke is called per transaction on the chaincode.
func (t *testCC) Invoke(stub shim.ChaincodeStubInterface) (response pb.Response) {
	response = shim.Success([]byte(""))
	return
}
