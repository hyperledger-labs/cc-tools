package test

import (
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
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
