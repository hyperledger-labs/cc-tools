package test

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func invokeAndVerify(stub *shimtest.MockStub, txName string, req, expectedRes interface{}, expectedStatus int32) error {
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
		log.Println(res.GetMessage())
		return fmt.Errorf("expected %d got %d", expectedStatus, res.GetStatus())
	}

	var resData interface{}
	if expectedStatus == 200 {
		err = json.Unmarshal(res.GetPayload(), &resData)
	} else {
		resData = res.GetMessage()
	}
	if err != nil {
		log.Println(err)
		return err
	}
	if !reflect.DeepEqual(resData, expectedRes) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resData)
		log.Printf("%#v\n", expectedRes)
		return fmt.Errorf("unexpected response")
	}

	return nil
}

func isEmpty(stub *shimtest.MockStub, key string) bool {
	stub.MockTransactionStart("ensureDeletion")
	defer stub.MockTransactionEnd("ensureDeletion")
	state := stub.State[key]
	return state == nil
}
