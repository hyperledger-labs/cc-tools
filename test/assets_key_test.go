package test

import (
	"log"
	"reflect"
	"testing"

	"github.com/hyperledger-labs/cc-tools/assets"
)

func TestKeyString(t *testing.T) {
	a := assets.Key{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}

	expectedString := `{"@assetType":"person","@key":"person:47061146-c642-51a1-844a-bf0b17cb5e19"}`

	if !reflect.DeepEqual(a.String(), expectedString) {
		log.Println("these should be deeply equal")
		log.Println(a.String())
		log.Println(expectedString)
		t.FailNow()
	}
}

func TestKeyJSON(t *testing.T) {
	a := assets.Key{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}

	expectedJSON := []byte(`{"@assetType":"person","@key":"person:47061146-c642-51a1-844a-bf0b17cb5e19"}`)

	if !reflect.DeepEqual(a.JSON(), expectedJSON) {
		log.Println("these should be deeply equal")
		log.Println(a.String())
		log.Println(expectedJSON)
		t.FailNow()
	}
}

func TestKeyGenUsingAssetType(t *testing.T) {
	m := map[string]interface{}{
		"@assetType": "person",
		"id":         "31820792048",
	}

	key, err := assets.NewKey(m)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	expectedKey := assets.Key{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}

	if !reflect.DeepEqual(key, expectedKey) {
		log.Println("these should be deeply equal")
		log.Println(key)
		log.Println(expectedKey)
		t.FailNow()
	}
}

func TestKeyGenUsingKey(t *testing.T) {
	m := map[string]interface{}{
		"@key": "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}

	key, err := assets.NewKey(m)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	expectedKey := assets.Key{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}

	if !reflect.DeepEqual(key, expectedKey) {
		log.Println("these should be deeply equal")
		log.Println(key)
		log.Println(expectedKey)
		t.FailNow()
	}
}
