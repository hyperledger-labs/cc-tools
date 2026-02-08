// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mock

// evaluateEq evaluates equality (explicit or implicit)
func evaluateEq(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return false
	}
	return compareValues(actualValue, expectedValue) == 0
}

// evaluateNe evaluates not equal
func evaluateNe(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return true
	}
	return compareValues(actualValue, expectedValue) != 0
}

// evaluateLt evaluates less than
func evaluateLt(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return false
	}
	return compareValues(actualValue, expectedValue) < 0
}

// evaluateLte evaluates less than or equal
func evaluateLte(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return false
	}
	return compareValues(actualValue, expectedValue) <= 0
}

// evaluateGt evaluates greater than
func evaluateGt(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return false
	}
	return compareValues(actualValue, expectedValue) > 0
}

// evaluateGte evaluates greater than or equal
func evaluateGte(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return false
	}
	return compareValues(actualValue, expectedValue) >= 0
}

// evaluateIn evaluates $in operator (value in array)
func evaluateIn(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return false
	}

	expectedArray, ok := expectedValue.([]interface{})
	if !ok {
		return false
	}

	for _, item := range expectedArray {
		if compareValues(actualValue, item) == 0 {
			return true
		}
	}

	return false
}

// evaluateNin evaluates $nin operator (value not in array)
func evaluateNin(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return true
	}

	return !evaluateIn(actualValue, exists, expectedValue)
}
