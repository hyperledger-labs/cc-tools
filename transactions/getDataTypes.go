package transactions

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// GetDataTypes returns the primitive data type map
var GetDataTypes = Transaction{
	Tag:         "getDataTypes",
	Label:       "Get DataTypes",
	Description: "GetDataTypes returns the primary data type map",
	Method:      "GET",

	ReadOnly: true,
	MetaTx:   true,
	Args:     ArgList{},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		dataTypeMap := assets.DataTypeMap()

		dataTypeMapBytes, err := json.Marshal(dataTypeMap)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling data type map", 500)
		}
		return dataTypeMapBytes, nil
	},
}
