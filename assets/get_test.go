package assets

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestGetAsset(t *testing.T) {
	assetJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}")
	var a Asset
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub := shim.NewMockStub("testcc", new(testCC))
	stub.MockTransactionStart("TestGetAsset")
	_, err = a.put(stub)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	gotAsset, err := a.Get(stub)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if !reflect.DeepEqual(a, *gotAsset) {
		fmt.Println("these should be deeply equal")
		fmt.Println(a)
		fmt.Println(*gotAsset)
		t.FailNow()
	}
}

func TestGetRecursive(t *testing.T) {
	stub := shim.NewMockStub("testcc", new(testCC))

	stub.MockTransactionStart("TestGetRecursive")
	var a Asset
	assetJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}")
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	_, err = a.put(stub)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetRecursive")

	stub.MockTransactionStart("TestGetRecursive")
	var b Asset
	assetJSON = []byte("{\"@assetType\": \"sampleBook\",\"title\": \"Meu Nome é Maria\",\"author\": \"Maria Viana\",\"currentTenant\": {\"name\": \"Maria\"},\"genres\": [\"biography\", \"non-fiction\"],\"published\": \"2019-05-06T22:12:41Z\"}")
	err = json.Unmarshal(assetJSON, &b)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	_, err = b.put(stub)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetRecursive")

	gotAsset, err := b.GetRecursive(stub)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(gotAsset.GetProp("currentTenant"), a) {
		fmt.Println("these should be deeply equal")
		fmt.Println(a)
		fmt.Println(gotAsset.GetProp("currentTenant"))
		t.FailNow()
	}
}

func TestGetRecursiveWithPvtData(t *testing.T) {
	stub := shim.NewMockStub("testcc", new(testCC))

	stub.MockTransactionStart("TestGetRecursiveWithPvtData")
	var a Asset
	assetJSON := []byte("{\"@assetType\": \"sampleSecret\",\"secretName\": \"mySecret\",\"secret\": \"VERYVERYSECRET\"}")
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	_, err = a.put(stub)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetRecursiveWithPvtData")

	stub.MockTransactionStart("TestGetRecursiveWithPvtData")
	var b Asset
	assetJSON = []byte("{\"@assetType\": \"sampleBook\",\"title\": \"Meu Nome é Maria\",\"author\": \"Maria Viana\",\"genres\": [\"biography\", \"non-fiction\"],\"published\": \"2019-05-06T22:12:41Z\",\"secret\":{\"@assetType\": \"sampleSecret\",\"secretName\": \"mySecret\",\"secret\": \"VERYVERYSECRET\"}}")
	err = json.Unmarshal(assetJSON, &b)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	_, err = b.put(stub)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetRecursiveWithPvtData")

	gotAsset, err := b.GetRecursive(stub)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(gotAsset.GetProp("secret"), a) {
		fmt.Println("these should be deeply equal")
		fmt.Println(a)
		fmt.Println(gotAsset.GetProp("currentTenant"))
		t.FailNow()
	}
}
