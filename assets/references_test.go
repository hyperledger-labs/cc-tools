package assets

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAssetRefs(t *testing.T) {
	assetJSON := []byte("{\"@assetType\": \"sampleBook\",\"title\": \"Meu Nome é Maria\",\"author\": \"Maria Viana\",\"currentTenant\": {\"name\": \"Maria\"},\"genres\": [\"biography\", \"non-fiction\"],\"published\": \"2019-05-06T22:12:41Z\"}")
	var a Asset
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	refs, err := a.Refs()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if len(refs) == 0 {
		fmt.Println("no refs return by asset.Refs")
		t.FailNow()
	}
	if len(refs) > 1 {
		fmt.Println("too many refs return by asset.Refs")
		t.FailNow()
	}
	refKey := refs[0]
	if refKey.TypeTag() != "samplePerson" {
		fmt.Println("reference returned by asset.Refs should be of type samplePerson")
		t.FailNow()
	}
}
