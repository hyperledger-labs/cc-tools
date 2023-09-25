package events

import (
	"fmt"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

type EventType float64

const (
	EventLog EventType = iota
	EventTransaction
	EventCustom
)

// Event is the struct defining a primitive event.
type Event struct {
	// Tag is how the event will be referenced
	Tag string `json:"tag"`

	// Label is the pretty event name for logs
	Label string `json:"label"`

	// Description is a simple explanation describing the meaning of the event.
	Description string `json:"description"`

	// BaseLog is the basisc log message for the event
	BaseLog string `json:"baseLog"`

	// Type is the type of event
	Type EventType `json:"type"`

	// Receivers is an array that specifies which organizations will receive the event.
	// Accepts either basic strings for exact matches
	// eg. []string{'org1MSP', 'org2MSP'}
	// or regular expressions
	// eg. []string{`$org\dMSP`} and cc-tools will
	// check for a match with regular expression `org\dMSP`
	Receivers []string `json:"receivers,omitempty"`

	// Transaction is the transaction that the event triggers (if of type EventTransaction)
	Transaction string `json:"transaction"`

	// Channel is the channel of the transaction that the event triggers (if of type EventTransaction)
	// If empty, the event will trigger on the same channel as the transaction that calls the event
	Channel string `json:"channel"`

	// Chaincode is the chaincode of the transaction that the event triggers (if of type EventTransaction)
	// If empty, the event will trigger on the same chaincode as the transaction that calls the event
	Chaincode string `json:"chaincode"`

	// CustomFunction is used an event of type "EventCustom" is called.
	// It is a function that receives a stub and a payload and returns an error.
	CustomFunction func(*sw.StubWrapper, []byte) error `json:"-"`

	// ReadOnly indicates if the CustomFunction has the ability to alter the world state (if of type EventCustom).
	ReadOnly bool `json:"readOnly"`
}

func (event Event) CallEvent(stub *sw.StubWrapper, payload []byte) errors.ICCError {
	err := stub.SetEvent(event.Tag, payload)
	if err != nil {
		return errors.WrapError(err, "stub.SetEvent call error")
	}

	return nil
}

func CallEvent(stub *sw.StubWrapper, eventTag string, payload []byte) errors.ICCError {
	event := FetchEvent(eventTag)
	if event == nil {
		return errors.NewCCError(fmt.Sprintf("event named %s does not exist", eventTag), http.StatusBadRequest)
	}

	return event.CallEvent(stub, payload)
}
