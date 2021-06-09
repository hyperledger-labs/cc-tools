package assets

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestPutAsset(t *testing.T) {
	assetJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}")
	var a Asset
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub := shim.NewMockStub("testcc", new(testCC))
	stub.MockTransactionStart("TestPutAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err = a.put(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestPutAsset")
}

func TestPutAssetWithSubAsset(t *testing.T) {
	stub := shim.NewMockStub("testcc", new(testCC))

	stub.MockTransactionStart("TestPutAssetWithSubAsset")
	assetJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}")
	var a Asset
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
	stub.MockTransactionEnd("TestPutAssetWithSubAsset")

	stub.MockTransactionStart("TestPutAssetWithSubAsset")
	assetJSON = []byte("{\"@assetType\": \"author\",\"person\": {\"@assetType\": \"samplePerson\",\"name\": \"Maria\"}}")
	var b Asset
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
	stub.MockTransactionEnd("TestPutAssetWithSubAsset")

}

func TestPutAssetRecursive(t *testing.T) {
	stub := shim.NewMockStub("testcc", new(testCC))

	stub.MockTransactionStart("TestPutAssetRecursive")
	assetJSON := []byte("{\"@assetType\": \"author\",\"person\": {\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}}")
	var a Asset
	var obj map[string]interface{}
	err := json.Unmarshal(assetJSON, &obj)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	a, err = NewAsset(obj)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	_, err = PutNewRecursive(sw, obj)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestPutAssetRecursive")

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

	subAssetKey, err := NewKey(a.GetProp("person").(map[string]interface{}))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	gotSubAsset, err := subAssetKey.Get(sw)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if gotSubAsset == nil {
		fmt.Println("subasset not found")
		t.FailNow()
	}
}
