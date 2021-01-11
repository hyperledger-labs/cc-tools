package assets

import (
	"fmt"
	"reflect"
	"testing"
)

func testParseValid(t *testing.T, dtype DataType, inputVal interface{}, expectedKey string, expectedVal interface{}) {
	var key string
	var val interface{}
	var err error
	key, val, err = dtype.Parse(inputVal)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if key != expectedKey {
		fmt.Printf("parsing %v expected key: %q but got %q\n", inputVal, expectedKey, key)
		t.FailNow()
	}
	if val != expectedVal {
		fmt.Printf("parsing %v expected parsed val: \"%v\" of type %s but got \"%v\" of type %s\n", inputVal, expectedVal, reflect.TypeOf(expectedVal), val, reflect.TypeOf(val))
		t.FailNow()
	}
}

func testParseInvalid(t *testing.T, dtype DataType, inputVal interface{}, expectedErr int32) {
	_, _, err := dtype.Parse(inputVal)
	if err == nil {
		fmt.Println("expected error but DataType.Parse was successful")
		t.FailNow()
	}
	if err.Status() != expectedErr {
		fmt.Printf("expected error code %d but got %d\n", expectedErr, err.Status())
	}
}

func TestDataTypeString(t *testing.T) {
	dtypeName := "string"
	dtype, exists := DataTypeMap()[dtypeName]
	if !exists {
		fmt.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
		t.FailNow()
	}
	testParseValid(t, dtype, "string válida", "string válida", "string válida")
	testParseInvalid(t, dtype, 32.0, 400)
}

func TestDataTypeNumber(t *testing.T) {
	dtypeName := "number"
	dtype, exists := DataTypeMap()[dtypeName]
	if !exists {
		fmt.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
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
	dtype, exists := DataTypeMap()[dtypeName]
	if !exists {
		fmt.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
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
	dtype, exists := DataTypeMap()[dtypeName]
	if !exists {
		fmt.Printf("%s datatype not declared in DataTypeMap\n", dtypeName)
		t.FailNow()
	}
	testParseValid(t, dtype, true, "t", true)
	testParseValid(t, dtype, false, "f", false)
	testParseValid(t, dtype, "true", "t", true)
	testParseValid(t, dtype, "false", "f", false)
	testParseInvalid(t, dtype, "True", 400)
	testParseInvalid(t, dtype, 37.3, 400)
}
