package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestCreateAsset(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	person := map[string]interface{}{
		"@assetType": "person",
		"name":       "Maria",
		"id":         "318.207.920-48",
		"info": map[string]interface{}{
			"passport": "1234",
		},
	}
	req := map[string]interface{}{
		"asset": []map[string]interface{}{person},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("createAsset", [][]byte{
		[]byte("createAsset"),
		reqBytes,
	})
	lastUpdated, _ := stub.GetTxTimestamp()
	expectedResponse := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
		"info": map[string]interface{}{
			"@assetType": "@object",
			"passport":   "1234",
		},
	}

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resPayload []map[string]interface{}
	err = json.Unmarshal(res.GetPayload(), &resPayload)
	if err != nil {
		log.Println(err)
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

	var state map[string]interface{}
	stateBytes := stub.State["person:47061146-c642-51a1-844a-bf0b17cb5e19"]
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(state, expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", state)
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}
}

func TestCreateAssetEmptyList(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	req := map[string]interface{}{
		"asset": []map[string]interface{}{},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("CreateAsset", [][]byte{
		[]byte("createAsset"),
		reqBytes,
	})

	if res.GetStatus() != 400 {
		log.Println(res)
		t.FailNow()
	}

	if res.GetMessage() != "unable to get args: required argument 'asset' must be non-empty" {
		log.Printf("error message different from expected: %s", res.GetMessage())
		t.FailNow()
	}
}

func TestCreatePrivate(t *testing.T) {
	stub := mock.NewMockStub("org2MSP", new(testCC))
	secret := map[string]interface{}{
		"@assetType": "secret",
		"secretName": "testSecret",
		"secret":     "this is very secret",
	}
	expectedResponse := map[string]interface{}{
		"@key":       "secret:73a3f9a7-eb91-5f4d-b1bb-c0487e90f40b",
		"@assetType": "secret",
	}
	req := map[string]interface{}{
		"asset": []map[string]interface{}{secret},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("createAsset", [][]byte{
		[]byte("createAsset"),
		reqBytes,
	})
	lastUpdated, _ := stub.GetTxTimestamp()
	expectedState := map[string]interface{}{
		"@key":         "secret:73a3f9a7-eb91-5f4d-b1bb-c0487e90f40b",
		"@assetType":   "secret",
		"@lastTouchBy": "org2MSP",
		"@lastTx":      "createAsset",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"secretName":   "testSecret",
		"secret":       "this is very secret",
	}

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resPayload []map[string]interface{}
	err = json.Unmarshal(res.GetPayload(), &resPayload)
	if err != nil {
		log.Println(err)
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

	var state map[string]interface{}
	stateBytes := stub.PvtState["secret"]["secret:73a3f9a7-eb91-5f4d-b1bb-c0487e90f40b"]
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(state, expectedState) {
		log.Println("these should be equal")
		log.Printf("%#v\n", state)
		log.Printf("%#v\n", expectedState)
		t.FailNow()
	}
}
