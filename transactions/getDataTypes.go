package transactions

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GetDataTypes returns the data type map
var GetDataTypes = Transaction{
	Tag:         "getDataTypes",
	Label:       "Get DataTypes",
	Description: "GetDataTypes returns the primary data type map",
	Method:      "GET",

	ReadOnly: true,
	MetaTx:   true,
	Args:     []Argument{},
	Routine: func(stub shim.ChaincodeStubInterface, req map[string]interface{}) ([]byte, errors.ICCError) {
		dataTypeMap := assets.DataTypeMap()

		dataTypeMapBytes, err := json.Marshal(dataTypeMap)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling data type map", 500)
		}
		return dataTypeMapBytes, nil
	},
}
