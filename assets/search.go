package assets

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type SearchResponse struct {
	Result   []map[string]interface{}  `json:"result"`
	Metadata *pb.QueryResponseMetadata `json:"metadata"`
}

func Search(stub *sw.StubWrapper, request map[string]interface{}, privateCollection string, resolve bool) (*SearchResponse, errors.ICCError) {
	var bookmark string
	var pageSize int32

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

		if resolve {
			key, err := NewKey(data)
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

	response := SearchResponse{
		Result:   searchResult,
		Metadata: responseMetadata,
	}

	return &response, nil
}
