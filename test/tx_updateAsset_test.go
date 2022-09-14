package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestUpdateAsset(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	person := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}
	setupState, _ := json.Marshal(person)

	stub.MockTransactionStart("setupUpdateAsset")
	err := stub.PutState("person:47061146-c642-51a1-844a-bf0b17cb5e19", setupState)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("setupUpdateAsset")

	personUpdate := map[string]interface{}{
		"@assetType":  "person",
		"name":        "Maria",
		"id":          "318.207.920-48",
		"dateOfBirth": "1999-05-06T22:12:41Z",
		"height":      1.66,
		"info":        map[string]interface{}{},
	}

	req := map[string]interface{}{
		"update": personUpdate,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}
	res := stub.MockInvoke("updateAsset", [][]byte{
		[]byte("updateAsset"),
		reqBytes,
	})
	lastUpdated, _ := stub.GetTxTimestamp()

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resPayload map[string]interface{}
	err = json.Unmarshal(res.GetPayload(), &resPayload)
	if err != nil {
		log.Println(resPayload)
		t.FailNow()
	}

	expectedPerson := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "updateAsset",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"dateOfBirth":  "1999-05-06T22:12:41Z",
		"height":       1.66,
	}

	if !reflect.DeepEqual(resPayload, expectedPerson) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resPayload)
		log.Printf("%#v\n", expectedPerson)
		t.FailNow()
	}

	var state map[string]interface{}
	err = json.Unmarshal(stub.State["person:47061146-c642-51a1-844a-bf0b17cb5e19"], &state)
	if err != nil {
		log.Println(resPayload)
		t.FailNow()
	}

	if !reflect.DeepEqual(state, expectedPerson) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resPayload)
		log.Printf("%#v\n", expectedPerson)
		t.FailNow()
	}
}
