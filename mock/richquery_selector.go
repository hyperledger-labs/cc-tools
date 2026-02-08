package mock

import (
	"fmt"
	"strings"
)

// evaluateSelector evaluates a selector against a document
// Returns true if the document matches the selector
func evaluateSelector(doc map[string]interface{}, selector map[string]interface{}) (bool, error) {
	if selector == nil {
		return true, nil
	}

	// Empty selector matches all documents
	if len(selector) == 0 {
		return true, nil
	}

	// Multiple fields in selector imply AND
	for field, condition := range selector {
		// Check if this is an operator or a field
		if strings.HasPrefix(field, "$") {
			// This is a combination operator
			matched, err := evaluateCombinationOperator(doc, field, condition)
			if err != nil {
				return false, err
			}
			if !matched {
				return false, nil
			}
		} else {
			// This is a field selector
			matched, err := evaluateFieldCondition(doc, field, condition)
			if err != nil {
				return false, err
			}
			if !matched {
				return false, nil
			}
		}
	}

	return true, nil
}

// evaluateFieldCondition evaluates a condition for a specific field
func evaluateFieldCondition(doc map[string]interface{}, field string, condition interface{}) (bool, error) {
	// Validate field name
	if field == "" {
		return false, fmt.Errorf("empty field names are invalid")
	}

	// Get the actual value from the document
	actualValue, exists := getFieldValue(doc, field)

	// Check if condition is a simple value (implicit $eq) or an operator map
	switch cond := condition.(type) {
	case map[string]interface{}:
		// Operator-based condition
		return evaluateConditionOperators(actualValue, exists, cond)
	default:
		// Implicit equality
		if !exists {
			return false, nil
		}
		return compareValues(actualValue, condition) == 0, nil
	}
}

// evaluateCombinationOperator evaluates combination operators ($and, $or, $not, $nor)
func evaluateCombinationOperator(doc map[string]interface{}, operator string, condition interface{}) (bool, error) {
	switch operator {
	case "$and":
		return evaluateAnd(doc, condition)
	case "$or":
		return evaluateOr(doc, condition)
	case "$not":
		return evaluateNot(doc, condition)
	case "$nor":
		return evaluateNor(doc, condition)
	default:
		return false, fmt.Errorf("unknown combination operator: %s", operator)
	}
}

// evaluateConditionOperators evaluates condition operators for a field
func evaluateConditionOperators(actualValue interface{}, exists bool, operators map[string]interface{}) (bool, error) {
	for op, expectedValue := range operators {
		matched, err := evaluateConditionOperator(actualValue, exists, op, expectedValue)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

// evaluateConditionOperator dispatches to the appropriate operator evaluator
func evaluateConditionOperator(actualValue interface{}, exists bool, operator string, expectedValue interface{}) (bool, error) {
	switch operator {
	case "$eq":
		return evaluateEq(actualValue, exists, expectedValue), nil
	case "$ne":
		return evaluateNe(actualValue, exists, expectedValue), nil
	case "$lt":
		return evaluateLt(actualValue, exists, expectedValue), nil
	case "$lte":
		return evaluateLte(actualValue, exists, expectedValue), nil
	case "$gt":
		return evaluateGt(actualValue, exists, expectedValue), nil
	case "$gte":
		return evaluateGte(actualValue, exists, expectedValue), nil
	case "$in":
		return evaluateIn(actualValue, exists, expectedValue), nil
	case "$nin":
		return evaluateNin(actualValue, exists, expectedValue), nil
	case "$exists":
		return evaluateExists(exists, expectedValue)
	case "$type":
		return evaluateType(actualValue, exists, expectedValue)
	case "$mod":
		return evaluateMod(actualValue, exists, expectedValue)
	case "$regex":
		return evaluateRegex(actualValue, exists, expectedValue)
	case "$beginsWith":
		return evaluateBeginsWith(actualValue, exists, expectedValue), nil
	case "$size":
		return evaluateSize(actualValue, exists, expectedValue), nil
	case "$all":
		return evaluateAll(actualValue, exists, expectedValue), nil
	case "$elemMatch":
		return evaluateElemMatch(actualValue, exists, expectedValue)
	case "$allMatch":
		return evaluateAllMatch(actualValue, exists, expectedValue)
	case "$keyMapMatch":
		return evaluateKeyMapMatch(actualValue, exists, expectedValue)
	case "$text":
		return false, fmt.Errorf("$text operator is not supported in derived mode; use query result registration for text searches")
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}
