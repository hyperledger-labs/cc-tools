// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mock

import "fmt"

// evaluateExists evaluates $exists operator
func evaluateExists(exists bool, expectedValue interface{}) (bool, error) {
	shouldExist, ok := expectedValue.(bool)
	if !ok {
		return false, fmt.Errorf("$exists requires a boolean value")
	}

	return exists == shouldExist, nil
}

// evaluateType evaluates $type operator
func evaluateType(actualValue interface{}, exists bool, expectedValue interface{}) (bool, error) {
	if !exists {
		return false, nil
	}

	expectedType, ok := expectedValue.(string)
	if !ok {
		return false, fmt.Errorf("$type requires a string value")
	}

	actualType := getValueType(actualValue)
	return actualType == expectedType, nil
}

// getValueType returns the CouchDB type name for a value
func getValueType(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch v.(type) {
	case bool:
		return "boolean"
	case float64, int, int32, int64, uint, uint32, uint64:
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "object"
	}
}
