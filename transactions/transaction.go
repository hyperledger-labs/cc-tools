package transactions

import (
	"github.com/goledgerdev/cc-tools/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Transaction defines the object containing tx definitions
type Transaction struct {
	// List of all MSPs allowed to run this transaction.
	// Regexp is supported by putting '$' before the MSP regexp e.g. []string{`$org\dMSP`}.
	// Please note this restriction DOES NOT protect ledger data from being read by unauthorized organizations.
	// This should be done with Private Data.
	Callers []string `json:"callers,omitempty"`

	Tag         string `json:"tag"`
	Label       string `json:"label"`
	Description string `json:"description"`

	Args     map[string]Argument `json:"args"`
	Method   string              `json:"method"`
	ReadOnly bool                `json:"readOnly"`
	MetaTx   bool                `json:"metaTx"`

	Routine func(shim.ChaincodeStubInterface, map[string]interface{}) ([]byte, errors.ICCError) `json:"-"`
}
