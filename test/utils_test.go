package test

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/goledgerdev/cc-tools/mock"
)

func invokeAndVerify(stub *mock.MockStub, txName string, req, expectedRes interface{}, expectedStatus int32) error {
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

	resData = clearLastUpdated(resData)
	if !reflect.DeepEqual(resData, expectedRes) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resData)
		log.Printf("%#v\n", expectedRes)
		return fmt.Errorf("unexpected response")
	}

	return nil
}

func isEmpty(stub *mock.MockStub, key string) bool {
	state := stub.State[key]
	return state == nil
}

// This is done like this because invokeAndVerify does not allow
// us to access tx timestamp before calling it. A refactor is
// recommended but not urgent.
func clearLastUpdated(in interface{}) interface{} {
	var out interface{}
	switch input := in.(type) {
	case map[string]interface{}:
		delete(input, "@lastUpdated")
		for k := range input {
			input[k] = clearLastUpdated(input[k])
		}
		out = input
	case []interface{}:
		for k := range input {
			input[k] = clearLastUpdated(input[k])
		}
		out = input
	default:
		out = input
	}
	return out
}
