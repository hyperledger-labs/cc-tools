package test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/mock"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

func TestDeleteWithRef(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	setupPerson := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}
	setupBook := map[string]interface{}{
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
	setupPersonJSON, _ := json.Marshal(setupPerson)
	setupBookJSON, _ := json.Marshal(setupBook)

	stub.MockTransactionStart("setupDeleteWithRef")
	stub.PutState("person:47061146-c642-51a1-844a-bf0b17cb5e19", setupPersonJSON)
	stub.PutState("book:a36a2920-c405-51c3-b584-dcd758338cb5", setupBookJSON)
	refIdx, err := stub.CreateCompositeKey("person:47061146-c642-51a1-844a-bf0b17cb5e19", []string{"book:a36a2920-c405-51c3-b584-dcd758338cb5"})
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	stub.PutState(refIdx, []byte{0x00})
	stub.MockTransactionEnd("setupDeleteWithRef")

	personAsset := assets.Asset{
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@assetType": "person",
		"id":         "31820792048",
	}

	stub.MockTransactionStart("TestDeleteWithRef")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, ccerr := personAsset.Delete(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if ccerr.Status() != 400 {
		log.Printf("expected err status: %d\n", 400)
		log.Printf("     got err status: %d\n", ccerr.Status())
		t.FailNow()
	}
	if ccerr.Message() != "another asset holds a reference to this one" {
		log.Printf("expected err msg: %s\n", "another asset holds a reference to this one")
		log.Printf("     got err msg: %s\n", ccerr.Message())
		t.FailNow()
	}

	stub.MockTransactionEnd("TestDeleteWithRef")
}
