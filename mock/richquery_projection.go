// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"encoding/json"
)

// applySkipLimit applies skip and limit to a list of keys
func applySkipLimit(keys []string, skip int, limit int) []string {
	// Apply skip
	if skip > 0 {
		if skip >= len(keys) {
			return []string{}
		}
		keys = keys[skip:]
	}

	// Apply limit
	if limit > 0 && limit < len(keys) {
		keys = keys[:limit]
	}

	return keys
}

// applyFieldsProjection applies field projection to a document
// Returns the projected document as JSON bytes
func applyFieldsProjection(doc map[string]interface{}, fields []string) ([]byte, error) {
	if len(fields) == 0 {
		// No projection, return full document
		return json.Marshal(doc)
	}

	// Create new document with only requested fields
	projected := make(map[string]interface{})
	for _, field := range fields {
		if value, exists := getFieldValue(doc, field); exists {
			setFieldValue(projected, field, value)
		}
	}

	return json.Marshal(projected)
}
