package assets

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCheckWriters(t *testing.T) {
	assetJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"readerScore\": 70}")
	var a Asset
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	err = a.checkWriters("org1MSP")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	a.SetProp("cpf", "318.207.920-48")

	err = a.checkWriters("org2MSP")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	err = a.checkWriters("orgInvalidMSP")
	if err == nil {
		fmt.Println("expected asset.checkWriters to fail")
		t.FailNow()
	}
}
