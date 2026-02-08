package mock

import (
	"encoding/json"
	"fmt"
)

// FindRequest represents a CouchDB _find request shape
type FindRequest struct {
	Selector map[string]interface{} `json:"selector"`
	Limit    int                    `json:"limit,omitempty"`
	Skip     int                    `json:"skip,omitempty"`
	Sort     []interface{}          `json:"sort,omitempty"`
	Fields   []string               `json:"fields,omitempty"`
}

// ParseFindRequest parses and validates a CouchDB _find query string
func ParseFindRequest(query string) (*FindRequest, error) {
	if query == "" {
		return nil, fmt.Errorf("query string cannot be empty")
	}

	var req FindRequest
	if err := json.Unmarshal([]byte(query), &req); err != nil {
		return nil, fmt.Errorf("invalid query JSON: %w", err)
	}

	// Validate required fields
	if req.Selector == nil {
		return nil, fmt.Errorf("selector is required")
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 25 // CouchDB default
	}

	return &req, nil
}
