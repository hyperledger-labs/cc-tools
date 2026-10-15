package transactions

import (
	"github.com/hyperledger-labs/cc-tools/v2/accesscontrol"
	"github.com/hyperledger-labs/cc-tools/v2/errors"
	sw "github.com/hyperledger-labs/cc-tools/v2/stubwrapper"
)

// Transaction defines the object containing tx definitions
type Transaction struct {
	// List of all MSPs allowed to run this transaction.
	// Regexp is supported by putting '$' before the MSP regexp e.g. []string{`$org\dMSP`}.
	// Please note this restriction DOES NOT protect ledger data from being
	// read by unauthorized organizations, this should be done with Private Data.
	Callers []accesscontrol.Caller `json:"callers,omitempty"`

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
	// A successful Routine always responds with status 200.
	Routine func(*sw.StubWrapper, map[string]interface{}) ([]byte, errors.ICCError) `json:"-"`

	// ResponseRoutine is an alternative to Routine that also sets the success status of the response.
	// Only one of Routine and ResponseRoutine may be defined. The status is only propagated
	// when the chaincode calls RunResponse instead of Run.
	ResponseRoutine func(*sw.StubWrapper, map[string]interface{}) (*Response, errors.ICCError) `json:"-"`
}

// Response is the successful result of a ResponseRoutine.
type Response struct {
	// Status must be a 2xx status. A zero value defaults to 200.
	// It is recorded in the ledger, so it must be deterministic across endorsers.
	Status int32

	Payload []byte
}
