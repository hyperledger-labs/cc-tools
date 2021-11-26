package test

import (
	"encoding/json"
	"testing"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestSearchEmptyQuery(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	// expectedResponse := map[string]interface{}{
	// 	"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
	// 	"@lastTouchBy": "org1MSP",
	// 	"@lastTx":      "createAsset",
	// 	"@assetType":   "person",
	// 	"name":         "Maria",
	// 	"id":           "31820792048",
	// 	"height":       0.0,
	// }

	req := map[string]interface{}{
		// "query": map[string]interface{}{},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.FailNow()
	}
	res := stub.MockInvoke("search", [][]byte{
		[]byte("search"),
		reqBytes,
	})

	if res.GetStatus() != 400 {
		t.FailNow()
	}

	if res.GetMessage() != "unable to get args: missing argument 'query'" {
		t.FailNow()
	}
}
