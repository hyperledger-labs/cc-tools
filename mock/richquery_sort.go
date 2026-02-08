package mock

import (
	"fmt"
	"sort"
	"strings"
)

// sortField represents a single sort field with direction
type sortField struct {
	field     string
	ascending bool
}

// parseSortSpec parses the sort specification from the request
func parseSortSpec(sortSpec []interface{}) ([]sortField, error) {
	if len(sortSpec) == 0 {
		return nil, nil
	}

	fields := make([]sortField, 0, len(sortSpec))
	seenDirection := ""

	for _, spec := range sortSpec {
		switch s := spec.(type) {
		case string:
			// Simple field name implies ascending
			fields = append(fields, sortField{field: s, ascending: true})
			if seenDirection == "" {
				seenDirection = "asc"
			} else if seenDirection != "asc" {
				return nil, fmt.Errorf("mixed sort directions are not supported")
			}

		case map[string]interface{}:
			// Field with explicit direction: {"field": "asc"} or {"field": "desc"}
			if len(s) != 1 {
				return nil, fmt.Errorf("sort field object must have exactly one key")
			}

			for field, dirVal := range s {
				direction, ok := dirVal.(string)
				if !ok {
					return nil, fmt.Errorf("sort direction must be a string (asc or desc)")
				}

				// Handle typed field names like "field:string"
				fieldName := field
				if idx := strings.Index(field, ":"); idx > 0 {
					fieldName = field[:idx]
				}

				ascending := true
				switch direction {
				case "asc", "ascending":
					ascending = true
				case "desc", "descending":
					ascending = false
				default:
					return nil, fmt.Errorf("invalid sort direction: %s (must be asc or desc)", direction)
				}

				if seenDirection == "" {
					if ascending {
						seenDirection = "asc"
					} else {
						seenDirection = "desc"
					}
				} else {
					currentDirection := "asc"
					if !ascending {
						currentDirection = "desc"
					}
					if seenDirection != currentDirection {
						return nil, fmt.Errorf("mixed sort directions are not supported")
					}
				}

				fields = append(fields, sortField{field: fieldName, ascending: ascending})
			}

		default:
			return nil, fmt.Errorf("invalid sort specification: must be string or object")
		}
	}

	return fields, nil
}

// sortDocuments sorts documents by the given sort fields
func sortDocuments(docs []map[string]interface{}, sortFields []sortField) {
	if len(sortFields) == 0 {
		return
	}

	sort.SliceStable(docs, func(i, j int) bool {
		for _, sf := range sortFields {
			valI, existsI := getFieldValue(docs[i], sf.field)
			valJ, existsJ := getFieldValue(docs[j], sf.field)

			// Handle missing values (treat as null, which sorts first)
			if !existsI && !existsJ {
				continue
			}
			if !existsI {
				return sf.ascending // null < any value
			}
			if !existsJ {
				return !sf.ascending // value > null
			}

			cmp := compareValues(valI, valJ)
			if cmp != 0 {
				if sf.ascending {
					return cmp < 0
				}
				return cmp > 0
			}
		}
		return false
	})
}
