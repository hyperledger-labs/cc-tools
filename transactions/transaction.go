package transactions

import (
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// Transaction defines the object containing tx definitions
type Transaction struct {
	// List of all MSPs allowed to run this transaction.
	// Regexp is supported by putting '$' before the MSP regexp e.g. []string{`$org\dMSP`}.
	// Please note this restriction DOES NOT protect ledger data from being
	// read by unauthorized organizations, this should be done with Private Data.
	Callers []string `json:"callers,omitempty"`

	// Tag is how the tx will be called.
	Tag string `json:"tag"`

	// Label is the pretty tx name for front-end rendering.
	Label string `json:"label"`

	// Description is a simple explanation describing what the tx does.
	Description string `json:"description"`

	// Args is a list of argument formats accepted by the tx.
	Args ArgList `json:"args"`

	// Method indicates the HTTP method which should be used to call the tx when using an HTTP API.
	Method string `json:"method"`

	// ReadOnly indicates that the tx does not alter the world state.
	ReadOnly bool `json:"readOnly"`

	// MetaTx indicates that the tx does not encode a business-specific rule,
	// but an internal process of the chaincode e.g. listing available asset types.
	MetaTx bool `json:"metaTx"`

	// Routine is the function called when running the tx. It is where the tx logic can be programmed.
	Routine func(*sw.StubWrapper, map[string]interface{}) ([]byte, errors.ICCError) `json:"-"`
}
