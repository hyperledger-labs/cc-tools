// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mock

import "fmt"

// evaluateAll evaluates $all operator (array contains all elements)
func evaluateAll(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return false
	}

	actualArray, ok := actualValue.([]interface{})
	if !ok {
		return false
	}

	expectedArray, ok := expectedValue.([]interface{})
	if !ok {
		return false
	}

	// Check that all expected elements are in the actual array
	for _, expected := range expectedArray {
		found := false
		for _, actual := range actualArray {
			if compareValues(actual, expected) == 0 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// evaluateElemMatch evaluates $elemMatch operator (any array element matches selector)
func evaluateElemMatch(actualValue interface{}, exists bool, expectedValue interface{}) (bool, error) {
	if !exists {
		return false, nil
	}

	actualArray, ok := actualValue.([]interface{})
	if !ok {
		return false, nil
	}

	selector, ok := expectedValue.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("$elemMatch requires an object selector")
	}

	// Check if any element matches the selector
	for _, elem := range actualArray {
		// Element must be an object for selector matching
		elemObj, ok := elem.(map[string]interface{})
		if !ok {
			continue
		}

		matched, err := evaluateSelector(elemObj, selector)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}

// evaluateAllMatch evaluates $allMatch operator (all array elements match selector)
func evaluateAllMatch(actualValue interface{}, exists bool, expectedValue interface{}) (bool, error) {
	if !exists {
		return false, nil
	}

	actualArray, ok := actualValue.([]interface{})
	if !ok {
		return false, nil
	}

	selector, ok := expectedValue.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("$allMatch requires an object selector")
	}

	// Empty array matches
	if len(actualArray) == 0 {
		return true, nil
	}

	// Check if all elements match the selector
	for _, elem := range actualArray {
		elemObj, ok := elem.(map[string]interface{})
		if !ok {
			return false, nil
		}

		matched, err := evaluateSelector(elemObj, selector)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}

	return true, nil
}

// evaluateKeyMapMatch evaluates $keyMapMatch operator (any map key matches selector)
func evaluateKeyMapMatch(actualValue interface{}, exists bool, expectedValue interface{}) (bool, error) {
	if !exists {
		return false, nil
	}

	actualMap, ok := actualValue.(map[string]interface{})
	if !ok {
		return false, nil
	}

	selector, ok := expectedValue.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("$keyMapMatch requires an object selector")
	}

	// Check if any value in the map matches the selector
	for _, value := range actualMap {
		valueObj, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		matched, err := evaluateSelector(valueObj, selector)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}

// evaluateSize evaluates $size operator (array length)
func evaluateSize(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return false
	}

	actualArray, ok := actualValue.([]interface{})
	if !ok {
		return false
	}

	expectedSize, ok := expectedValue.(float64)
	if !ok {
		// Try int
		if intSize, ok := expectedValue.(int); ok {
			expectedSize = float64(intSize)
		} else {
			return false
		}
	}

	return float64(len(actualArray)) == expectedSize
}
