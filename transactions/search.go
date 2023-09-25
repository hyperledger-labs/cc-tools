package transactions

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// Search makes a rich query against CouchDB
var Search = Transaction{
	Tag:         "search",
	Label:       "Search World State",
	Description: "",
	Method:      "GET",

	MetaTx: true,
	Args: ArgList{
		{
			Tag:         "query",
			Description: "Query string according to CouchDB specification: https://docs.couchdb.org/en/stable/api/database/find.html.",
			DataType:    "@query",
			Required:    true,
		},
		{
			Tag:         "collection",
			Description: "Name of the private collection to be searched.",
			DataType:    "string",
		},
		{
			Tag:         "resolve",
			Description: "Resolve references recursively.",
			DataType:    "boolean",
		},
	},
	ReadOnly: true,
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var err error

		query, _ := req["query"].(map[string]interface{})

		// Check if search is inside a private collection
		privateCollection, _ := req["collection"].(string)

		resolve, _ := req["resolve"].(bool)

		response, err := assets.Search(stub, query, privateCollection, resolve)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "query error", 500)
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", 500)
		}

		return responseJSON, nil
	},
}
