package test

import (
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
)

func TestGetPropDef(t *testing.T) {
	propDef := *testAssetList[0].GetPropDef("id")
	expectedPropDef := testAssetList[0].Props[0]

	if !reflect.DeepEqual(propDef, expectedPropDef) {
		log.Println("these should be deeply equal")
		log.Println(propDef)
		log.Println(expectedPropDef)
		t.FailNow()
	}
}

func TestAssetTypeToMap(t *testing.T) {
	var emptySlice []string
	typeMap := testAssetList[0].ToMap()
	expectedMap := map[string]interface{}{
		"tag":         "person",
		"label":       "Person",
		"description": "Personal data of someone",
		"props": []map[string]interface{}{
			{
				"tag":          "id",
				"label":        "CPF (Brazilian ID)",
				"description":  "",
				"isKey":        true,
				"required":     true,
				"readOnly":     false,
				"defaultValue": nil,
				"dataType":     "cpf",
				"writers":      []string{"org1MSP"},
			},
			{
				"tag":          "name",
				"label":        "Name of the person",
				"description":  "",
				"isKey":        false,
				"required":     true,
				"readOnly":     false,
				"defaultValue": nil,
				"dataType":     "string",
				"writers":      emptySlice,
			},
			{
				"tag":          "dateOfBirth",
				"label":        "Date of Birth",
				"description":  "",
				"isKey":        false,
				"required":     false,
				"readOnly":     false,
				"defaultValue": nil,
				"dataType":     "datetime",
				"writers":      []string{"org1MSP"},
			},
			{
				"tag":          "height",
				"label":        "Person's height",
				"description":  "",
				"isKey":        false,
				"required":     false,
				"readOnly":     false,
				"defaultValue": 0,
				"dataType":     "number",
				"writers":      emptySlice,
			},
			{
				"tag":          "info",
				"label":        "Other Info",
				"description":  "",
				"isKey":        false,
				"required":     false,
				"readOnly":     false,
				"defaultValue": nil,
				"dataType":     "@object",
				"writers":      emptySlice,
			},
		},
		"readers": emptySlice,
		"dynamic": false,
	}

	if !reflect.DeepEqual(typeMap, expectedMap) {
		log.Println("these should be deeply equal")
		log.Println(typeMap)
		log.Println(expectedMap)
		t.FailNow()
	}
}

func TestAssetTypeFromMap(t *testing.T) {
	testMap := map[string]interface{}{
		"tag":         "secret",
		"label":       "Secret",
		"description": "Secret between Org2 and Org3",
		"props": []interface{}{
			map[string]interface{}{
				"tag":      "secretName",
				"isKey":    true,
				"label":    "Secret Name",
				"dataType": "string",
				"writers":  []interface{}{"org2MSP"},
			},
			map[string]interface{}{
				"tag":      "secret",
				"label":    "Secret",
				"dataType": "string",
				"required": true,
			},
		},
		"readers": []interface{}{"org2MSP", "org3MSP"},
	}
	testAssetType := assets.AssetTypeFromMap(testMap)
	expectedType := testAssetList[3]

	if !reflect.DeepEqual(testAssetType, expectedType) {
		log.Println("these should be deeply equal")
		log.Println(testAssetType)
		log.Println(expectedType)
		t.FailNow()
	}
}

func TestAssetTypeListToMap(t *testing.T) {
	assetList := []assets.AssetType{
		testAssetList[0],
		testAssetList[3],
	}

	mapList := assets.ArrayFromAssetTypeList(assetList)
	var emptySlice []string
	expectedMap := []map[string]interface{}{
		{
			"tag":         "person",
			"label":       "Person",
			"description": "Personal data of someone",
			"props": []map[string]interface{}{
				{
					"tag":          "id",
					"label":        "CPF (Brazilian ID)",
					"description":  "",
					"isKey":        true,
					"required":     true,
					"readOnly":     false,
					"defaultValue": nil,
					"dataType":     "cpf",
					"writers":      []string{"org1MSP"},
				},
				{
					"tag":          "name",
					"label":        "Name of the person",
					"description":  "",
					"isKey":        false,
					"required":     true,
					"readOnly":     false,
					"defaultValue": nil,
					"dataType":     "string",
					"writers":      emptySlice,
				},
				{
					"tag":          "dateOfBirth",
					"label":        "Date of Birth",
					"description":  "",
					"isKey":        false,
					"required":     false,
					"readOnly":     false,
					"defaultValue": nil,
					"dataType":     "datetime",
					"writers":      []string{"org1MSP"},
				},
				{
					"tag":          "height",
					"label":        "Person's height",
					"description":  "",
					"isKey":        false,
					"required":     false,
					"readOnly":     false,
					"defaultValue": 0,
					"dataType":     "number",
					"writers":      emptySlice,
				},
				{
					"tag":          "info",
					"label":        "Other Info",
					"description":  "",
					"isKey":        false,
					"required":     false,
					"readOnly":     false,
					"defaultValue": nil,
					"dataType":     "@object",
					"writers":      emptySlice,
				},
			},
			"readers": emptySlice,
			"dynamic": false,
		},
		{
			"tag":         "secret",
			"label":       "Secret",
			"description": "Secret between Org2 and Org3",
			"props": []map[string]interface{}{
				{
					"tag":          "secretName",
					"label":        "Secret Name",
					"description":  "",
					"isKey":        true,
					"required":     false,
					"readOnly":     false,
					"defaultValue": nil,
					"dataType":     "string",
					"writers":      []string{"org2MSP"},
				},
				{
					"tag":          "secret",
					"label":        "Secret",
					"description":  "",
					"isKey":        false,
					"required":     true,
					"readOnly":     false,
					"defaultValue": nil,
					"dataType":     "string",
					"writers":      emptySlice,
				},
			},
			"readers": []string{"org2MSP", "org3MSP"},
			"dynamic": false,
		},
	}

	if !reflect.DeepEqual(mapList, expectedMap) {
		log.Println("these should be deeply equal")
		log.Println(mapList)
		log.Println(expectedMap)
		t.FailNow()
	}
}

func TestArrayFromAssetTypeList(t *testing.T) {
	testArray := []interface{}{
		map[string]interface{}{
			"tag":         "secret",
			"label":       "Secret",
			"description": "Secret between Org2 and Org3",
			"props": []interface{}{
				map[string]interface{}{
					"tag":      "secretName",
					"isKey":    true,
					"label":    "Secret Name",
					"dataType": "string",
					"writers":  []interface{}{"org2MSP"},
				},
				map[string]interface{}{
					"tag":      "secret",
					"label":    "Secret",
					"dataType": "string",
					"required": true,
				},
			},
			"readers": []interface{}{"org2MSP", "org3MSP"},
		},
		map[string]interface{}{
			"tag":         "library",
			"label":       "Library",
			"description": "Library as a collection of books",
			"props": []interface{}{
				map[string]interface{}{
					"tag":      "name",
					"isKey":    true,
					"required": true,
					"label":    "Library Name",
					"dataType": "string",
					"writers":  []interface{}{"org3MSP"},
				},
				map[string]interface{}{
					"tag":      "books",
					"label":    "Book Collection",
					"dataType": "[]->book",
				},
				map[string]interface{}{
					"tag":      "entranceCode",
					"label":    "Entrance Code for the Library",
					"dataType": "->secret",
				},
			},
		},
	}
	array := assets.AssetTypeListFromArray(testArray)
	expectedList := []assets.AssetType{
		testAssetList[3],
		testAssetList[1],
	}

	if !reflect.DeepEqual(array, expectedList) {
		log.Println("these should be deeply equal")
		log.Println(array)
		log.Println(expectedList)
		t.FailNow()
	}
}
