package test

import (
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
)

func TestBuildAssetPropValid(t *testing.T) {
	propMap := map[string]interface{}{
		"tag":          "id",
		"label":        "CPF (Brazilian ID)",
		"description":  "",
		"isKey":        true,
		"required":     true,
		"readOnly":     false,
		"defaultValue": nil,
		"dataType":     "cpf",
		"writers":      []interface{}{"org1MSP"},
	}
	prop, err := assets.BuildAssetProp(propMap, nil)

	expectedProp := assets.AssetProp{
		Tag:          "id",
		Label:        "CPF (Brazilian ID)",
		Description:  "",
		IsKey:        true,
		Required:     true,
		ReadOnly:     false,
		DefaultValue: nil,
		DataType:     "cpf",
		Writers:      []string{"org1MSP"},
	}

	if err != nil {
		log.Println("an error should not have occurred")
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(prop, expectedProp) {
		log.Println("these should be deeply equal")
		log.Println(prop)
		log.Println(expectedProp)
		t.FailNow()
	}
}

func TestBuildAssetPropInvalid(t *testing.T) {
	propMap := map[string]interface{}{
		"tag":          "id",
		"label":        "CPF (Brazilian ID)",
		"description":  "",
		"isKey":        true,
		"required":     true,
		"readOnly":     false,
		"defaultValue": nil,
		"dataType":     "inexistant",
		"writers":      []interface{}{"org1MSP"},
	}
	_, err := assets.BuildAssetProp(propMap, nil)

	err.Status()
	if err.Status() != 400 {
		log.Println(err)
		t.FailNow()
	}

	if err.Message() != "failed checking data type: invalid dataType value 'inexistant'" {
		log.Printf("error message different from expected: %s", err.Message())
		t.FailNow()
	}
}

func TestHandlePropUpdate(t *testing.T) {
	prop := assets.AssetProp{
		Tag:          "id",
		Label:        "CPF (Brazilian ID)",
		Description:  "",
		IsKey:        true,
		Required:     true,
		ReadOnly:     false,
		DefaultValue: nil,
		DataType:     "cpf",
	}

	propUpdateMap := map[string]interface{}{
		"writers":      []interface{}{"org1MSP"},
		"defaultValue": "12345678901",
	}

	updatedProp, err := assets.HandlePropUpdate(prop, propUpdateMap)

	if err != nil {
		log.Println("an error should not have occurred")
		log.Println(err)
		t.FailNow()
	}

	prop.DefaultValue = "12345678901"
	prop.Writers = []string{"org1MSP"}

	if !reflect.DeepEqual(updatedProp, prop) {
		log.Println("these should be deeply equal")
		log.Println(updatedProp)
		log.Println(prop)
		t.FailNow()
	}
}

func TestBuildAssetPropWithReferenceList(t *testing.T) {
	newTypeList := []interface{}{
		map[string]interface{}{
			"tag":   "newType",
			"label": "New Type",
		},
	}

	propMap := map[string]interface{}{
		"tag":          "id",
		"label":        "CPF (Brazilian ID)",
		"description":  "",
		"isKey":        true,
		"required":     true,
		"readOnly":     false,
		"defaultValue": nil,
		"dataType":     "->newType",
		"writers":      []interface{}{"org1MSP"},
	}
	prop, err := assets.BuildAssetProp(propMap, newTypeList)

	expectedProp := assets.AssetProp{
		Tag:          "id",
		Label:        "CPF (Brazilian ID)",
		Description:  "",
		IsKey:        true,
		Required:     true,
		ReadOnly:     false,
		DefaultValue: nil,
		DataType:     "->newType",
		Writers:      []string{"org1MSP"},
	}

	if err != nil {
		log.Println("an error should not have occurred")
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(prop, expectedProp) {
		log.Println("these should be deeply equal")
		log.Println(prop)
		log.Println(expectedProp)
		t.FailNow()
	}
}
