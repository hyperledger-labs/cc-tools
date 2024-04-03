// This package provides a way to query a database, possibly using a CouchDB-like interface,
// and defines several types, methods, and functions for querying. It facilitates the creation of queries
// without having to write queries in a raw form.
//
// Configuration
// -------------
// The Config struct contains the settings used for querying.
// These settings include RemoveTags (a list of tags to be removed from the results),
// AssetName (the name of the asset being queried), PageSize (the number of results per page),
// BookMark (the bookmark for pagination), Resolve (a list of relations to be resolved),
// Sort (a list of fields to sort by), IndexDoc (a document containing the index design and name),
// NoRemoveTagsTransaction (a boolean indicating whether to remove tags during the transaction),
// and CallBack (a function to be called after the query is executed).
package querysearch

import (
	"encoding/json"
	"time"

	"github.com/hyperledger-labs/cc-tools/errors"
	sw "github.com/hyperledger-labs/cc-tools/stubwrapper"
)

type Config struct {
	RemoveTags              []string
	removeDefaultTags       []string
	AssetName               string
	PageSize                int32
	BookMark                string
	Resolve                 []string
	Sort                    []map[string]string
	IndexDoc                IndexDocument
	NoRemoveTagsTransaction bool
	CallBack                func(*sw.StubWrapper, map[string]interface{}, map[string]interface{}) error
	CallBackList            func(*sw.StubWrapper, map[string]interface{}, []map[string]interface{}, map[string]interface{}) (map[string]interface{}, error)
}

type IndexDocument struct {
	Design    string
	IndexName string
}

type QuerySearch struct {
	Query       map[string]interface{}
	config      Config
	QueryParser string
}

type FieldSearch struct {
	KeyName string
	Value   interface{}
}

func NewQuery(cfg Config) *QuerySearch {
	if len(cfg.AssetName) == 0 {
		return nil
	}

	// ** Set default tags to remove
	cfg.removeDefaultTags = *GetDefaultRemoveTags()
	return &QuerySearch{
		Query: map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType": cfg.AssetName,
			},
		},
		config: cfg,
	}
}

func (q *QuerySearch) AddCallBack(f func(*sw.StubWrapper, map[string]interface{}, map[string]interface{}) error) {
	q.config.CallBack = f
}

func (q *QuerySearch) AddCallBackList(f func(*sw.StubWrapper, map[string]interface{}, []map[string]interface{}, map[string]interface{}) (map[string]interface{}, error)) {
	q.config.CallBackList = f
}

func (q *QuerySearch) SetPagination(bookMark string, pageSize int32) {
	q.config.BookMark = bookMark
	q.config.PageSize = pageSize
}

func (q *QuerySearch) SetIndexDoc(designDoc, nameIndex string) {
	q.config.IndexDoc.Design = designDoc
	q.config.IndexDoc.IndexName = nameIndex
}

func (q *QuerySearch) SetIndexDocSimple(nameIndex string) {
	q.config.IndexDoc.Design = nameIndex
	q.config.IndexDoc.IndexName = nameIndex
}

func (q *QuerySearch) SetSort(sortFilter []map[string]string) {
	q.config.Sort = sortFilter
}

func (q *QuerySearch) SetResolve(res []string) {
	q.config.Resolve = res
}

func (q *QuerySearch) NoContains(value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$ne": value,
	}
}

func (q *QuerySearch) NoContainsValues(values []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$nin": values,
	}
}

func (q *QuerySearch) AddFieldsOR(fields map[string]interface{}) {
	aux := q.Query["selector"].(map[string]interface{})
	listOfElementsFinds := make([]map[string]interface{}, 0)

	for k, v := range fields {
		listOfElementsFinds = append(listOfElementsFinds, map[string]interface{}{
			k: v,
		})
	}

	aux["$or"] = listOfElementsFinds
	q.Query["selector"] = aux
}

func (q *QuerySearch) AddFieldKeyValue(key string, value interface{}) {
	q.AddField(FieldSearch{
		KeyName: key,
		Value:   value,
	})
}

func (q *QuerySearch) AddField(fields ...FieldSearch) {
	for _, f := range fields {
		if len(f.KeyName) > 0 {
			aux := q.Query["selector"].(map[string]interface{})
			aux[f.KeyName] = f.Value
			q.Query["selector"] = aux
		}
	}
}

func (q *QuerySearch) AddDateRange(nameField string, dataStart, dataEnd time.Time) {
	q.AddField(FieldSearch{
		KeyName: nameField,
		Value:   PeriodDay(dataStart, dataEnd),
	})
}

func (q *QuerySearch) Parser() (string, error) {
	qAux := q.Query

	//** Add Fields for sort query
	if len(q.config.Sort) > 0 {
		qAux["sort"] = q.config.Sort
	}

	//** Add indexDoc and indexName for use in query
	if len(q.config.IndexDoc.Design) > 0 && len(q.config.IndexDoc.IndexName) > 0 {
		qAux["use_index"] = []string{
			"_design/" + q.config.IndexDoc.Design,
			q.config.IndexDoc.IndexName,
		}
	}

	queryJson, err := json.Marshal(qAux)
	if err != nil {
		return "", errors.WrapError(err, "error marshaling query")
	}

	return string(queryJson), nil
}
