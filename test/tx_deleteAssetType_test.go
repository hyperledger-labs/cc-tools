package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestDeleteAssetType(t *testing.T) {
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

	// Delete Type
	deleteReq := map[string]interface{}{
		"tag":   "magazine",
		"force": true,
	}
	req = map[string]interface{}{
		"assetTypes": []map[string]interface{}{deleteReq},
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res = stub.MockInvoke("deleteAssetType", [][]byte{
		[]byte("deleteAssetType"),
		reqBytes,
	})
	expectedResponse := map[string]interface{}{
		"assetType": map[string]interface{}{
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
		},
	}

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	var resDeletePayload []map[string]interface{}
	err = json.Unmarshal(res.GetPayload(), &resDeletePayload)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if len(resDeletePayload) != 1 {
		log.Println("response length should be 1")
		t.FailNow()
	}

	if !reflect.DeepEqual(resDeletePayload[0], expectedResponse) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resDeletePayload[0])
		log.Printf("%#v\n", expectedResponse)
		t.FailNow()
	}
}

func TestDeleteAssetTypeEmptyList(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	req := map[string]interface{}{
		"assetTypes": []map[string]interface{}{},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("deleteAssetType", [][]byte{
		[]byte("deleteAssetType"),
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

func TestDeleteNonExistingAssetType(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	deleteTag := map[string]interface{}{
		"tag":   "inexistent",
		"force": true,
	}
	req := map[string]interface{}{
		"assetTypes": []map[string]interface{}{deleteTag},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}

	res := stub.MockInvoke("deleteAssetType", [][]byte{
		[]byte("deleteAssetType"),
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
