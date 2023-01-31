package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestCreateAssetType(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	newType := map[string]interface{}{
		"tag":         "magazine",
		"label":       "Magazine",
		"description": "Magazine definition",
		"props": []map[string]interface{}{
			{
				"tag":      "name",
				"label":    "Name",
				"dataType": "string",
				"required": true,
				"writers":  []string{"org1MSP"},
				"isKey":    true,
			},
			{
				"tag":      "images",
				"label":    "Images",
				"dataType": "[]string",
			},
		},
	}
	req := map[string]interface{}{
		"assetTypes": []map[string]interface{}{newType},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("createAssetType", [][]byte{
		[]byte("createAssetType"),
		reqBytes,
	})
	expectedResponse := map[string]interface{}{
		"description": "Magazine definition",
		"dynamic":     true,
		"label":       "Magazine",
		"props": []interface{}{
			map[string]interface{}{
				"dataType":    "string",
				"description": "",
				"isKey":       true,
				"label":       "Name",
				"readOnly":    false,
				"required":    true,
				"tag":         "name",
				"writers":     []interface{}{"org1MSP"},
			},
			map[string]interface{}{
				"dataType":    "[]string",
				"description": "",
				"isKey":       false,
				"label":       "Images",
				"readOnly":    false,
				"required":    false,
				"tag":         "images",
				"writers":     nil,
			},
		},
		"tag": "magazine",
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

	// Create Asset
	asset := map[string]interface{}{
		"@assetType": "magazine",
		"name":       "MAG",
		"images":     []string{"url.com/1", "url.com/2"},
	}
	req = map[string]interface{}{
		"asset": []map[string]interface{}{asset},
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res = stub.MockInvoke("createAsset", [][]byte{
		[]byte("createAsset"),
		reqBytes,
	})
	lastUpdated, _ := stub.GetTxTimestamp()
	expectedResponse = map[string]interface{}{
		"@key":         "magazine:236a29db-f53c-59e1-ac6d-a4f264dbc477",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"@assetType":   "magazine",
		"name":         "MAG",
		"images": []interface{}{
			"url.com/1",
			"url.com/2",
		},
	}

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resAssetPayload []map[string]interface{}
	err = json.Unmarshal(res.GetPayload(), &resAssetPayload)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if len(resAssetPayload) != 1 {
		log.Println("response length should be 1")
		t.FailNow()
	}

	if !reflect.DeepEqual(resAssetPayload[0], expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resAssetPayload[0])
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}

	var state map[string]interface{}
	stateBytes := stub.State["magazine:236a29db-f53c-59e1-ac6d-a4f264dbc477"]
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

func TestCreateAssetTypeEmptyList(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	req := map[string]interface{}{
		"assetTypes": []map[string]interface{}{},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("createAssetType", [][]byte{
		[]byte("createAssetType"),
		reqBytes,
	})

	if res.GetStatus() != 400 {
		log.Println(res)
		t.FailNow()
	}

	if res.GetMessage() != "unable to get args: required argument 'assetTypes' must be non-empty" {
		log.Printf("error message different from expected: %s", res.GetMessage())
		t.FailNow()
	}
}

func TestCreateExistingAssetType(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	newType := map[string]interface{}{
		"tag":         "library",
		"label":       "Library",
		"description": "Library definition",
		"props": []map[string]interface{}{
			{
				"tag":      "name",
				"label":    "Name",
				"dataType": "string",
				"required": true,
				"isKey":    true,
			},
		},
	}
	req := map[string]interface{}{
		"assetTypes": []map[string]interface{}{newType},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("createAssetType", [][]byte{
		[]byte("createAssetType"),
		reqBytes,
	})

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

	if len(resPayload) != 0 {
		log.Println("response length should be 0")
		t.FailNow()
	}
}

func TestCreateAssetTypeWithoutKey(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	newType := map[string]interface{}{
		"tag":         "library",
		"label":       "Library",
		"description": "Library definition",
		"props": []map[string]interface{}{
			{
				"tag":      "name",
				"label":    "Name",
				"dataType": "string",
				"required": true,
			},
		},
	}
	req := map[string]interface{}{
		"assetTypes": []map[string]interface{}{newType},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("createAssetType", [][]byte{
		[]byte("createAssetType"),
		reqBytes,
	})

	if res.GetStatus() != 400 {
		log.Println(res)
		t.FailNow()
	}

	if res.GetMessage() != "failed to build asset type: asset type must have a key" {
		log.Printf("error message different from expected: %s", res.GetMessage())
		t.FailNow()
	}
}
