package test

import (
	"log"
	"testing"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestGetDataTypes(t *testing.T) {
	stub := mock.NewMockStub("testcc", new(testCC))

	expectedResponse := map[string]interface{}{
		"boolean": map[string]interface{}{
			"acceptedFormats": []interface{}{
				"boolean",
			},
			"DropDownValues": nil,
		},
		"cpf": map[string]interface{}{
			"acceptedFormats": nil,
			"DropDownValues":  nil,
		},
		"datetime": map[string]interface{}{
			"acceptedFormats": []interface{}{
				"string",
			},
			"DropDownValues": nil,
		},
		"integer": map[string]interface{}{
			"acceptedFormats": []interface{}{
				"number",
			},
			"DropDownValues": nil,
		},
		"number": map[string]interface{}{
			"acceptedFormats": []interface{}{
				"number",
			},
			"DropDownValues": nil,
		},
		"string": map[string]interface{}{
			"acceptedFormats": []interface{}{
				"string",
			},
			"DropDownValues": nil,
		},
		"@object": map[string]interface{}{
			"acceptedFormats": []interface{}{
				"@object",
			},
			"DropDownValues": nil,
		},
	}
	err := invokeAndVerify(stub, "getDataTypes", nil, expectedResponse, 200)
	if err != nil {
		log.Println("getDataTypes fail")
		t.FailNow()
	}
}
