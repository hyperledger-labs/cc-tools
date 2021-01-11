package assets

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAssetUnmarshal(t *testing.T) {
	assetJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}")
	var a Asset
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if a.Key() == "" {
		fmt.Println("asset does not have @key property after unmarshal")
		t.FailNow()
	}
	if a.String() == "" {
		fmt.Println("failed to print asset as string")
	}
}

func TestAssetIsPrivate(t *testing.T) {
	assetJSON := []byte("{\"@assetType\": \"sampleBook\",\"title\": \"Meu Nome é Maria\",\"author\": \"Maria Viana\",\"currentTenant\": {\"name\": \"Maria\"},\"genres\": [\"biography\", \"non-fiction\"],\"published\": \"2019-05-06T22:12:41Z\"}")
	var a Asset
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if a.IsPrivate() {
		fmt.Println("false positive in Asset.IsPrivate")
		t.FailNow()
	}

	pvtAssetJSON := []byte("{\"@assetType\": \"sampleSecret\",\"secretName\": \"mySecret\",\"secret\": \"VERYVERYSECRET\"}")
	var pvtAsset Asset
	err = json.Unmarshal(pvtAssetJSON, &pvtAsset)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if !pvtAsset.IsPrivate() {
		fmt.Println("false negative in Asset.IsPrivate")
		t.FailNow()
	}
}
