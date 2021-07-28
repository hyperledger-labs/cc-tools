package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/mock"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

func TestGetAsset(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	expectedResponse := assets.Asset{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}
	stub.MockTransactionStart("setupGetAsset")
	setupState, _ := json.Marshal(expectedResponse)
	stub.PutState("person:47061146-c642-51a1-844a-bf0b17cb5e19", setupState)
	stub.MockTransactionEnd("setupGetAsset")

	personKey := assets.Key{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}

	stub.MockTransactionStart("TestGetAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	gotAsset, err := personKey.Get(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if !reflect.DeepEqual(*gotAsset, expectedResponse) {
		log.Println("these should be deeply equal")
		log.Println(expectedResponse)
		log.Println(*gotAsset)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetAsset")
}

func TestGetCommittedAsset(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	expectedResponse := assets.Asset{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}
	stub.MockTransactionStart("setupGetAsset")
	setupState, _ := json.Marshal(expectedResponse)
	stub.PutState("person:47061146-c642-51a1-844a-bf0b17cb5e19", setupState)
	stub.MockTransactionEnd("setupGetAsset")

	personKey := assets.Key{
		"@assetType": "person",
		"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}

	stub.MockTransactionStart("TestGetAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	gotAsset, err := personKey.GetCommitted(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if !reflect.DeepEqual(*gotAsset, expectedResponse) {
		log.Println("these should be deeply equal")
		log.Println(expectedResponse)
		log.Println(*gotAsset)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetAsset")
}

func TestGetRecursive(t *testing.T) {
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

	stub.MockTransactionStart("setupReadAsset")
	stub.PutState("person:47061146-c642-51a1-844a-bf0b17cb5e19", setupPersonJSON)
	stub.PutState("book:a36a2920-c405-51c3-b584-dcd758338cb5", setupBookJSON)
	refIdx, err := stub.CreateCompositeKey("person:47061146-c642-51a1-844a-bf0b17cb5e19", []string{"book:a36a2920-c405-51c3-b584-dcd758338cb5"})
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	stub.PutState(refIdx, []byte{0x00})
	stub.MockTransactionEnd("setupReadAsset")

	bookKey := assets.Key{
		"@assetType": "book",
		"@key":       "book:a36a2920-c405-51c3-b584-dcd758338cb5",
	}
	expectedResponse := map[string]interface{}{
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org2MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "book",
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			"@lastTouchBy": "org1MSP",
			"@lastTx":      "createAsset",
			"@assetType":   "person",
			"name":         "Maria",
			"id":           "31820792048",
			"height":       0.0,
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	stub.MockTransactionStart("TestGetAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	gotAsset, err := bookKey.GetRecursive(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if !reflect.DeepEqual(gotAsset, expectedResponse) {
		log.Println("these should be deeply equal")
		log.Println(expectedResponse)
		log.Println(gotAsset)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetAsset")
}

func TestGetRecursiveWithPvtData(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	setupSecret := map[string]interface{}{
		"@assetType":   "secret",
		"@key":         "secret:73a3f9a7-eb91-5f4d-b1bb-c0487e90f40b",
		"@lastTouchBy": "org2MSP",
		"@lastTx":      "createAsset",
		"secretName":   "testSecret",
		"secret":       "this is very secret",
	}
	setupLibrary := map[string]interface{}{
		"@assetType":   "library",
		"@key":         "library:37262f3f-5f08-5649-b488-e5abaad266e1",
		"@lastTouchBy": "org3MSP",
		"@lastTx":      "createAsset",
		"name":         "Biblioteca Maria da Silva",
		"entranceCode": map[string]interface{}{
			"@assetType": "secret",
			"@key":       "secret:73a3f9a7-eb91-5f4d-b1bb-c0487e90f40b",
		},
	}
	setupSecretJSON, _ := json.Marshal(setupSecret)
	setupLibraryJSON, _ := json.Marshal(setupLibrary)

	stub.MockTransactionStart("setupReadAsset")
	stub.PutPrivateData("secret", "secret:73a3f9a7-eb91-5f4d-b1bb-c0487e90f40b", setupSecretJSON)
	stub.PutState("library:37262f3f-5f08-5649-b488-e5abaad266e1", setupLibraryJSON)
	refIdx, err := stub.CreateCompositeKey("secret:73a3f9a7-eb91-5f4d-b1bb-c0487e90f40b", []string{"library:37262f3f-5f08-5649-b488-e5abaad266e1"})
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	stub.PutState(refIdx, []byte{0x00})
	stub.MockTransactionEnd("setupReadAsset")

	libraryKey := assets.Key{
		"@assetType": "library",
		"@key":       "library:37262f3f-5f08-5649-b488-e5abaad266e1",
	}
	expectedResponse := map[string]interface{}{
		"@assetType":   "library",
		"@key":         "library:37262f3f-5f08-5649-b488-e5abaad266e1",
		"@lastTouchBy": "org3MSP",
		"@lastTx":      "createAsset",
		"name":         "Biblioteca Maria da Silva",
		"entranceCode": map[string]interface{}{
			"@assetType":   "secret",
			"@key":         "secret:73a3f9a7-eb91-5f4d-b1bb-c0487e90f40b",
			"@lastTouchBy": "org2MSP",
			"@lastTx":      "createAsset",
			"secretName":   "testSecret",
			"secret":       "this is very secret",
		},
	}

	stub.MockTransactionStart("TestGetAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	gotAsset, err := libraryKey.GetRecursive(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if !reflect.DeepEqual(gotAsset, expectedResponse) {
		log.Println("these should be deeply equal")
		log.Println(expectedResponse)
		log.Println(gotAsset)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetAsset")
}

func TestGetAssetNoKey(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	personKey := assets.Key{
		"@assetType": "person",
		// "@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	}

	stub.MockTransactionStart("TestGetAssetNoKey")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err := personKey.Get(sw)
	if err.Status() != 500 || err.Message() != "key cannot be empty" {
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetAssetNoKey")
}
