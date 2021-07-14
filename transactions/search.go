package transactions

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
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
		var bookmark string
		var pageSize int32
		var privateCollection string

		// Check if search is inside a private collection
		privateCollectionInterface, ok := req["collection"]
		if ok {
			privateCollection, ok = privateCollectionInterface.(string)
			if !ok {
				return nil, errors.NewCCError("optional argument 'collection' must be a string", 400)
			}
		}

		requestInterface, ok := req["query"]
		if !ok {
			return nil, errors.NewCCError("argument 'query' is required", 400)
		}
		request, ok := requestInterface.(map[string]interface{})
		if !ok {
			return nil, errors.NewCCError("argument 'query' must be a JSON object", 400)
		}

		// Evaluate special pagination parameters
		bookmarkInt, bookmarkExists := request["bookmark"]
		limit, limitExists := request["limit"]

		// Validate special pagination parameters
		if limitExists {
			limit64, ok := limit.(float64)
			if !ok {
				return nil, errors.NewCCError("limit must be an integer", 400)
			}
			pageSize = int32(limit64)
		}

		if bookmarkExists {
			var ok bool
			bookmark, ok = bookmarkInt.(string)
			if !ok {
				return nil, errors.NewCCError("bookmark must be a string", 400)
			}
		}

		// The "bookmark" and "limit" values are passed as arguments to chaincode API so we delete it from the request
		delete(request, "bookmark")
		delete(request, "limit")

		// Marshal query string
		query, err := json.Marshal(request)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed marshaling JSON-encoded asset", 500)
		}
		queryString := string(query)

		var resultsIterator shim.StateQueryIteratorInterface
		var responseMetadata *pb.QueryResponseMetadata

		if !limitExists {
			// If limit does not exist, search should not be paginated
			if privateCollection == "" {
				resultsIterator, err = stub.GetQueryResult(queryString)
			} else {
				resultsIterator, err = stub.GetPrivateDataQueryResult(privateCollection, queryString)
			}
		} else {
			if privateCollection != "" {
				return nil, errors.NewCCError("private data pagination is not implemented", 501)
			}
			// If it is paginated, call proper API function
			resultsIterator, responseMetadata, err = stub.GetQueryResultWithPagination(queryString, pageSize, bookmark)
		}
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed to get query result", 500)
		}
		defer resultsIterator.Close()

		searchResult := make([]map[string]interface{}, 0)

		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "error iterating response", 500)
			}

			var data map[string]interface{}

			err = json.Unmarshal(queryResponse.Value, &data)
			if err != nil {
				return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal queryResponse values", 500)
			}

			resolve, ok := req["resolve"].(bool)
			if ok && resolve {
				key, err := assets.NewKey(data)
				if err != nil {
					return nil, errors.WrapError(err, "failed to create key object to resolve result")
				}
				asset, err := key.GetRecursive(stub)
				if err != nil {
					return nil, errors.WrapError(err, "failed to resolve result")
				}
				data = asset
			}

			searchResult = append(searchResult, data)
		}

		response := make(map[string]interface{})

		// If query was paginated, responseMetadata is added to
		if responseMetadata != nil {
			response["metadata"] = *responseMetadata
		} else {
			response["metadata"] = make(map[string]string)
		}

		response["result"] = searchResult

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", 500)
		}

		return responseJSON, nil
	},
}
