package test

import (
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
)

func TestAssetPropToMap(t *testing.T) {
	propMap := testAssetList[0].GetPropDef("id").ToMap()
	expectedMap := map[string]interface{}{
		"tag":          "id",
		"label":        "CPF (Brazilian ID)",
		"description":  "",
		"isKey":        true,
		"required":     true,
		"readOnly":     false,
		"defaultValue": nil,
		"dataType":     "cpf",
		"writers":      []string{"org1MSP"},
	}

	if !reflect.DeepEqual(propMap, expectedMap) {
		log.Println("these should be deeply equal")
		log.Println(propMap)
		log.Println(expectedMap)
		t.FailNow()
	}
}

func TestAssetPropFromMap(t *testing.T) {
	testMap := map[string]interface{}{
		"tag":      "secretName",
		"isKey":    true,
		"label":    "Secret Name",
		"dataType": "string",
		"writers":  []interface{}{"org2MSP"},
	}
	testAssetProp := assets.AssetPropFromMap(testMap)
	expectedProp := *testAssetList[3].GetPropDef("secretName")

	if !reflect.DeepEqual(testAssetProp, expectedProp) {
		log.Println("these should be deeply equal")
		log.Println(testAssetProp)
		log.Println(expectedProp)
		t.FailNow()
	}
}
