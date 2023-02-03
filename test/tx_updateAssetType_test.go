package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestUpdateAssetType(t *testing.T) {
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

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	// Update Type
	updateReq := map[string]interface{}{
		"tag":   "magazine",
		"label": "Magazines",
		"props": []map[string]interface{}{
			{
				"tag":    "images",
				"delete": true,
			},
			{
				"tag":   "name",
				"label": "Magazine Name",
			},
			{
				"tag":      "pages",
				"label":    "Pages",
				"dataType": "[]string",
			},
		},
	}
	req = map[string]interface{}{
		"assetTypes":               []map[string]interface{}{updateReq},
		"skipAssetEmptyValidation": true,
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res = stub.MockInvoke("updateAssetType", [][]byte{
		[]byte("updateAssetType"),
		reqBytes,
	})
	expectedResponse := map[string]interface{}{
		"description": "Magazine definition",
		"dynamic":     true,
		"label":       "Magazines",
		"props": []interface{}{
			map[string]interface{}{
				"dataType":    "string",
				"description": "",
				"isKey":       true,
				"label":       "Magazine Name",
				"readOnly":    false,
				"required":    true,
				"tag":         "name",
				"writers":     []interface{}{"org1MSP"},
			},
			map[string]interface{}{
				"dataType":    "[]string",
				"description": "",
				"isKey":       false,
				"label":       "Pages",
				"readOnly":    false,
				"required":    false,
				"tag":         "pages",
				"writers":     nil,
			},
		},
		"tag": "magazine",
	}

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resUpdatePayload map[string]interface{}
	err = json.Unmarshal(res.GetPayload(), &resUpdatePayload)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	assetTypesPayload := resUpdatePayload["assetTypes"].([]interface{})

	if len(assetTypesPayload) != 1 {
		log.Println("response length should be 1")
		t.FailNow()
	}

	if !reflect.DeepEqual(assetTypesPayload[0].(map[string]interface{}), expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", assetTypesPayload[0].(map[string]interface{}))
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}
}

func TestUpdateAssetTypeEmptyList(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	req := map[string]interface{}{
		"assetTypes": []map[string]interface{}{},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("updateAssetType", [][]byte{
		[]byte("updateAssetType"),
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

func TestUpdateNonExistingAssetType(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	updateTag := map[string]interface{}{
		"tag":   "inexistent",
		"label": "New Label",
	}
	req := map[string]interface{}{
		"assetTypes": []map[string]interface{}{updateTag},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("updateAssetType", [][]byte{
		[]byte("updateAssetType"),
		reqBytes,
	})

	if res.GetStatus() != 400 {
		log.Println(res)
		t.FailNow()
	}

	if res.GetMessage() != "asset type 'inexistent' not found" {
		log.Printf("error message different from expected: %s", res.GetMessage())
		t.FailNow()
	}
}

func TestDeleteNonExistingProp(t *testing.T) {
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

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	// Update prop
	updateReq := map[string]interface{}{
		"tag":   "magazine",
		"label": "Magazines",
		"props": []map[string]interface{}{
			{
				"tag":    "inexistant",
				"delete": true,
			},
		},
	}
	req = map[string]interface{}{
		"assetTypes":               []map[string]interface{}{updateReq},
		"skipAssetEmptyValidation": true,
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res = stub.MockInvoke("updateAssetType", [][]byte{
		[]byte("updateAssetType"),
		reqBytes,
	})

	if res.GetStatus() != 400 {
		log.Println(res)
		t.FailNow()
	}

	if res.GetMessage() != "invalid props array: attempt to delete inexistent prop" {
		log.Printf("error message different from expected: %s", res.GetMessage())
		t.FailNow()
	}
}

func TestDeleteKeyProp(t *testing.T) {
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

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	// Update prop
	updateReq := map[string]interface{}{
		"tag":   "magazine",
		"label": "Magazines",
		"props": []map[string]interface{}{
			{
				"tag":    "name",
				"delete": true,
			},
		},
	}
	req = map[string]interface{}{
		"assetTypes":               []map[string]interface{}{updateReq},
		"skipAssetEmptyValidation": true,
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res = stub.MockInvoke("updateAssetType", [][]byte{
		[]byte("updateAssetType"),
		reqBytes,
	})

	if res.GetStatus() != 400 {
		log.Println(res)
		t.FailNow()
	}

	if res.GetMessage() != "invalid props array: cannot delete key prop" {
		log.Printf("error message different from expected: %s", res.GetMessage())
		t.FailNow()
	}
}

func TestAttemptToUpdateInvalidPropSpecs(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	newType := map[string]interface{}{
		"tag":         "rack",
		"label":       "Rack",
		"description": "Rack definition",
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

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	// Update prop
	updateReq := map[string]interface{}{
		"tag": "rack",
		"props": []map[string]interface{}{
			{
				"tag":      "images",
				"dataType": "[]string",
				"isKey":    true,
			},
		},
	}
	req = map[string]interface{}{
		"assetTypes":               []map[string]interface{}{updateReq},
		"skipAssetEmptyValidation": true,
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res = stub.MockInvoke("updateAssetType", [][]byte{
		[]byte("updateAssetType"),
		reqBytes,
	})
	expectedResponse := map[string]interface{}{
		"description": "Rack definition",
		"dynamic":     true,
		"label":       "Rack",
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
		"tag": "rack",
	}

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resUpdatePayload map[string]interface{}
	err = json.Unmarshal(res.GetPayload(), &resUpdatePayload)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	assetTypesPayload := resUpdatePayload["assetTypes"].([]interface{})

	if len(assetTypesPayload) != 1 {
		log.Println("response length should be 1")
		t.FailNow()
	}

	if !reflect.DeepEqual(assetTypesPayload[0].(map[string]interface{}), expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", assetTypesPayload[0].(map[string]interface{}))
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}
}

func TestCreateKeyProp(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	newType := map[string]interface{}{
		"tag":         "page",
		"label":       "Page",
		"description": "Page definition",
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

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	// Update prop
	updateReq := map[string]interface{}{
		"tag": "page",
		"props": []map[string]interface{}{
			{
				"tag":      "index",
				"label":    "Index",
				"dataType": "[]number",
				"isKey":    true,
			},
		},
	}
	req = map[string]interface{}{
		"assetTypes":               []map[string]interface{}{updateReq},
		"skipAssetEmptyValidation": true,
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res = stub.MockInvoke("updateAssetType", [][]byte{
		[]byte("updateAssetType"),
		reqBytes,
	})

	if res.GetStatus() != 400 {
		log.Println(res)
		t.FailNow()
	}

	if res.GetMessage() != "invalid props array: cannot create key prop" {
		log.Printf("error message different from expected: %s", res.GetMessage())
		t.FailNow()
	}
}
