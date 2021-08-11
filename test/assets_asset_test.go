package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
)

func TestAssetUnmarshal(t *testing.T) {
	assetJSON := []byte("{\"@assetType\": \"person\",\"name\": \"Maria\",\"id\": \"318.207.920-48\",\"height\": 0}")
	expectedAsset := assets.Asset{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"id":         "31820792048",
		"name":       "Maria",
		"height":     0.0,
	}
	var a assets.Asset
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if !reflect.DeepEqual(a, expectedAsset) {
		log.Println("these should be deeply equal")
		log.Printf("%#v\n", expectedAsset)
		log.Printf("%#v\n", a)
	}
}

func TestAssetIsPrivate(t *testing.T) {
	a := assets.Asset{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"name":       "Maria",
		"id":         "31820792048",
		"height":     0.0,
	}

	if a.IsPrivate() {
		log.Println("false positive in Asset.IsPrivate")
		t.FailNow()
	}

	s := assets.Asset{
		"@assetType": "secret",
		"@key":       "secret:73a3f9a7-eb91-5f4d-b1bb-c0487e90f40b",
		"secretName": "testSecret",
		"secret":     "this is very secret",
	}

	if !s.IsPrivate() {
		log.Println("false negative in Asset.IsPrivate")
		t.FailNow()
	}
}

func TestAssetString(t *testing.T) {
	a := assets.Asset{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"name":       "Maria",
		"id":         "31820792048",
		"height":     0.0,
	}

	expectedString := `{"@assetType":"person","@key":"person:47061146-c642-51a1-844a-bf0b17cb5e19","height":0,"id":"31820792048","name":"Maria"}`

	if !reflect.DeepEqual(a.String(), expectedString) {
		log.Println("these should be deeply equal")
		log.Println(a.String())
		log.Println(expectedString)
		t.FailNow()
	}
}

func TestAssetJSON(t *testing.T) {
	a := assets.Asset{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"name":       "Maria",
		"id":         "31820792048",
		"height":     0.0,
	}

	expectedJSON := []byte(`{"@assetType":"person","@key":"person:47061146-c642-51a1-844a-bf0b17cb5e19","height":0,"id":"31820792048","name":"Maria"}`)

	if !reflect.DeepEqual(a.JSON(), expectedJSON) {
		log.Println("these should be deeply equal")
		log.Println(a.String())
		log.Println(expectedJSON)
		t.FailNow()
	}
}
