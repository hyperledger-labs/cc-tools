package transactions

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestGetArgs(t *testing.T) {
	stub := shim.NewMockStub("testcc", new(testCC))
	res := stub.MockInvoke("TestGetArgs", [][]byte{[]byte("getTx")})
	if res.GetStatus() != 200 {
		fmt.Println(res)
		t.FailNow()
	}
	res = stub.MockInvoke("TestGetArgs", [][]byte{[]byte("getTx"), []byte("{\"txName\": \"getTx\"}")})
	if res.GetStatus() != 200 {
		fmt.Println(res)
		t.FailNow()
	}
}
