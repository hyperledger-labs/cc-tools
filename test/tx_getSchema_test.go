package test

import (
	"log"
	"testing"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestGetSchema(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	expectedResponse := []interface{}{
		map[string]interface{}{
			"description": "Personal data of someone",
			"label":       "Person",
			"tag":         "person",
			"writers":     nil,
		},
		map[string]interface{}{
			"description": "Library as a collection of books",
			"label":       "Library",
			"tag":         "library",
			"writers":     nil,
		},
		map[string]interface{}{
			"description": "Book",
			"label":       "Book",
			"tag":         "book",
			"writers":     nil,
		},
		map[string]interface{}{
			"description": "Secret between Org2 and Org3",
			"label":       "Secret",
			"readers": []interface{}{
				"org2MSP",
				"org3MSP",
			},
			"tag":     "secret",
			"writers": nil,
		},
		map[string]interface{}{
			"description": "AssetTypeListData",
			"label":       "AssetTypeListData",
			"tag":         "assetTypeListData",
			"writers":     nil,
		},
	}
	err := invokeAndVerify(stub, "getSchema", nil, expectedResponse, 200)
	if err != nil {
		log.Println("getSchema fail")
		t.FailNow()
	}

	req := map[string]interface{}{
		"assetType": "person",
	}
	expectedPersonSchema := map[string]interface{}{
		"tag":         "person",
		"label":       "Person",
		"description": "Personal data of someone",
		"props": []interface{}{
			map[string]interface{}{
				"dataType":    "cpf",
				"description": "",
				"isKey":       true,
				"label":       "CPF (Brazilian ID)",
				"readOnly":    false,
				"required":    true,
				"tag":         "id",
				"writers": []interface{}{
					"org1MSP",
				},
			},
			map[string]interface{}{
				"dataType":    "string",
				"description": "",
				"isKey":       false,
				"label":       "Name of the person",
				"readOnly":    false,
				"required":    true,
				"tag":         "name",
				"writers":     nil,
			},
			map[string]interface{}{
				"dataType":    "datetime",
				"description": "",
				"isKey":       false,
				"label":       "Date of Birth",
				"readOnly":    false,
				"required":    false,
				"tag":         "dateOfBirth",
				"writers": []interface{}{
					"org1MSP",
				},
			},
			map[string]interface{}{
				"dataType":     "number",
				"defaultValue": 0.0,
				"description":  "",
				"isKey":        false,
				"label":        "Person's height",
				"readOnly":     false,
				"required":     false,
				"tag":          "height",
				"writers":      nil,
			},
			map[string]interface{}{
				"dataType":    "@object",
				"description": "",
				"isKey":       false,
				"label":       "Other Info",
				"readOnly":    false,
				"required":    false,
				"tag":         "info",
				"writers":     nil,
			},
		},
	}
	err = invokeAndVerify(stub, "getSchema", req, expectedPersonSchema, 200)
	if err != nil {
		log.Println("getSchema of person fail")
		t.FailNow()
	}
}

func TestGetSchema404(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	req := map[string]interface{}{
		"assetType": "inexistentAsset",
	}

	err := invokeAndVerify(stub, "getSchema", req, "asset type named inexistentAsset does not exist", 404)
	if err != nil {
		t.FailNow()
	}
}
