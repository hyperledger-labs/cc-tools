package assets

import (
	"encoding/json"
	"fmt"
	"testing"

	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestDeleteAsset(t *testing.T) {
	stub := shim.NewMockStub("testcc", new(testCC))

	stub.MockTransactionStart("TestDeleteAsset")
	var a Asset
	assetJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}")
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err = a.put(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestDeleteAsset")

	stub.MockTransactionStart("TestDeleteAsset")

	a.delete(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	stub.MockTransactionEnd("TestDeleteAsset")

	stub.MockTransactionStart("TestDeleteAsset")

	res, err := a.Get(sw)
	if err == nil {
		fmt.Println("should not be capable of getting asset", res)
		t.FailNow()
	}

	stub.MockTransactionEnd("TestDeleteAsset")
}

func TestDeleteAssetCascade(t *testing.T) {
	stub := shim.NewMockStub("testcc", new(testCC))

	stub.MockTransactionStart("TestDeleteAssetCascade")
	var personAsset Asset
	personJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}")
	err := json.Unmarshal(personJSON, &personAsset)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err = personAsset.put(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestDeleteAssetCascade")

	stub.MockTransactionStart("TestDeleteAssetCascade")
	tenant, err := NewKey(personAsset)
	bookString := fmt.Sprintf("{\"@assetType\": \"sampleBook\",\"title\": \"Pale Blue Dot\",\"author\": \"Carl Sagan\",\"currentTenant\": %s}", tenant)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	bookJSON := []byte(bookString)
	var bookAsset Asset
	err = json.Unmarshal(bookJSON, &bookAsset)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	_, err = bookAsset.put(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestDeleteAssetCascade")

	stub.MockTransactionStart("TestDeleteAssetCascade")
	var k Key
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	keyJSON := []byte(bookString)
	err = json.Unmarshal(keyJSON, &k)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	deletedKeys := make([]string, 0)
	err = deleteRecursive(sw, personAsset.Key(), &deletedKeys)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestDeleteAssetCascade")

	stub.MockTransactionStart("TestGetAsset")

	res, err := bookAsset.Get(sw)
	if err == nil {
		fmt.Println("should not be capable of getting asset", res)
		t.FailNow()
	}

	res, err = personAsset.Get(sw)
	if err == nil {
		fmt.Println("should not be capable of getting asset", res)
		t.FailNow()
	}

	stub.MockTransactionEnd("TestGetAsset")
}
