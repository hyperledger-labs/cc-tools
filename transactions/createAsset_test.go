package transactions

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestCreateAsset(t *testing.T) {
	stub := shim.NewMockStub("org1MSP", new(testCC))
	person := map[string]interface{}{
		"@assetType": "person",
		"name":       "Maria",
		"id":         "318.207.920-48",
	}
	expectedResponse := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}
	req := map[string]interface{}{
		"asset": []map[string]interface{}{person},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("CreateAsset", [][]byte{
		[]byte("createAsset"),
		reqBytes,
	})

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resPayload []map[string]interface{}
	err = json.Unmarshal(res.GetPayload(), &resPayload)
	if err != nil {
		log.Println(resPayload)
		t.FailNow()
	}

	if len(resPayload) != 1 {
		log.Println("response length should be 1")
		t.FailNow()
	}

	if !reflect.DeepEqual(resPayload[0], expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resPayload[0])
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}
}
