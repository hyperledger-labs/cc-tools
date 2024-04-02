package querysearch

import (
	"log"
	"net/http"
	"strings"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

func (q *QuerySearch) resolveAsset(stub *sw.StubWrapper, elemMap map[string]interface{}, subResolve *[]string) (map[string]interface{}, errors.ICCError) {
	assetKey, err := assets.NewKey(elemMap)
	if err != nil {
		return nil, errors.NewCCError("failed to get generate key", http.StatusBadRequest)
	}

	assetMap, err := assetKey.GetMap(stub)
	if err != nil {
		return nil, errors.NewCCError("failed to get asset in ledger", http.StatusBadRequest)
	}

	q.removeTags(&assetMap)

	// ** CAUTION **
	// ** This solutions is low performance
	if len(*subResolve) > 0 {
		for _, subKey := range *subResolve {

			subKey = strings.TrimSpace(subKey)
			if subAsset, ok := assetMap[subKey]; ok {
				if subAssetMap, ok := subAsset.(map[string]interface{}); ok {
					subAssetKey, err := assets.NewKey(subAssetMap)
					if err != nil {
						return nil, errors.NewCCError("failed to generate key for sub asset", http.StatusBadRequest)
					}

					subAssetMap, err := subAssetKey.GetMap(stub)
					if err != nil {
						return nil, errors.NewCCError("failed to get sub asset in ledger", http.StatusBadRequest)
					}
					q.removeTags(&subAssetMap)
					assetMap[subKey] = subAssetMap
				}
			}
		}
	}

	return assetMap, nil
}

func (q *QuerySearch) resolveSlice(stub *sw.StubWrapper, prop []interface{}, subResolve []string) ([]interface{}, errors.ICCError) {

	// ** Revolve assets into slice
	//! ** There can be 2 types of slices
	// ** Type 1
	//! **	- reference slice for asset
	// ** Type 2
	//! ** 	- slice containing a map with @key and other data at the same level

	lstResult := make([]interface{}, 0)
	for _, elem := range prop {

		elemMap, ok := elem.(map[string]interface{})

		// ** Not a map or not subResolve, append to result and continue
		if !ok {
			lstResult = append(lstResult, elem)
			continue
		}

		// **
		// ** TYPE 1 **
		// **
		// ** Check if element is a asset
		if _, ok := elemMap["@key"]; ok {
			res, err := q.resolveAsset(stub, elemMap, &subResolve)
			if err != nil {
				return nil, err
			}

			lstResult = append(lstResult, res)
			continue
		}

		// **
		// ** TYPE 2 **
		// **
		itemMap := make(map[string]interface{})
		for key, value := range elemMap {

			// ** check if value is a map of asset
			if _, ok := value.(map[string]interface{}); ok {

				assetMap, err := q.resolveAsset(stub, value.(map[string]interface{}), &subResolve)
				if err != nil {
					return nil, err
				}

				//! ** add key (with resolved asset) to itemMap
				itemMap[key] = assetMap
				continue
			}

			//! ** add key OTHER type ( no asset ) to itemMap
			itemMap[key] = value
		}

		lstResult = append(lstResult, itemMap)
	}

	return lstResult, nil
}

func (q *QuerySearch) ResolveExternalAsset(stub *sw.StubWrapper, m *map[string]interface{}) (map[string]interface{}, errors.ICCError) {

	data := *m

	q.removeTags(&data)

	// ** Resolve asset external
	for _, resolveStr := range q.config.Resolve {
		// Check exist sub asset for resolve
		var resolve string
		var subResolve []string

		tuplaKeys := strings.Split(resolveStr, ".")

		if len(tuplaKeys) > 1 {
			resolve = tuplaKeys[0]

			subAssetStr := strings.ReplaceAll(tuplaKeys[1], "{", "")
			subAssetStr = strings.ReplaceAll(subAssetStr, "}", "")

			//** create list asset to resolve
			subResolve = strings.Split(subAssetStr, ",")
		} else {
			resolve = resolveStr
		}

		key, has := data[resolve]

		if has {
			switch prop := key.(type) {
			case map[string]interface{}:
				{
					assetKey, err := assets.NewKey(prop)
					if err != nil {
						return nil, errors.NewCCError("failed to get asset in ledger", http.StatusBadRequest)
					}

					assetMap, err := assetKey.GetRecursive(stub)
					if err != nil {
						return nil, errors.NewCCError("failed to get asset in ledger", http.StatusBadRequest)
					}
					q.removeTags(&assetMap)
					data[resolve] = assetMap
				}

			case []interface{}:
				{
					res, err := q.resolveSlice(stub, prop, subResolve)
					if err != nil {
						return nil, errors.NewCCError(err.Error(), http.StatusBadRequest)
					}
					data[resolve] = res
				}
			default:
				log.Printf("resolveExternalAsset typeKey: %T", key)
			}
		}
	}
	return data, nil
}
