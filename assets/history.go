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
	var subAssets []AssetProp

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
			if subAssets == nil {
				key, err := NewKey(data)
				if err != nil {
					return nil, errors.WrapError(err, "failed to create key object to resolve result")
				}
				subAssets = key.Type().SubAssets()
			}

			err := resolveHistory(stub, data, subAssets)
			if err != nil {
				return nil, errors.WrapError(err, "failed to resolve result")
			}
		}

		historyResult = append(historyResult, data)
	}

	response := HistoryResponse{
		Result: historyResult,
	}

	return &response, nil
}

func resolveHistory(stub *sw.StubWrapper, data map[string]interface{}, subAssets []AssetProp) errors.ICCError {
	for _, refProp := range subAssets {
		ref, ok := data[refProp.Tag].(map[string]interface{})
		if !ok {
			continue
		}

		key, err := NewKey(ref)
		if err != nil {
			return errors.WrapError(err, "could not make subasset key")
		}

		resolved, err := key.GetRecursive(stub)
		if err != nil {
			return errors.WrapError(err, "failed to get subasset recursive")
		}

		data[refProp.Tag] = resolved
	}

	return nil
}
