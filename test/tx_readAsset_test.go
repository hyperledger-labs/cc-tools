package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestReadAsset(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	expectedResponse := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}
	setupState, _ := json.Marshal(expectedResponse)
	stub.State["person:47061146-c642-51a1-844a-bf0b17cb5e19"] = setupState

	personKey := map[string]interface{}{
		"@assetType": "person",
		"name":       "Maria",
		"id":         "318.207.920-48",
	}

	req := map[string]interface{}{
		"key": personKey,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}
	res := stub.MockInvoke("readAsset", [][]byte{
		[]byte("readAsset"),
		reqBytes,
	})

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

	if !reflect.DeepEqual(resPayload, expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resPayload)
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}
}

func TestReadRecursive(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	setupPerson := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}
	setupBook := map[string]interface{}{
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org2MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "book",
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}
	setupPersonJSON, _ := json.Marshal(setupPerson)
	setupBookJSON, _ := json.Marshal(setupBook)

	stub.MockTransactionStart("setupReadAsset")
	stub.PutState("person:47061146-c642-51a1-844a-bf0b17cb5e19", setupPersonJSON)
	stub.PutState("book:a36a2920-c405-51c3-b584-dcd758338cb5", setupBookJSON)
	refIdx, err := stub.CreateCompositeKey("person:47061146-c642-51a1-844a-bf0b17cb5e19", []string{"book:a36a2920-c405-51c3-b584-dcd758338cb5"})
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	stub.PutState(refIdx, []byte{0x00})
	stub.MockTransactionEnd("setupReadAsset")

	bookKey := map[string]interface{}{
		"@assetType": "book",
		"@key":       "book:a36a2920-c405-51c3-b584-dcd758338cb5",
	}
	expectedResponse := map[string]interface{}{
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org2MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "book",
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			"@lastTouchBy": "org1MSP",
			"@lastTx":      "createAsset",
			"@assetType":   "person",
			"name":         "Maria",
			"id":           "31820792048",
			"height":       0.0,
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	req := map[string]interface{}{
		"key":     bookKey,
		"resolve": true,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}
	res := stub.MockInvoke("readAsset", [][]byte{
		[]byte("readAsset"),
		reqBytes,
	})

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

	if !reflect.DeepEqual(resPayload, expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resPayload)
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}
}
