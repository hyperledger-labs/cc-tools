package transactions

import (
	b64 "encoding/base64"
	"net/http"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/goledgerdev/cc-tools/events"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// ExecuteEvent executes an event of type "EventCustom"
var ExecuteEvent = Transaction{
	Tag:         "executeEvent",
	Label:       "Execute Event",
	Description: "ExecuteEvent executes an event of type 'EventCustom'",
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

		event.CustomFunction(payload)

		return nil, nil
	},
}