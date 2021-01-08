package assets

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestKeyUnmarshal(t *testing.T) {
	keyJSON := []byte("{\"@assetType\": \"sampleBook\",\"title\": \"Meu Nome é Maria\",\"author\": \"Maria Viana\"}")
	var k Key
	err := json.Unmarshal(keyJSON, &k)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if k.Key() == "" {
		fmt.Println("key does not have @key property after unmarshal")
		t.FailNow()
	}
	if k.TypeTag() == "" {
		fmt.Println("key does not have @assetType property after unmarshal")
		t.FailNow()
	}
	if k.Type() == nil {
		fmt.Println("unable to fetch asset type data from key")
		t.FailNow()
	}
	for p := range k {
		if p != "@key" && p != "@assetType" {
			fmt.Println("unmarshaling key object should erase everything but @assetType and @key")
			t.FailNow()
		}
	}
}
