package test

import (
	"testing"

	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/cc-tools/mock"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	tx "github.com/hyperledger-labs/cc-tools/transactions"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// responseCC implements the shim.Chaincode interface using tx.RunResponse
type responseCC struct{}

func (t *responseCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *responseCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return tx.RunResponse(stub)
}

var responseTestTxs = []tx.Transaction{
	{
		Tag: "respondCreated",
		ResponseRoutine: func(stub *sw.StubWrapper, req map[string]interface{}) (*tx.Response, errors.ICCError) {
			return &tx.Response{Status: 201, Payload: []byte(`{"created":true}`)}, nil
		},
	},
	{
		Tag: "respondDefault",
		ResponseRoutine: func(stub *sw.StubWrapper, req map[string]interface{}) (*tx.Response, errors.ICCError) {
			return &tx.Response{Payload: []byte(`{}`)}, nil
		},
	},
	{
		Tag: "respondNil",
		ResponseRoutine: func(stub *sw.StubWrapper, req map[string]interface{}) (*tx.Response, errors.ICCError) {
			return nil, nil
		},
	},
	{
		Tag: "respondInvalid",
		ResponseRoutine: func(stub *sw.StubWrapper, req map[string]interface{}) (*tx.Response, errors.ICCError) {
			return &tx.Response{Status: 302}, nil
		},
	},
	{
		Tag: "respondError",
		ResponseRoutine: func(stub *sw.StubWrapper, req map[string]interface{}) (*tx.Response, errors.ICCError) {
			return nil, errors.NewCCError("asset not found", 404)
		},
	},
	{
		Tag: "routineError",
		Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
			return nil, errors.NewCCError("conflict", 409)
		},
	},
	{
		Tag: "routineErrorBelowThreshold",
		Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
			return nil, errors.NewCCError("wrong status", 300)
		},
	},
}

func initResponseTestTxs(t *testing.T) {
	txs := make([]tx.Transaction, 0, len(testTxList)+len(responseTestTxs))
	txs = append(txs, testTxList...)
	txs = append(txs, responseTestTxs...)
	tx.InitTxList(txs)
	t.Cleanup(func() { tx.InitTxList(testTxList) })
}

func TestRunResponse(t *testing.T) {
	initResponseTestTxs(t)
	stub := mock.NewMockStub("org1MSP", new(responseCC))

	cases := []struct {
		txName  string
		status  int32
		payload string
		message string
	}{
		{"respondCreated", 201, `{"created":true}`, ""},
		{"respondDefault", 200, `{}`, ""},
		{"respondNil", 200, "", ""},
		{"respondInvalid", 500, "", "invalid success status 302 returned by tx respondInvalid"},
		{"respondError", 404, "", "asset not found"},
		{"routineError", 409, "", "conflict"},
		{"routineErrorBelowThreshold", 500, "", "wrong status"},
		{"getHeader", 200, "", ""},
	}

	for _, c := range cases {
		res := stub.MockInvoke(c.txName, [][]byte{[]byte(c.txName)})
		if res.Status != c.status {
			t.Errorf("%s: expected status %d, got %d (%s)", c.txName, c.status, res.Status, res.Message)
		}
		if c.payload != "" && string(res.Payload) != c.payload {
			t.Errorf("%s: expected payload %s, got %s", c.txName, c.payload, res.Payload)
		}
		if res.Message != c.message {
			t.Errorf("%s: expected message %q, got %q", c.txName, c.message, res.Message)
		}
	}
}

func TestRunWithResponseRoutine(t *testing.T) {
	initResponseTestTxs(t)
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// Chaincodes that still call tx.Run keep responding 200 with the payload
	res := stub.MockInvoke("respondCreated", [][]byte{[]byte("respondCreated")})
	if res.Status != 200 {
		t.Errorf("expected status 200, got %d (%s)", res.Status, res.Message)
	}
	if string(res.Payload) != `{"created":true}` {
		t.Errorf("unexpected payload %s", res.Payload)
	}
}

func TestStartupCheckBothRoutines(t *testing.T) {
	routine := func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		return nil, nil
	}
	responseRoutine := func(stub *sw.StubWrapper, req map[string]interface{}) (*tx.Response, errors.ICCError) {
		return nil, nil
	}

	tx.InitTxList([]tx.Transaction{{Tag: "both", Routine: routine, ResponseRoutine: responseRoutine}})
	t.Cleanup(func() { tx.InitTxList(testTxList) })

	err := tx.StartupCheck()
	if err == nil || err.Status() != 500 {
		t.Fatalf("expected startup check error, got %v", err)
	}
}

func TestCCErrorString(t *testing.T) {
	err := errors.NewCCError("limit must be greater than 0", 400)
	if err.Error() != "400 limit must be greater than 0" {
		t.Errorf("unexpected error string %q", err.Error())
	}
}
