package assets

import (
	"encoding/json"
	"fmt"
	"testing"

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
	_, err = a.put(stub)
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
	_, err = a.put(stub)
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
	_, err = b.put(stub)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestPutAssetWithSubAsset")

}
