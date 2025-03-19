package querysearch

import (
	"encoding/json"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

// *** Simple count
func (q *QuerySearch) GetCount(stub *sw.StubWrapper) (float64, error) {

	// ** Assembly query
	query, err := q.Parser()
	if err != nil {
		return 0, errors.WrapError(err, "failed to get query result")
	}

	// ** Execute query
	searchIterator, err := stub.GetQueryResult(query)
	if err != nil {
		return 0, errors.WrapError(err, "failed to get query search result")
	}

	var count float64

	//** Get records
	for searchIterator.HasNext() {

		_, err := searchIterator.Next()
		if err != nil {
			return 0, errors.WrapError(err, "error iterating order query response")
		}

		count++
	}

	return count, nil
}

// *** Simple Sum
func (q *QuerySearch) GetSum(stub *sw.StubWrapper, propertyNameForSum string) (float64, error) {

	// ** Assembly query
	query, err := q.Parser()
	if err != nil {
		return 0, errors.WrapError(err, "failed to get query result")
	}

	// ** Execute query
	searchIterator, err := stub.GetQueryResult(query)
	if err != nil {
		return 0, errors.WrapError(err, "failed to get query search result")
	}

	var sum float64

	//** Get records
	for searchIterator.HasNext() {

		searchResult, err := searchIterator.Next()
		if err != nil {
			return 0, errors.WrapError(err, "error iterating order query response")
		}

		var searchMap map[string]interface{}
		err = json.Unmarshal(searchResult.Value, &searchMap)
		if err != nil {
			return 0, errors.WrapError(err, "failed to unmarshal searchResult values")
		}

		if value, has := searchMap[propertyNameForSum].(float64); has {
			sum += value
		}
	}

	return sum, nil
}
