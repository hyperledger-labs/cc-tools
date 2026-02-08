package mock

import "fmt"

// evaluateAnd evaluates $and operator
func evaluateAnd(doc map[string]interface{}, condition interface{}) (bool, error) {
	conditions, ok := condition.([]interface{})
	if !ok {
		return false, fmt.Errorf("$and requires an array of conditions")
	}

	for _, cond := range conditions {
		condMap, ok := cond.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("$and condition must be an object")
		}

		matched, err := evaluateSelector(doc, condMap)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}

	return true, nil
}

// evaluateOr evaluates $or operator
func evaluateOr(doc map[string]interface{}, condition interface{}) (bool, error) {
	conditions, ok := condition.([]interface{})
	if !ok {
		return false, fmt.Errorf("$or requires an array of conditions")
	}

	if len(conditions) == 0 {
		return false, nil
	}

	for _, cond := range conditions {
		condMap, ok := cond.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("$or condition must be an object")
		}

		matched, err := evaluateSelector(doc, condMap)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}

// evaluateNot evaluates $not operator
func evaluateNot(doc map[string]interface{}, condition interface{}) (bool, error) {
	condMap, ok := condition.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("$not requires an object")
	}

	matched, err := evaluateSelector(doc, condMap)
	if err != nil {
		return false, err
	}

	return !matched, nil
}

// evaluateNor evaluates $nor operator
func evaluateNor(doc map[string]interface{}, condition interface{}) (bool, error) {
	conditions, ok := condition.([]interface{})
	if !ok {
		return false, fmt.Errorf("$nor requires an array of conditions")
	}

	for _, cond := range conditions {
		condMap, ok := cond.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("$nor condition must be an object")
		}

		matched, err := evaluateSelector(doc, condMap)
		if err != nil {
			return false, err
		}
		if matched {
			return false, nil
		}
	}

	return true, nil
}
