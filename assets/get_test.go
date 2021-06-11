package assets

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	sw "github.com/goledgerdev/cc-tools/stubwrapper"
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
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err = a.put(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	gotAsset, err := a.Get(sw)
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
	stub.MockTransactionEnd("TestGetAsset")
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
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err = a.put(sw)
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
	_, err = b.put(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetRecursive")

	gotAsset, err := b.GetRecursive(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(gotAsset["currentTenant"], (map[string]interface{})(a)) {
		fmt.Println("these should be deeply equal")
		fmt.Println((map[string]interface{})(a))
		fmt.Println(gotAsset["currentTenant"])
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
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err = a.put(sw)
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
	_, err = b.put(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetRecursiveWithPvtData")

	gotAsset, err := b.GetRecursive(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(gotAsset["secret"], (map[string]interface{})(a)) {
		fmt.Println("these should be deeply equal")
		fmt.Println((map[string]interface{})(a))
		fmt.Println(gotAsset["secret"])
		t.FailNow()
	}
}

func TestGetRecursiveWithListOfPvtData(t *testing.T) {
	stub := shim.NewMockStub("testcc", new(testCC))

	stub.MockTransactionStart("TestGetRecursiveWithListOfPvtData")
	var a Asset
	assetJSON := []byte("{\"@assetType\": \"sampleSecret\",\"secretName\": \"mySecret\",\"secret\": \"VERYVERYSECRET\"}")
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
	stub.MockTransactionEnd("TestGetRecursiveWithListOfPvtData")

	stub.MockTransactionStart("TestGetRecursiveWithListOfPvtData")
	var b Asset
	assetJSON = []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70,\"secrets\":[{\"@assetType\": \"sampleSecret\",\"secretName\": \"mySecret\",\"secret\": \"VERYVERYSECRET\"}]}")
	err = json.Unmarshal(assetJSON, &b)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	_, err = b.put(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestGetRecursiveWithListOfPvtData")

	gotAsset, err := b.GetRecursive(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	secretProp := gotAsset["secrets"]
	secretList, ok := secretProp.([]interface{})
	if !ok {
		fmt.Println("secretList should be of type []interface{}")
		t.FailNow()
	}

	if len(secretList) != 1 {
		fmt.Println("secretList should have length equal to 1")
	}

	if !reflect.DeepEqual(secretList[0], (map[string]interface{})(a)) {
		fmt.Println("these should be deeply equal")
		fmt.Println((map[string]interface{})(a))
		fmt.Println(secretList[0])
		t.FailNow()
	}
}
