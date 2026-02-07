# Contracts: Mock Rich Query Results (Internal)

This feature does not add public chaincode transactions or REST endpoints.
It defines an internal contract for test behavior of the mock stub.

## Operation: `GetQueryResult(query)`

### Inputs

- `query`: string (expected to be JSON)

Minimum supported JSON shape:

```json
{
  "selector": {
    "fieldA": "valueA",
    "fieldB": "valueB"
  },
  "fields": ["fieldA", "fieldB"],
  "skip": 0,
  "limit": 25,
  "sort": [
    {"fieldA": "asc"},
    {"fieldB": "desc"}
  ]
}
```

Selector semantics follow CouchDB Mango (`_find`) behavior, including logical
operators like `$and`/`$or` and the standard operator set from the CouchDB docs
referenced by `transactions/search.go`.

Note: `$text` requires CouchDB text/nouveau indexes; the mock derived evaluator
should fail clearly for `$text` selectors and tests should use the override
registry for those queries.

### Outputs

- Success: returns an iterator implementing `shim.StateQueryIteratorInterface`.
  - Each `Next()` yields `queryresult.KV{Key, Value}`
  - Ordering is deterministic:
    - Override mode: registered order.
    - Derived mode: key order from `MockStub.Keys`.

### Errors

- Invalid JSON: return error describing invalid query JSON.
- Missing `selector` or non-object selector: return error indicating selector is required.
- Unknown/unsupported selector operator: return error indicating the operator is not supported.

### Override registry (test-only helper)

- Operation: `RegisterQueryResult(query string, keys []string)` (exact name TBD by implementation)
- Behavior: exact match on `query` string.
- Rules:
  - Stores a copy of `keys`.
  - Retrieval returns a copy of the stored keys.
