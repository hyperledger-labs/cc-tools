package transactions

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func invokeAndVerify(stub *shim.MockStub, txName string, req, expectedRes interface{}, expectedStatus int32) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return err
	}

	res := stub.MockInvoke(txName, [][]byte{
		[]byte(txName),
		reqBytes,
	})

	if res.GetStatus() != expectedStatus {
		log.Println(res.Message)
		return fmt.Errorf("expected %d got %d", expectedStatus, res.GetStatus())
	}

	var resPayload interface{}
	err = json.Unmarshal(res.GetPayload(), &resPayload)
	if err != nil {
		log.Println(res.GetPayload())
		log.Println(err)
		return err
	}

	if !reflect.DeepEqual(resPayload, expectedRes) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resPayload)
		log.Printf("%#v\n", expectedRes)
		return fmt.Errorf("unexpected response")
	}

	return nil
}

func ensureEmpty(stub *shim.MockStub, key string) error {
	stub.MockTransactionStart("ensureDeletion")
	defer stub.MockTransactionEnd("ensureDeletion")
	state, err := stub.GetState(key)
	if err != nil {
		return fmt.Errorf("mock GetState error: %w", err)
	}
	if state != nil {
		return fmt.Errorf("key not deleted, state is %s", string(state))
	}
	return nil
}
