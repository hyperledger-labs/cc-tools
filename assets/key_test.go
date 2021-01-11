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
	if k.String() == "" {
		fmt.Println("failed to print asset as string")
	}
}

func TestKeyIsPrivate(t *testing.T) {
	keyJSON := []byte("{\"@assetType\": \"sampleBook\",\"title\": \"Meu Nome é Maria\",\"author\": \"Maria Viana\"}")
	var k Key
	err := json.Unmarshal(keyJSON, &k)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if k.IsPrivate() {
		fmt.Println("false positive in Key.IsPrivate")
		t.FailNow()
	}

	pvtKeyJSON := []byte("{\"@assetType\": \"sampleSecret\",\"secretName\": \"mySecret\"}")
	var pvtKey Key
	err = json.Unmarshal(pvtKeyJSON, &pvtKey)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !pvtKey.IsPrivate() {
		fmt.Println("false negative in Key.IsPrivate")
		t.FailNow()
	}
}

func TestNewKeyWithNilMap(t *testing.T) {
	_, err := NewKey(nil)
	if err == nil {
		fmt.Println("NewKey should fail if nil map is passed as argument")
		t.FailNow()
	}
}
