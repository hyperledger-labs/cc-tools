// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"container/list"
)

// executeDerivedQuery executes a rich query by deriving results from mock state
// Returns the list of matching keys in the correct order
func executeDerivedQuery(stub *MockStub, req *FindRequest) ([]string, error) {
	// Step 1: Iterate over all keys and match documents against selector
	matchedDocs := make([]map[string]interface{}, 0)
	matchedKeys := make([]string, 0)

	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		key := elem.Value.(string)
		value, err := stub.GetState(key)
		if err != nil {
			continue
		}

		// Skip empty values
		if len(value) == 0 {
			continue
		}

		// Decode as JSON document
		doc, err := decodeDocument(value)
		if err != nil {
			// Not a JSON object, skip
			continue
		}

		// Evaluate selector
		matched, err := evaluateSelector(doc, req.Selector)
		if err != nil {
			return nil, err
		}

		if matched {
			matchedDocs = append(matchedDocs, doc)
			matchedKeys = append(matchedKeys, key)
		}
	}

	// Step 2: Apply sorting if specified
	if len(req.Sort) > 0 {
		sortFields, err := parseSortSpec(req.Sort)
		if err != nil {
			return nil, err
		}

		// Sort documents and reorder keys accordingly
		if len(sortFields) > 0 {
			sortDocumentsWithKeys(matchedDocs, matchedKeys, sortFields)
		}
	}

	// Step 3: Apply skip and limit
	resultKeys := applySkipLimit(matchedKeys, req.Skip, req.Limit)

	return resultKeys, nil
}

// docKeyPair is a helper struct for sorting documents with their keys
type docKeyPair struct {
	doc map[string]interface{}
	key string
}

// sortDocumentsWithKeys sorts documents and their corresponding keys together
func sortDocumentsWithKeys(docs []map[string]interface{}, keys []string, sortFields []sortField) {
	// Create pairs of doc and key
	pairs := make([]docKeyPair, len(docs))
	for i := range docs {
		pairs[i] = docKeyPair{doc: docs[i], key: keys[i]}
	}

	// Sort pairs by documents
	sortDocumentPairs(pairs, sortFields)

	// Unpack sorted pairs
	for i, pair := range pairs {
		docs[i] = pair.doc
		keys[i] = pair.key
	}
}

// sortDocumentPairs sorts document-key pairs
func sortDocumentPairs(pairs []docKeyPair, sortFields []sortField) {
	// Extract just the docs for sorting
	docs := make([]map[string]interface{}, len(pairs))
	for i, pair := range pairs {
		docs[i] = pair.doc
	}

	// Create index map before sorting
	indices := make([]int, len(pairs))
	for i := range indices {
		indices[i] = i
	}

	// Sort indices based on documents
	sortDocuments(docs, sortFields)

	// Now we need to reorder pairs based on how docs were sorted
	// We'll use a simple approach: sort pairs using the same comparison
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if shouldSwap(pairs[i].doc, pairs[j].doc, sortFields) {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
}

// shouldSwap determines if two documents should be swapped based on sort fields
func shouldSwap(docI, docJ map[string]interface{}, sortFields []sortField) bool {
	for _, sf := range sortFields {
		valI, existsI := getFieldValue(docI, sf.field)
		valJ, existsJ := getFieldValue(docJ, sf.field)

		if !existsI && !existsJ {
			continue
		}
		if !existsI {
			return sf.ascending
		}
		if !existsJ {
			return !sf.ascending
		}

		cmp := compareValues(valI, valJ)
		if cmp != 0 {
			if sf.ascending {
				return cmp > 0 // Swap if i > j for ascending
			}
			return cmp < 0 // Swap if i < j for descending
		}
	}
	return false
}

// applyProjectionToKeys applies field projection to a list of keys
// Returns a new iterator with projected values
func applyProjectionToKeys(stub *MockStub, keys []string, fields []string) (*StateQueryIterator, error) {
	if len(fields) == 0 {
		// No projection needed
		return NewStateQueryIterator(stub, keys), nil
	}

	// Create a custom iterator that will apply projection on Next()
	// For simplicity, we'll pre-apply projection and create a temporary state
	projectedStub := &MockStub{
		State: make(map[string][]byte),
		Keys:  list.New(),
	}

	for _, key := range keys {
		value, err := stub.GetState(key)
		if err != nil || len(value) == 0 {
			continue
		}

		// Decode document
		doc, err := decodeDocument(value)
		if err != nil {
			continue
		}

		// Apply projection
		projectedBytes, err := applyFieldsProjection(doc, fields)
		if err != nil {
			continue
		}

		projectedStub.State[key] = projectedBytes
	}

	// Return iterator over projected results
	return NewStateQueryIterator(projectedStub, keys), nil
}
