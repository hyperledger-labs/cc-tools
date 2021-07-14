package test

import (
	"log"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func TestGetHeader(t *testing.T) {
	stub := shimtest.NewMockStub("org1MSP", new(testCC))

	expectedResponse := map[string]interface{}{
		"ccToolsVersion": "v0.7.0",
		"colors": []interface{}{
			"#4267B2",
			"#34495E",
			"#ECF0F1",
		},
		"name":     "CC Tools Test",
		"orgMSP":   "org1MSP",
		"orgTitle": "CC Tools Demo",
		"version":  "v0.7.0",
	}
	err := invokeAndVerify(stub, "getHeader", nil, expectedResponse, 200)
	if err != nil {
		log.Println("getHeader fail")
		t.FailNow()
	}
}
