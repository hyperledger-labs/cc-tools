package test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/mock"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

func TestPutAsset(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))
	person := assets.Asset{
		"@assetType": "person",
		"name":       "Maria",
		"id":         "31820792048",
	}
	stub.MockTransactionStart("TestPutAsset")
	lastUpdated, _ := stub.GetTxTimestamp()
	expectedState := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
	}
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	var err error
	_, err = person.Put(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	stub.MockTransactionEnd("TestPutAsset")

	stateJSON := stub.State["person:47061146-c642-51a1-844a-bf0b17cb5e19"]
	var state map[string]interface{}
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expectedState, state) {
		log.Println("these should be deeply equal")
		log.Println(expectedState)
		log.Println(state)
		t.FailNow()
	}
}

func TestPutAssetWithSubAsset(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	// State setup
	setupPerson := map[string]interface{}{
		"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "createAsset",
		"@assetType":   "person",
		"name":         "Maria",
		"id":           "31820792048",
		"height":       0.0,
	}

	setupPersonJSON, _ := json.Marshal(setupPerson)

	stub.MockTransactionStart("setupReadAsset")
	stub.State["person:47061146-c642-51a1-844a-bf0b17cb5e19"] = setupPersonJSON
	stub.MockTransactionEnd("setupReadAsset")

	stub.MockTransactionStart("TestPutAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	book := assets.Asset{
		"@assetType": "book",
		"@key":       "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"title":      "Meu Nome é Maria",
		"author":     "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"id":         "31820792048",
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	var err error
	putBook, err := book.PutNew(sw)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	expectedState := (map[string]interface{})(book)
	if !reflect.DeepEqual(expectedState, putBook) {
		log.Println("these should be deeply equal")
		log.Println(expectedState)
		log.Println(putBook)
		t.FailNow()
	}

	stateJSON := stub.State["book:a36a2920-c405-51c3-b584-dcd758338cb5"]
	var state map[string]interface{}
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expectedState, state) {
		log.Println("these should be deeply equal")
		log.Println(expectedState)
		log.Println(state)
		t.FailNow()
	}
}

func TestPutNewAssetRecursive(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	stub.MockTransactionStart("TestPutAsset")
	sw := &sw.StubWrapper{
		Stub: stub,
	}
	book := map[string]interface{}{
		"@assetType": "book",
		"title":      "Meu Nome é Maria",
		"author":     "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"name":       "Maria",
			"id":         "31820792048",
			"height":     1.66,
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	var err error
	putBook, err := assets.PutNewRecursive(sw, book)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	publishedTime, _ := time.Parse(time.RFC3339, "2019-05-06T22:12:41Z")
	lastUpdated, _ := stub.GetTxTimestamp()
	expectedBook := map[string]interface{}{
		"@assetType":   "book",
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType":   "person",
			"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
			"@lastTouchBy": "org1MSP",
			"@lastTx":      "",
			"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
			"name":         "Maria",
			"id":           "31820792048",
			"height":       1.66,
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": publishedTime,
	}

	if !reflect.DeepEqual(expectedBook, putBook) {
		log.Println("these should be deeply equal")
		log.Println(expectedBook)
		log.Println(putBook)
		t.FailNow()
	}

	expectedState := map[string]interface{}{
		"@assetType":   "book",
		"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
		"@lastTouchBy": "org1MSP",
		"@lastTx":      "",
		"@lastUpdated": lastUpdated.AsTime().Format(time.RFC3339),
		"title":        "Meu Nome é Maria",
		"author":       "Maria Viana",
		"currentTenant": map[string]interface{}{
			"@assetType": "person",
			"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
		},
		"genres":    []interface{}{"biography", "non-fiction"},
		"published": "2019-05-06T22:12:41Z",
	}

	stateJSON := stub.State["book:a36a2920-c405-51c3-b584-dcd758338cb5"]
	var state map[string]interface{}
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expectedState, state) {
		log.Println("these should be deeply equal")
		log.Println(expectedState)
		log.Println(state)
		t.FailNow()
	}
}

func TestUpdateRecursive(t *testing.T) {
	tests := []struct {
		description string

		assetkey         string
		asset            map[string]interface{}
		updateReq        map[string]interface{}
		expectedResponse func(lastUpdated string) map[string]interface{}
		expectedState    func(lastUpdated string) map[string]interface{}
	}{
		{
			description: "update recursive book",

			assetkey: "book:a36a2920-c405-51c3-b584-dcd758338cb5",
			asset: map[string]interface{}{
				"@assetType": "book",
				"title":      "Meu Nome é Maria",
				"author":     "Maria Viana",
				"currentTenant": map[string]interface{}{
					"@assetType": "person",
					"name":       "Maria",
					"id":         "31820792048",
					"height":     1.66,
				},
				"genres":    []interface{}{"biography", "non-fiction"},
				"published": "2019-05-06T22:12:41Z",
			},
			updateReq: map[string]interface{}{
				"@assetType": "book",
				"author":     "Maria Viana",
				"title":      "Meu Nome é Maria",
				"published":  "2022-05-06T22:12:41Z",
				"currentTenant": map[string]interface{}{
					"@assetType": "person",
					"id":         "31820792048",
					"height":     1.88,
				},
			},
			expectedResponse: func(lastUpdated string) map[string]interface{} {
				publishedTime, _ := time.Parse(time.RFC3339, "2022-05-06T22:12:41Z")
				return map[string]interface{}{
					"@assetType":   "book",
					"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
					"@lastTouchBy": "org1MSP",
					"@lastTx":      "",
					"@lastUpdated": lastUpdated,
					"title":        "Meu Nome é Maria",
					"author":       "Maria Viana",
					"currentTenant": map[string]interface{}{
						"@assetType":   "person",
						"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
						"@lastTouchBy": "org1MSP",
						"@lastTx":      "",
						"@lastUpdated": lastUpdated,
						"name":         "Maria",
						"id":           "31820792048",
						"height":       1.88,
					},
					"genres":    []interface{}{"biography", "non-fiction"},
					"published": publishedTime,
				}
			},
			expectedState: func(lastUpdated string) map[string]interface{} {
				return map[string]interface{}{
					"@assetType":   "book",
					"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
					"@lastTouchBy": "org1MSP",
					"@lastTx":      "",
					"@lastUpdated": lastUpdated,
					"title":        "Meu Nome é Maria",
					"author":       "Maria Viana",
					"currentTenant": map[string]interface{}{
						"@assetType": "person",
						"@key":       "person:47061146-c642-51a1-844a-bf0b17cb5e19",
					},
					"genres":    []interface{}{"biography", "non-fiction"},
					"published": "2022-05-06T22:12:41Z",
				}
			},
		},
		{
			description: "update books in library",

			assetkey: "library:9c5ffeb3-2491-5a88-858c-653b1ea8dbc5",
			asset: map[string]interface{}{
				"@assetType": "library",
				"name":       "biographies",
				"books": []interface{}{
					map[string]interface{}{
						"@assetType": "book",
						"title":      "Meu Nome é Maria",
						"author":     "Maria Viana",
						"currentTenant": map[string]interface{}{
							"@assetType": "person",
							"name":       "Maria",
							"id":         "31820792048",
							"height":     1.66,
						},
						"genres":    []interface{}{"biography", "non-fiction"},
						"published": "2019-05-06T22:12:41Z",
					},
					map[string]interface{}{
						"@assetType": "book",
						"title":      "Meu Nome é João",
						"author":     "João Viana",
						"currentTenant": map[string]interface{}{
							"@assetType": "person",
							"name":       "João",
							"id":         "42931801159",
							"height":     1.90,
						},
						"genres":    []interface{}{"biography", "non-fiction"},
						"published": "2020-06-06T22:12:41Z",
					},
				},
			},
			updateReq: map[string]interface{}{
				"@assetType": "library",
				"name":       "biographies",
				"books": []interface{}{
					map[string]interface{}{ // This book will not be updated
						"@assetType": "book",
						"@key":       "book:a36a2920-c405-51c3-b584-dcd758338cb5",
					},
					map[string]interface{}{ // This book will be updated
						"@assetType": "book",
						"title":      "Meu Nome é João",
						"author":     "João Viana",
						"currentTenant": map[string]interface{}{ // This person will be updated
							"@assetType": "person",
							"id":         "42931801159",
							"height":     1.92,
						},
						"published": "2020-06-10T22:12:41Z",
					},
				},
			},
			expectedResponse: func(lastUpdated string) map[string]interface{} {
				publishedTimeBook2, _ := time.Parse(time.RFC3339, "2020-06-10T22:12:41Z")

				return map[string]interface{}{
					"@assetType":   "library",
					"@key":         "library:9c5ffeb3-2491-5a88-858c-653b1ea8dbc5",
					"@lastTouchBy": "org1MSP",
					"@lastTx":      "",
					"@lastUpdated": lastUpdated,
					"name":         "biographies",
					"books": []interface{}{
						map[string]interface{}{
							"@assetType":   "book",
							"@key":         "book:a36a2920-c405-51c3-b584-dcd758338cb5",
							"@lastTouchBy": "org1MSP",
							"@lastTx":      "",
							"@lastUpdated": lastUpdated,
							"title":        "Meu Nome é Maria",
							"author":       "Maria Viana",
							"currentTenant": map[string]interface{}{
								"@assetType":   "person",
								"@key":         "person:47061146-c642-51a1-844a-bf0b17cb5e19",
								"@lastTouchBy": "org1MSP",
								"@lastTx":      "",
								"@lastUpdated": lastUpdated,
								"name":         "Maria",
								"id":           "31820792048",
								"height":       1.66,
							},
							"genres":    []interface{}{"biography", "non-fiction"},
							"published": "2019-05-06T22:12:41Z",
						},
						map[string]interface{}{
							"@assetType":   "book",
							"@key":         "book:679f58a4-578f-563c-9c22-3f51b9fab6d5",
							"@lastTouchBy": "org1MSP",
							"@lastTx":      "",
							"@lastUpdated": lastUpdated,
							"title":        "Meu Nome é João",
							"author":       "João Viana",
							"currentTenant": map[string]interface{}{
								"@assetType":   "person",
								"@key":         "person:09c4f266-3bac-5d2f-813b-db3c41ab3375",
								"@lastTouchBy": "org1MSP",
								"@lastTx":      "",
								"@lastUpdated": lastUpdated,
								"name":         "João",
								"id":           "42931801159",
								"height":       1.92,
							},
							"genres":    []interface{}{"biography", "non-fiction"},
							"published": publishedTimeBook2,
						},
					},
				}
			},
			expectedState: func(lastUpdated string) map[string]interface{} {
				return map[string]interface{}{
					"@assetType":   "library",
					"@key":         "library:9c5ffeb3-2491-5a88-858c-653b1ea8dbc5",
					"@lastTouchBy": "org1MSP",
					"@lastTx":      "",
					"@lastUpdated": lastUpdated,
					"name":         "biographies",
					"books": []interface{}{
						map[string]interface{}{
							"@assetType": "book",
							"@key":       "book:a36a2920-c405-51c3-b584-dcd758338cb5",
						},
						map[string]interface{}{
							"@assetType": "book",
							"@key":       "book:679f58a4-578f-563c-9c22-3f51b9fab6d5",
						},
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			stub := mock.NewMockStub("org1MSP", new(testCC))
			sw := &sw.StubWrapper{
				Stub: stub,
			}

			// Put asset
			stub.MockTransactionStart("TestPutAsset")
			var err error
			_, err = assets.PutNewRecursive(sw, tt.asset)
			if err != nil {
				log.Println(err)
				t.FailNow()
			}
			stub.MockTransactionEnd("TestPutAsset")

			// Update asset
			stub.MockTransactionStart("TestUpdateAsset")
			updateResult, err := assets.UpdateRecursive(sw, tt.updateReq)
			if err != nil {
				log.Println(err)
				t.FailNow()
			}

			lastUpdatedTimestamp, _ := stub.GetTxTimestamp()
			lastUpdated := lastUpdatedTimestamp.AsTime().Format(time.RFC3339)
			if !reflect.DeepEqual(tt.expectedResponse(lastUpdated), updateResult) {
				log.Println("these should be deeply equal")
				log.Println(tt.expectedResponse(lastUpdated))
				log.Println(updateResult)
				t.FailNow()
			}
			stub.MockTransactionEnd("TestUpdateAsset")

			// Check state
			stateJSON := stub.State[tt.assetkey]
			var state map[string]interface{}
			err = json.Unmarshal(stateJSON, &state)
			if err != nil {
				log.Println(err)
				t.FailNow()
			}

			if !reflect.DeepEqual(tt.expectedState(lastUpdated), state) {
				log.Println("these should be deeply equal")
				log.Println(tt.expectedState(lastUpdated))
				log.Println(state)
				t.FailNow()
			}
		})
	}
}
