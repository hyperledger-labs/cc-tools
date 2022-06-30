package assets

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

//
// Get all history from asset
//
type HistoryResponse struct {
	Result []map[string]interface{} `json:"result"`
}

func GetHistoryForKey(stub *sw.StubWrapper, keyAsset string, resolve bool) (*HistoryResponse, errors.ICCError) {

	if len(keyAsset) == 0 {
		return nil, errors.WrapErrorWithStatus(nil, "keyAsset is empty", 500)
	}

	var err error
	var historyIterator shim.HistoryQueryIteratorInterface

	historyIterator, err = stub.GetHistoryForKey(keyAsset)
	if err != nil {
		return nil, errors.WrapError(err, "failed to read asset from blockchain")
	}

	// Result
	historyResult := make([]map[string]interface{}, 0)

	// Return is empty
	if historyIterator == nil {
		response := HistoryResponse{
			Result: historyResult,
		}
		return &response, nil
	}

	defer historyIterator.Close()

	for historyIterator.HasNext() {

		queryResponse, err := historyIterator.Next()
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error iterating response", 500)
		}

		var data = make(map[string]interface{})
		if queryResponse.IsDelete {
			continue
		}

		err = json.Unmarshal(queryResponse.Value, &data)
		if err != nil {
			return nil, errors.WrapError(err, "failed to unmarshal queryResponse's values")
		}

		// If active get the current state of the references assets
		if resolve {

			for k, v := range data {

				switch prop := v.(type) {

				// Direct reference
				case map[string]interface{}:
					propKey, err := NewKey(prop)
					if err != nil {
						return nil, errors.WrapErrorWithStatus(err, "failed to resolve asset references", 500)
					}

					if subAsset, err := propKey.GetRecursive(stub); err != nil {
						return nil, errors.WrapErrorWithStatus(err, "failed to resolve asset references", 500)
					} else {
						data[k] = subAsset
					}

				// needs to be tested
				// Arrays of references
				case []interface{}:
					for idx, elem := range prop {
						if elemMap, ok := elem.(map[string]interface{}); ok {
							elemKey, err := NewKey(elemMap)
							if err != nil {
								return nil, errors.WrapErrorWithStatus(err, "failed to resolve asset references", 500)
							}

							var subAsset map[string]interface{}
							subAsset, err = elemKey.GetRecursive(stub)
							if err != nil {
								return nil, errors.WrapErrorWithStatus(err, "failed to get subasset", 500)
							}
							prop[idx] = subAsset
						}
					}
				}
			}
		}

		historyResult = append(historyResult, data)
	}

	response := HistoryResponse{
		Result: historyResult,
	}

	return &response, nil
}
