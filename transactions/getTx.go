package transactions

import (
	"encoding/json"
	"fmt"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// getTx returns tx definitions
var getTx = Transaction{
	Tag:         "getTx",
	Label:       "Get Tx",
	Description: "",
	Method:      "GET",

	ReadOnly: true,
	MetaTx:   true,
	Args: map[string]Argument{
		"txName": {
			DataType:    "string",
			Description: "The name of the transaction of which you want to fetch the definition. Leave empty to fetch a list of possible transactions.",
		},
	},
	Routine: func(stub shim.ChaincodeStubInterface, req map[string]interface{}) ([]byte, errors.ICCError) {
		var txName string

		txNameInterface, ok := req["txName"]
		if ok {
			txName, ok = txNameInterface.(string)
			if !ok {
				return nil, errors.NewCCError(fmt.Sprintf("argument 'txName' must be a string"), 400)
			}
		}

		txList := TxList()

		// If user requested a specific transaction definition
		if txName != "" {
			txDef := FetchTx(txName)
			if txDef == nil {
				errMsg := fmt.Sprintf("transaction named %s does not exist", txName)
				return nil, errors.NewCCError(errMsg, 400)
			}
			txDefBytes, err := json.Marshal(txDef)
			if err != nil {
				errMsg := fmt.Sprintf("error marshaling transaction definition: %s", err)
				return nil, errors.NewCCError(errMsg, 400)
			}
			return txDefBytes, nil
		}

		// If user requested asset list
		type txListElem struct {
			Tag         string   `json:"tag"`
			Label       string   `json:"label"`
			Description string   `json:"description"`
			Callers     []string `json:"callers,omitempty"`
		}
		var txRetList []txListElem
		for _, tx := range txList {
			txRetList = append(txRetList, txListElem{
				Tag:         tx.Tag,
				Label:       tx.Label,
				Description: tx.Description,
				Callers:     tx.Callers,
			})
		}

		txListBytes, err := json.Marshal(txRetList)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling transaction list", 500)
		}
		return txListBytes, nil
	},
}
