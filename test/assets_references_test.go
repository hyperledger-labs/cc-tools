package test

import (
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/mock"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

func TestReferrers(t *testing.T) {
	fabricStub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	setupPerson := assets.Asset{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}
	setupBook := assets.Asset{
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org2MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "book",
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	stub := &sw.StubWrapper{
		Stub: fabricStub,
	}
	fabricStub.MockTransactionStart("setupReadAsset")
	_, err := setupPerson.PutNew(stub)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	_, err = setupBook.PutNew(stub)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	fabricStub.MockTransactionEnd("setupReadAsset")

	personKey := assets.Key{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}
	expectedReferrers := []assets.Key{{
		"@assetType": "book",
		"@key":       "book:a36a2920-c405-51c3-b584-dcd758338cb5",
	}}

	fabricStub.MockTransactionStart("TestReferrersWithSameStub")
	referrers, err := setupPerson.Referrers(stub)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	fabricStub.MockTransactionEnd("TestReferrersWithSameStub")

	if !reflect.DeepEqual(referrers, expectedReferrers) {
		log.Println("these should be deeply equal")
		log.Println(expectedReferrers)
		log.Println(referrers)
		t.FailNow()
	}

	stub = &sw.StubWrapper{
		Stub: fabricStub,
	}
	fabricStub.MockTransactionStart("TestReferrers")
	referrers, err = personKey.Referrers(stub)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	fabricStub.MockTransactionEnd("TestReferrers")

	if !reflect.DeepEqual(referrers, expectedReferrers) {
		log.Println("these should be deeply equal")
		log.Println(expectedReferrers)
		log.Println(referrers)
		t.FailNow()
	}
}
