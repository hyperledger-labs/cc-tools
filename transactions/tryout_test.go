package transactions

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestTryout(t *testing.T) {
	var err error
	stub := shim.NewMockStub("org1MSP", new(testCC))

	// Create Asset
	reqPerson := map[string]interface{}{
		"asset": []map[string]interface{}{
			{
				"@assetType": "person",
				"name":       "Maria",
				"id":         "318.207.920-48",
			},
		},
	}
	expectedPerson := []interface{}{
		map[string]interface{}{
			"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			"@lastTouchBy": "org1MSP",
			"@lastTx":      "createAsset",
			"@assetType":   "person",
			"name":         "Maria",
			"id":           "31820792048",
			"height":       0.0,
		},
	}

	err = invokeAndVerify(stub, "createAsset", reqPerson, expectedPerson, 200)
	if err != nil {
		log.Println("createAsset fail")
		t.FailNow()
	}

	// Create book
	stub.Name = "org2MSP"

	reqBook := map[string]interface{}{
		"asset": []map[string]interface{}{
			{
				"@assetType": "book",
				"title":      "Meu Nome é Maria",
				"author":     "Maria Viana",
				"currentTenant": map[string]interface{}{
					"id": "318.207.920-48",
				},
				"genres":    []string{"biography", "non-fiction"},
				"published": "2019-05-06T22:12:41Z",
			},
		},
	}

	expectedBook := []interface{}{
		map[string]interface{}{
			"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
			"@lastTouchBy": "org2MSP",
			"@lastTx":      "createAsset",
			"@assetType":   "book",
			"title":        "Meu Nome é Maria",
			"author":       "Maria Viana",
			"currentTenant": map[string]interface{}{
				"@assetType": "person",
				"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			},
			"genres":    []interface{}{"biography", "non-fiction"},
			"published": "2019-05-06T22:12:41Z",
		},
	}

	err = invokeAndVerify(stub, "createAsset", reqBook, expectedBook, 200)
	if err != nil {
		log.Println("createAsset fail")
		t.FailNow()
	}
}

func invokeAndVerify(stub *shim.MockStub, txName string, req, expectedRes interface{}, expectedStatus int32) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return err
	}

	res := stub.MockInvoke(txName, [][]byte{
		[]byte(txName),
		reqBytes,
	})

	if res.GetStatus() != expectedStatus {
		log.Println(res.Message)
		return fmt.Errorf("expected %d got %d", expectedStatus, res.GetStatus())
	}

	var resPayload interface{}
	err = json.Unmarshal(res.GetPayload(), &resPayload)
	if err != nil {
		log.Println(res.GetPayload())
		log.Println(err)
		return err
	}

	if !reflect.DeepEqual(resPayload, expectedRes) {
		log.Println("these should be equal")
		log.Printf("%#v\n", resPayload)
		log.Printf("%#v\n", expectedRes)
		return fmt.Errorf("unexpected response")
	}

	return nil
}
