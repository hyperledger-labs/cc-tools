package transactions

import (
	"encoding/json"
	"net/http"

	"github.com/goledgerdev/cc-tools/errors"
	"github.com/goledgerdev/cc-tools/events"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
)

// GetEvents returns the events map
var GetEvents = Transaction{
	Tag:         "getEvents",
	Label:       "Get Events",
	Description: "GetEvents returns the events map",
	Method:      "GET",

	ReadOnly: true,
	MetaTx:   true,
	Args:     ArgList{},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		eventList := events.EventList()

		eventListBytes, err := json.Marshal(eventList)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling event list", http.StatusInternalServerError)
		}
		return eventListBytes, nil
	},
}
