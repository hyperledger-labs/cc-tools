// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"fmt"
	"regexp"
	"strings"
)

// evaluateMod evaluates $mod operator (value % divisor == remainder)
func evaluateMod(actualValue interface{}, exists bool, expectedValue interface{}) (bool, error) {
	if !exists {
		return false, nil
	}

	// expectedValue should be [divisor, remainder]
	modArray, ok := expectedValue.([]interface{})
	if !ok || len(modArray) != 2 {
		return false, fmt.Errorf("$mod requires an array of two integers [divisor, remainder]")
	}

	divisor := toFloat64(modArray[0])
	remainder := toFloat64(modArray[1])

	if divisor == 0 {
		return false, fmt.Errorf("$mod divisor cannot be zero")
	}

	actualNum := toFloat64(actualValue)
	actualRemainder := int(actualNum) % int(divisor)

	return float64(actualRemainder) == remainder, nil
}

// evaluateRegex evaluates $regex operator
func evaluateRegex(actualValue interface{}, exists bool, expectedValue interface{}) (bool, error) {
	if !exists {
		return false, nil
	}

	pattern, ok := expectedValue.(string)
	if !ok {
		return false, fmt.Errorf("$regex requires a string pattern")
	}

	actualStr, ok := actualValue.(string)
	if !ok {
		return false, nil
	}

	matched, err := regexp.MatchString(pattern, actualStr)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return matched, nil
}

// evaluateBeginsWith evaluates $beginsWith operator
func evaluateBeginsWith(actualValue interface{}, exists bool, expectedValue interface{}) bool {
	if !exists {
		return false
	}

	prefix, ok := expectedValue.(string)
	if !ok {
		return false
	}

	actualStr, ok := actualValue.(string)
	if !ok {
		return false
	}

	return strings.HasPrefix(actualStr, prefix)
}
