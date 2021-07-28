package test

import (
	"log"
	"testing"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestTryout(t *testing.T) {
	var err error
	stub := mock.NewMockStub("org1MSP", new(testCC))

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
		log.Println("create person fail")
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
		log.Println("create book fail")
		t.FailNow()
	}

	// Read book
	stub.Name = "org1MSP"

	reqReadBook := map[string]interface{}{
		"key": map[string]interface{}{
			"@assetType": "book",
			"author":     "Maria Viana",
			"title":      "Meu Nome é Maria",
		},
	}

	expectedReadBook := expectedBook[0]

	err = invokeAndVerify(stub, "readAsset", reqReadBook, expectedReadBook, 200)
	if err != nil {
		log.Println("readAsset fail")
		t.FailNow()
	}

	// Update person
	stub.Name = "org2MSP"

	reqUpdatePerson := map[string]interface{}{
		"update": map[string]interface{}{
			"@assetType": "person",
			"id":         "318.207.920-48",
			"name":       "Maria",
			"height":     1.67,
		},
	}

	expectedUpdatePerson := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org2MSP",
		"@lastTx":      "updateAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       1.67,
	}

	err = invokeAndVerify(stub, "updateAsset", reqUpdatePerson, expectedUpdatePerson, 200)
	if err != nil {
		log.Println("updateAsset fail")
		t.FailNow()
	}

	// Read person to check if it was updated
	stub.Name = "org1MSP"

	reqReadPerson := map[string]interface{}{
		"key": map[string]interface{}{
			"@assetType": "person",
			"id":         "318.207.920-48",
		},
	}

	expectedReadPerson := expectedUpdatePerson

	err = invokeAndVerify(stub, "readAsset", reqReadPerson, expectedReadPerson, 200)
	if err != nil {
		log.Println("readAsset fail")
		t.FailNow()
	}

	// Delete book
	stub.Name = "org2MSP"

	reqDeleteBook := reqReadBook

	expectedDeleteBook := expectedBook[0]

	err = invokeAndVerify(stub, "deleteAsset", reqDeleteBook, expectedDeleteBook, 200)
	if err != nil {
		log.Println("deleteAsset fail")
		t.FailNow()
	}

	// Delete person
	stub.Name = "org1MSP"

	reqDeletePerson := reqReadPerson

	expectedDeletePerson := expectedUpdatePerson

	err = invokeAndVerify(stub, "deleteAsset", reqDeletePerson, expectedDeletePerson, 200)
	if err != nil {
		log.Println("deleteAsset fail")
		t.FailNow()
	}
}
