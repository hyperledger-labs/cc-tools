package transactions

import (
	b64 "encoding/base64"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/errors"
	"github.com/hyperledger-labs/cc-tools/events"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// RunEvent runs an event of type "EventCustom" as readOnly
var RunEvent = Transaction{
	Tag:         "runEvent",
	Label:       "Run Event",
	Description: "RunEvent runs an event of type 'EventCustom' as readOnly",
	Method:      "GET",

	ReadOnly: true,
	MetaTx:   true,
	Args: ArgList{
		{
			Tag:         "eventTag",
			Description: "Event Tag",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "payload",
			Description: "Paylod in Base64 encoding",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var payload []byte

		payloadEncoded, ok := req["payload"].(string)
		if ok {
			payloadDecoded, er := b64.StdEncoding.DecodeString(payloadEncoded)
			if er != nil {
				return nil, errors.WrapError(er, "error decoding payload")
			}

			payload = payloadDecoded
		}
		event := events.FetchEvent(req["eventTag"].(string))

		if event.Type != events.EventCustom {
			return nil, errors.NewCCError("event is not of type 'EventCustom'", http.StatusBadRequest)
		}

		err := event.CustomFunction(stub, payload)
		if err != nil {
			return nil, errors.WrapError(err, "error executing custom function")
		}

		return nil, nil
	},
}
