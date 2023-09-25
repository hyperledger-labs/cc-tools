package assets

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type HistoryResponse struct {
	Result   []map[string]interface{}  `json:"result"`
	Metadata *pb.QueryResponseMetadata `json:"metadata"`
}

func History(stub *sw.StubWrapper, key string, resolve bool) (*HistoryResponse, errors.ICCError) {
	var resultsIterator shim.HistoryQueryIteratorInterface

	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return nil, errors.WrapErrorWithStatus(err, "failed to get history for key", http.StatusInternalServerError)
	}
	defer resultsIterator.Close()

	historyResult := make([]map[string]interface{}, 0)

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

		historyResult = append(historyResult, data)
	}

	response := HistoryResponse{
		Result: historyResult,
	}

	return &response, nil
}
