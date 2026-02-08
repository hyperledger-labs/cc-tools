package mock

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// collationType represents the CouchDB collation type order
type collationType int

const (
	typeNull collationType = iota
	typeFalse
	typeTrue
	typeNumber
	typeString
	typeArray
	typeObject
)

// getCollationType returns the collation type for a value
func getCollationType(v interface{}) collationType {
	if v == nil {
		return typeNull
	}

	switch val := v.(type) {
	case bool:
		if val {
			return typeTrue
		}
		return typeFalse
	case float64, int, int32, int64, uint, uint32, uint64:
		return typeNumber
	case string:
		return typeString
	case []interface{}:
		return typeArray
	case map[string]interface{}:
		return typeObject
	default:
		// Use reflection for other types
		kind := reflect.TypeOf(v).Kind()
		switch kind {
		case reflect.Slice, reflect.Array:
			return typeArray
		case reflect.Map, reflect.Struct:
			return typeObject
		default:
			return typeString
		}
	}
}

// compareValues compares two values using CouchDB collation rules
// Returns: -1 if a < b, 0 if a == b, 1 if a > b
func compareValues(a, b interface{}) int {
	typeA := getCollationType(a)
	typeB := getCollationType(b)

	// Different types: compare by type order
	if typeA != typeB {
		if typeA < typeB {
			return -1
		}
		return 1
	}

	// Same type: compare by value
	switch typeA {
	case typeNull:
		return 0
	case typeFalse:
		return 0
	case typeTrue:
		return 0
	case typeNumber:
		return compareNumbers(a, b)
	case typeString:
		return strings.Compare(fmt.Sprint(a), fmt.Sprint(b))
	case typeArray:
		return compareArrays(a, b)
	case typeObject:
		return compareObjects(a, b)
	}

	return 0
}

// compareNumbers compares two numeric values
func compareNumbers(a, b interface{}) int {
	aVal := toFloat64(a)
	bVal := toFloat64(b)

	if aVal < bVal {
		return -1
	} else if aVal > bVal {
		return 1
	}
	return 0
}

// toFloat64 converts various numeric types to float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	default:
		return 0
	}
}

// compareArrays compares two arrays element by element
func compareArrays(a, b interface{}) int {
	aArr, aOk := a.([]interface{})
	bArr, bOk := b.([]interface{})

	if !aOk || !bOk {
		return 0
	}

	minLen := len(aArr)
	if len(bArr) < minLen {
		minLen = len(bArr)
	}

	for i := 0; i < minLen; i++ {
		cmp := compareValues(aArr[i], bArr[i])
		if cmp != 0 {
			return cmp
		}
	}

	// If all compared elements are equal, shorter array comes first
	if len(aArr) < len(bArr) {
		return -1
	} else if len(aArr) > len(bArr) {
		return 1
	}
	return 0
}

// compareObjects compares two objects by sorted key order
func compareObjects(a, b interface{}) int {
	aObj, aOk := a.(map[string]interface{})
	bObj, bOk := b.(map[string]interface{})

	if !aOk || !bOk {
		return 0
	}

	// Get sorted keys
	aKeys := make([]string, 0, len(aObj))
	for k := range aObj {
		aKeys = append(aKeys, k)
	}
	sort.Strings(aKeys)

	bKeys := make([]string, 0, len(bObj))
	for k := range bObj {
		bKeys = append(bKeys, k)
	}
	sort.Strings(bKeys)

	// Compare keys first
	minLen := len(aKeys)
	if len(bKeys) < minLen {
		minLen = len(bKeys)
	}

	for i := 0; i < minLen; i++ {
		cmp := strings.Compare(aKeys[i], bKeys[i])
		if cmp != 0 {
			return cmp
		}

		// Keys are equal, compare values
		cmp = compareValues(aObj[aKeys[i]], bObj[bKeys[i]])
		if cmp != 0 {
			return cmp
		}
	}

	// If all compared keys/values are equal, object with fewer keys comes first
	if len(aKeys) < len(bKeys) {
		return -1
	} else if len(aKeys) > len(bKeys) {
		return 1
	}
	return 0
}
