package assets

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSetProp(t *testing.T) {
	assetJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}")
	var a Asset
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if a.GetProp("readerScore") != 70.0 {
		fmt.Println("expected readerScore to be 70.0")
		t.FailNow()
	}

	err = a.SetProp("readerScore", 75)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if a.GetProp("readerScore") != 75.0 {
		fmt.Println("expected readerScore to be 75.0")
		t.FailNow()
	}
}

func TestSetPropFail(t *testing.T) {
	assetJSON := []byte("{\"@assetType\": \"samplePerson\",\"name\": \"Maria\",\"cpf\": \"318.207.920-48\",\"readerScore\": 70}")
	var a Asset
	err := json.Unmarshal(assetJSON, &a)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if a.SetProp("", 2) == nil {
		fmt.Println("expected SetProp to fail when receiving empty prop tag")
		t.FailNow()
	}

	if a.SetProp("@assetType", "sampleBook") == nil {
		fmt.Println("expected SetProp to fail when updating internal (@) props")
		t.FailNow()
	}

	if a.SetProp("favoriteColor", "green") == nil {
		fmt.Println("expected SetProp to fail when asset does not have a prop with given tag")
		t.FailNow()
	}

	ccerr := a.SetProp("name", "Bruno")
	if ccerr == nil {
		fmt.Println("expected SetProp to fail when setting key property")
		t.FailNow()
	}
	if ccerr.Status() != 501 {
		fmt.Println("expected SetProp error code 501 when setting key property")
	}

	if a.SetProp("readerScore", "seventy") == nil {
		fmt.Println("expected SetProp to fail when new prop value is invalid")
		t.FailNow()
	}
}
