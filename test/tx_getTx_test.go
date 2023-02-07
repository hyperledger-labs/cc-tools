package test

import (
	"log"
	"testing"

	"github.com/goledgerdev/cc-tools/mock"
)

func TestGetTx(t *testing.T) {
	stub := mock.NewMockStub("testcc", new(testCC))

	expectedResponse := []interface{}{
		map[string]interface{}{
			"description": "",
			"label":       "Create Asset",
			"tag":         "createAsset",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Update Asset",
			"tag":         "updateAsset",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Delete Asset",
			"tag":         "deleteAsset",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Create Asset Type",
			"tag":         "createAssetType",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Update Asset Type",
			"tag":         "updateAssetType",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Delete Asset Type",
			"tag":         "deleteAssetType",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Load Asset Type List from blockchain",
			"tag":         "loadAssetTypeList",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Get Tx",
			"tag":         "getTx",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Get Header",
			"tag":         "getHeader",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Get Schema",
			"tag":         "getSchema",
		},
		map[string]interface{}{
			"description": "GetDataTypes returns the primary data type map",
			"label":       "Get DataTypes",
			"tag":         "getDataTypes",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Read Asset",
			"tag":         "readAsset",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Read Asset History",
			"tag":         "readAssetHistory",
		},
		map[string]interface{}{
			"description": "",
			"label":       "Search World State",
			"tag":         "search",
		},
	}
	err := invokeAndVerify(stub, "getTx", nil, expectedResponse, 200)
	if err != nil {
		log.Println("getSchema fail")
		t.FailNow()
	}

	req := map[string]interface{}{
		"txName": "getTx",
	}
	expectedGetTx := map[string]interface{}{
		"tag":         "getTx",
		"label":       "Get Tx",
		"description": "",
		"method":      "GET",
		"metaTx":      true,
		"readOnly":    true,
		"args": []interface{}{
			map[string]interface{}{
				"dataType":    "string",
				"description": "The name of the transaction of which you want to fetch the definition. Leave empty to fetch a list of possible transactions.",
				"label":       "",
				"private":     false,
				"required":    false,
				"tag":         "txName",
			},
		},
	}
	err = invokeAndVerify(stub, "getTx", req, expectedGetTx, 200)
	if err != nil {
		log.Println("getSchema fail")
		t.FailNow()
	}
}

func TestGetTx404(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	req := map[string]interface{}{
		"txName": "inexistentTx",
	}

	err := invokeAndVerify(stub, "getTx", req, "transaction named inexistentTx does not exist", 404)
	if err != nil {
		t.FailNow()
	}
}

func TestGetTxInvalidArg(t *testing.T) {
	stub := mock.NewMockStub("org1MSP", new(testCC))

	req := map[string]interface{}{
		"txName": 2,
	}

	err := invokeAndVerify(stub, "getTx", req, "unable to get args: invalid argument 'txName': invalid argument format: property must be a string", 400)
	if err != nil {
		t.FailNow()
	}
}
