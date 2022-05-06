package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/mock"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

func TestPutAsset(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	person := assets.Asset{
		"@assetType": "person",
		"name":       "Maria",
		"id":         "31820792048",
	}
	stub.MockTransactionStart("TestPutAsset")
	lastUpdated, _ := stub.GetTxTimestamp()
	expectedState := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
	}
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	var err error
	_, err = person.Put(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestPutAsset")

	stateJSON := stub.State["person:47061146-c642-51a1-844a-bf0b17cb5e19"]
	var state map[string]interface{}
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expectedState, state) {
		log.Println("these should be deeply equal")
		log.Println(expectedState)
		log.Println(state)
		t.FailNow()
	}
}

func TestPutAssetWithSubAsset(t *testing.T) {
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

	setupPersonJSON, _ := json.Marshal(setupPerson)

	stub.MockTransactionStart("setupReadAsset")
	stub.State["person:47061146-c642-51a1-844a-bf0b17cb5e19"] = setupPersonJSON
	stub.MockTransactionEnd("setupReadAsset")

	stub.MockTransactionStart("TestPutAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	book := assets.Asset{
		"@assetType": "book",
		"@key":       "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"title":      "Meu Nome é Maria",
		"author":     "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"id":         "31820792048",
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	var err error
	putBook, err := book.PutNew(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	expectedState := (map[string]interface{})(book)
	if !reflect.DeepEqual(expectedState, putBook) {
		log.Println("these should be deeply equal")
		log.Println(expectedState)
		log.Println(putBook)
		t.FailNow()
	}

	stateJSON := stub.State["book:a36a2920-c405-51c3-b584-dcd758338cb5"]
	var state map[string]interface{}
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expectedState, state) {
		log.Println("these should be deeply equal")
		log.Println(expectedState)
		log.Println(state)
		t.FailNow()
	}
}

func TestPutNewAssetRecursive(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	stub.MockTransactionStart("TestPutAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	book := map[string]interface{}{
		"@assetType": "book",
		"title":      "Meu Nome é Maria",
		"author":     "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"name":       "Maria",
			"id":         "31820792048",
			"height":     1.66,
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	var err error
	putBook, err := assets.PutNewRecursive(sw, book)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	publishedTime, _ := time.Parse(time.RFC3339, "2019-05-06T22:12:41Z")
	lastUpdated, _ := stub.GetTxTimestamp()
	expectedBook := map[string]interface{}{
		"@assetType":   "book",
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType":   "person",
			"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			"@lastTouchBy": "org1MSP",
			"@lastTx":      "",
			"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
			"name":         "Maria",
			"id":           "31820792048",
			"height":       1.66,
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": publishedTime,
	}

	if !reflect.DeepEqual(expectedBook, putBook) {
		log.Println("these should be deeply equal")
		log.Println(expectedBook)
		log.Println(putBook)
		t.FailNow()
	}

	expectedState := map[string]interface{}{
		"@assetType":   "book",
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	stateJSON := stub.State["book:a36a2920-c405-51c3-b584-dcd758338cb5"]
	var state map[string]interface{}
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expectedState, state) {
		log.Println("these should be deeply equal")
		log.Println(expectedState)
		log.Println(state)
		t.FailNow()
	}
}

func TestUpdateRecursive(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// Create a book
	stub.MockTransactionStart("TestPutAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	book := map[string]interface{}{
		"@assetType": "book",
		"title":      "Meu Nome é Maria",
		"author":     "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"name":       "Maria",
			"id":         "31820792048",
			"height":     1.66,
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	var err error
	_, err = assets.PutNewRecursive(sw, book)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestPutAsset")

	// Update the book and the tenant
	stub.MockTransactionStart("TestUpdateAsset")
	book["published"] = "2022-05-06T22:12:41Z"
	book["currentTenant"] = map[string]interface{}{
		"@assetType": "person",
		"name":       "Maria",
		"id":         "31820792048",
		"height":     1.88,
	}

	putBook, err := assets.UpdateRecursive(sw, book)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	publishedTime, _ := time.Parse(time.RFC3339, "2022-05-06T22:12:41Z")
	lastUpdated, _ := stub.GetTxTimestamp()
	expectedBook := map[string]interface{}{
		"@assetType":   "book",
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType":   "person",
			"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			"@lastTouchBy": "org1MSP",
			"@lastTx":      "",
			"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
			"name":         "Maria",
			"id":           "31820792048",
			"height":       1.88,
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": publishedTime,
	}

	// Check if the book is updated
	if !reflect.DeepEqual(expectedBook, putBook) {
		log.Println("these should be deeply equal")
		log.Println(expectedBook)
		log.Println(putBook)
		t.FailNow()
	}

	expectedState := map[string]interface{}{
		"@assetType":   "book",
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2022-05-06T22:12:41Z",
	}

	stateJSON := stub.State["book:a36a2920-c405-51c3-b584-dcd758338cb5"]
	var state map[string]interface{}
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expectedState, state) {
		log.Println("these should be deeply equal")
		log.Println(expectedState)
		log.Println(state)
		t.FailNow()
	}
}
