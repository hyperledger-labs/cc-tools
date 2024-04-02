package querysearch

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger-labs/cc-tools/assets"
	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func (q *QuerySearch) GetFirstKey(stub *sw.StubWrapper) (assets.Key, errors.ICCError) {
	itemMap, err := q.GetFirst(stub)
	if err != nil {
		return nil, err
	}

	if itemMap == nil {
		return nil, nil
	}

	itemKey, err := assets.NewKey(map[string]interface{}{
		"@assetType": itemMap["@assetType"].(string),
		"@key":       itemMap["@key"].(string),
	})

	if err != nil {
		return nil, errors.WrapError(err, "failed to get item key")
	}

	return itemKey, nil
}

func (q *QuerySearch) GetFirst(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {
	result, err := q.getResultsNoPagination(stub)
	if err != nil {
		return nil, err
	}

	data, isExist := result["data"].([]map[string]interface{})
	if !isExist || len(data) == 0 {
		return nil, nil
	}

	fisrtItem := make(map[string]interface{})
	for _, item := range data {
		fisrtItem = item
		break
	}

	return fisrtItem, nil
}

func (q *QuerySearch) GetResults(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {

	if q.config.PageSize > 0 {
		return q.getResultsWithPagination(stub)
	}

	return q.getResultsNoPagination(stub)
}

func (q *QuerySearch) GetResultsByCallBack(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {

	var err error
	if q.config.CallBack == nil {
		return nil, errors.WrapError(nil, "failed, func callback not defined")
	}

	// ** Create results
	res := map[string]interface{}{}

	// ** Assembly query
	q.QueryParser, err = q.Parser()
	if err != nil {
		return nil, errors.WrapError(err, "failed to get query result")
	}

	// ** Execute query
	searchIterator, err := stub.GetQueryResult(q.QueryParser)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get query search result")
	}

	// ** Get records
	for searchIterator.HasNext() {

		searchResult, err := searchIterator.Next()
		if err != nil {
			return nil, errors.WrapError(err, "error iterating result query response")
		}

		var assetMap map[string]interface{}
		err = json.Unmarshal(searchResult.Value, &assetMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to unmarshal order values")
		}

		assetMap, errResolve := q.ResolveExternalAsset(stub, &assetMap)
		if errResolve != nil {
			return nil, errors.WrapErrorWithStatus(errResolve, "failed resolve asset", http.StatusInternalServerError)
		}

		q.removeTags(&assetMap)
		errCallBack := q.config.CallBack(stub, assetMap, res)
		if errCallBack != nil {
			return nil, errors.WrapError(errCallBack, "")
		}
	}

	return res, nil
}

func (q *QuerySearch) GetResultsByCallBackPagination(stub *sw.StubWrapper) (map[string]interface{}, map[string]interface{}, errors.ICCError) {

	var err error
	if q.config.CallBackList == nil {
		return nil, nil, errors.WrapError(nil, "failed, func callback not defined")
	}

	//** Initialize auxiliary variables
	var results []map[string]interface{}
	ctx := map[string]interface{}{}

	// ** Assembly query
	q.QueryParser, err = q.Parser()

	if err != nil {
		return nil, nil, errors.WrapError(err, "failed assemble query search")
	}

	// ** Execute query
	searchIterator, queryResponse, err := stub.GetQueryResultWithPagination(q.QueryParser, q.config.PageSize, q.config.BookMark)
	if err != nil {
		return nil, nil, errors.WrapError(err, "failed to get query search result")
	}

	// ** Get records
	for searchIterator.HasNext() {

		searchResult, err := searchIterator.Next()
		if err != nil {
			return nil, nil, errors.WrapError(err, "error iterating result query response")
		}

		var assetMap map[string]interface{}
		err = json.Unmarshal(searchResult.Value, &assetMap)
		if err != nil {
			return nil, nil, errors.WrapError(err, "failed to unmarshal order values")
		}

		assetMap, errResolve := q.ResolveExternalAsset(stub, &assetMap)
		if errResolve != nil {
			return nil, nil, errors.WrapErrorWithStatus(errResolve, "failed resolve asset", http.StatusInternalServerError)
		}

		data, errCallBack := q.config.CallBackList(stub, assetMap, results, ctx)
		if errCallBack != nil {
			return nil, nil, errors.WrapError(errCallBack, "")
		}

		if data != nil {
			results = append(results, data)
		}
	}

	response := make(map[string]interface{})
	response["data"] = results
	response["bookmark"] = queryResponse.Bookmark
	return response, ctx, nil
}

func (q *QuerySearch) GetIterator(stub *sw.StubWrapper) (shim.StateQueryIteratorInterface, errors.ICCError) {

	// ** Assembly query
	var err error
	q.QueryParser, err = q.Parser()

	if err != nil {
		return nil, errors.WrapError(err, "failed to get query result")
	}

	// ** Execute query
	searchIterator, err := stub.GetQueryResult(q.QueryParser)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get query search result")
	}

	return searchIterator, nil
}

func (q *QuerySearch) getResultsNoPagination(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {

	// ** Create results
	var res []map[string]interface{}

	searchIterator, err := q.GetIterator(stub)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get query search result")
	}

	res, errParse := q.parseResults(stub, searchIterator)
	if errParse != nil {
		return nil, errParse
	}

	response := make(map[string]interface{})
	response["data"] = res
	return response, nil
}

func (q *QuerySearch) getResultsWithPagination(stub *sw.StubWrapper) (map[string]interface{}, errors.ICCError) {

	//** Initialize auxiliary variables
	var results []map[string]interface{}

	// ** Assembly query
	var err error
	q.QueryParser, err = q.Parser()

	if err != nil {
		return nil, errors.WrapError(err, "failed assemble query search")
	}

	//** Get results
	resultsIterator, queryResponse, err := stub.GetQueryResultWithPagination(q.QueryParser, q.config.PageSize, q.config.BookMark)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get query result")
	}

	results, errParse := q.parseResults(stub, resultsIterator)
	if errParse != nil {
		return nil, errParse
	}

	response := make(map[string]interface{})
	response["data"] = results
	response["bookmark"] = queryResponse.Bookmark
	return response, nil
}

func (q *QuerySearch) parseResults(stub *sw.StubWrapper, resultsIterator shim.StateQueryIteratorInterface) ([]map[string]interface{}, errors.ICCError) {

	var results []map[string]interface{}

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error iterating response", http.StatusInternalServerError)
		}

		var data map[string]interface{}

		err = json.Unmarshal(queryResponse.Value, &data)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed to unmarshal queryResponse values", http.StatusInternalServerError)
		}

		data, err = q.ResolveExternalAsset(stub, &data)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "failed resolve asset", http.StatusInternalServerError)
		}

		results = append(results, data)
	}

	return results, nil
}
