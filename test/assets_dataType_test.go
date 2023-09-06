package test

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/goledgerdev/cc-tools/assets"
)

func testParseValid(t *testing.T, dtype assets.DataType, inputVal interface{}, expectedKey string, expectedVal interface{}) {
	var key string
	var val interface{}
	var err error
	key, val, err = dtype.Parse(inputVal)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if key != expectedKey {
		log.Printf("parsing %v expected key: %q but got %q\n", inputVal, expectedKey, key)
		t.FailNow()
	}
	if !reflect.DeepEqual(val, expectedVal) {
		log.Printf("parsing %v expected parsed val: \"%v\" of type %s but got \"%v\" of type %s\n", inputVal, expectedVal, reflect.TypeOf(expectedVal), val, reflect.TypeOf(val))
		t.FailNow()
	}
}

func testParseInvalid(t *testing.T, dtype assets.DataType, inputVal interface{}, expectedErr int32) {
	_, _, err := dtype.Parse(inputVal)
	if err == nil {
		log.Println("expected error but DataType.Parse was successful")
		t.FailNow()
	}
	if err.Status() != expectedErr {
		log.Printf("expected error code %d but got %d\n", expectedErr, err.Status())
	}
}

func TestDataTypeString(t *testing.T) {
	dtypeName := "string"
	dtype, exists := assets.DataTypeMap()[dtypeName]
	if !exists {
		log.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
		t.FailNow()
	}
	testParseValid(t, dtype, "string válida", "string válida", "string válida")
	testParseInvalid(t, dtype, 32.0, 400)
}

func TestDataTypeNumber(t *testing.T) {
	dtypeName := "number"
	dtype, exists := assets.DataTypeMap()[dtypeName]
	if !exists {
		log.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
		t.FailNow()
	}
	testParseValid(t, dtype, 472.7, "407d8b3333333333", 472.7)
	testParseValid(t, dtype, 472, "407d800000000000", 472.0)
	testParseValid(t, dtype, "472", "407d800000000000", 472.0)
	testParseInvalid(t, dtype, "32d.0", 400)
	testParseInvalid(t, dtype, false, 400)
}

func TestDataTypeInteger(t *testing.T) {
	dtypeName := "integer"
	dtype, exists := assets.DataTypeMap()[dtypeName]
	if !exists {
		log.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
		t.FailNow()
	}
	testParseValid(t, dtype, 470, "470", int64(470))
	testParseValid(t, dtype, 412.0, "412", int64(412))
	testParseValid(t, dtype, "472", "472", int64(472))
	testParseInvalid(t, dtype, 472.1, 400)
	testParseInvalid(t, dtype, "32d.0", 400)
	testParseInvalid(t, dtype, false, 400)
}

func TestDataTypeBoolean(t *testing.T) {
	dtypeName := "boolean"
	dtype, exists := assets.DataTypeMap()[dtypeName]
	if !exists {
		log.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
		t.FailNow()
	}
	testParseValid(t, dtype, true, "t", true)
	testParseValid(t, dtype, false, "f", false)
	testParseValid(t, dtype, "true", "t", true)
	testParseValid(t, dtype, "false", "f", false)
	testParseInvalid(t, dtype, "True", 400)
	testParseInvalid(t, dtype, 37.3, 400)
}

func TestDataTypeObject(t *testing.T) {
	dtypeName := "@object"
	dtype, exists := assets.DataTypeMap()[dtypeName]
	if !exists {
		log.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
		t.FailNow()
	}

	testCase1 := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	testCaseByte1, _ := json.Marshal(testCase1)
	testCaseExpected1 := map[string]interface{}{
		"@assetType": "@object",
		"key1":       "value1",
		"key2":       "value2",
	}
	testCaseExpectedByte1, _ := json.Marshal(testCaseExpected1)

	testParseValid(t, dtype, testCase1, string(testCaseExpectedByte1), testCase1)
	testParseValid(t, dtype, testCaseByte1, string(testCaseExpectedByte1), testCase1)
	testParseValid(t, dtype, string(testCaseByte1), string(testCaseExpectedByte1), testCase1)
	testParseInvalid(t, dtype, "{'key': 'value'}", http.StatusBadRequest)
}
func TestDataTypeAsset(t *testing.T) {
	dtypeName := "@asset"
	dtype, exists := assets.DataTypeMap()[dtypeName]
	if !exists {
		log.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
		t.FailNow()
	}

	testCase1 := map[string]interface{}{
		"@assetType": "person",
		"id":         "42186475006",
	}
	testCaseExpected1 := map[string]interface{}{
		"@assetType": "person",
		"id":         "42186475006",
		"@key":       "person:a11e54a8-7e23-5d16-9fed-45523dd96bfa",
	}
	testCaseExpectedByte1, _ := json.Marshal(testCaseExpected1)
	testParseValid(t, dtype, testCase1, string(testCaseExpectedByte1), testCaseExpected1)

	testCase2 := map[string]interface{}{
		"@assetType": "book",
		"title":      "Book Name",
		"author":     "Author Name",
		"@key":       "book:983a78df-9f0e-5ecb-baf2-4a8698590c81",
	}
	testCaseExpectedByte2, _ := json.Marshal(testCase2)
	testParseValid(t, dtype, testCase2, string(testCaseExpectedByte2), testCase2)
	testParseValid(t, dtype, testCaseExpectedByte2, string(testCaseExpectedByte2), testCase2)
	testParseValid(t, dtype, string(testCaseExpectedByte2), string(testCaseExpectedByte2), testCase2)

	testCase3 := map[string]interface{}{
		"@assetType": "library",
	}
	testParseInvalid(t, dtype, testCase3, http.StatusBadRequest)

	testCase4 := map[string]interface{}{
		"@assetType": "inexistant",
	}
	testParseInvalid(t, dtype, testCase4, http.StatusBadRequest)
}
