package test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/mock"
)

func TestLoadAssetTypeList(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	stubOrg2 := mock.NewMockStub("org2MSP", new(testCC))
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

	// Load List
	reqBytes, err = json.Marshal(map[string]interface{}{})
	if err != nil {
		t.FailNow()
	}
	res = stubOrg2.MockInvoke("loadAssetTypeList", [][]byte{
		[]byte("loadAssetTypeList"),
		reqBytes,
	})

	if res.GetStatus() != 200 {
		log.Println(res)
		t.FailNow()
	}

	assetTypeList := assets.AssetTypeList()
	if len(assetTypeList) != 6 {
		log.Println("Expected 6 asset types, got", len(assetTypeList))
		log.Println(assetTypeList)
		t.FailNow()
	}
}
