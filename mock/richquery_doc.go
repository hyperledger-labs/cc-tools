package mock

import (
	"encoding/json"
	"fmt"
	"strings"
)

// decodeDocument attempts to decode a byte slice as a JSON document
func decodeDocument(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty document")
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// getFieldValue retrieves a field value from a document using dot notation
// e.g., "owner.id" retrieves doc["owner"]["id"]
func getFieldValue(doc map[string]interface{}, fieldPath string) (interface{}, bool) {
	if fieldPath == "" {
		return nil, false
	}

	// Handle dotted field paths
	parts := strings.Split(fieldPath, ".")
	current := interface{}(doc)

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			val, exists := v[part]
			if !exists {
				return nil, false
			}
			current = val
		default:
			return nil, false
		}
	}

	return current, true
}

// setFieldValue sets a field value in a document using dot notation (for projection)
func setFieldValue(doc map[string]interface{}, fieldPath string, value interface{}) {
	if fieldPath == "" {
		return
	}

	parts := strings.Split(fieldPath, ".")
	current := doc

	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}
		if nextMap, ok := current[part].(map[string]interface{}); ok {
			current = nextMap
		} else {
			return
		}
	}

	current[parts[len(parts)-1]] = value
}
