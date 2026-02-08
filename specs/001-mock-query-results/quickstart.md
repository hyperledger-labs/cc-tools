# Quickstart: Mock Rich Query Results

## Goal

Run unit tests that depend on `GetQueryResult` without a CouchDB backing store.

## Prereqs

- Go toolchain installed
- Import `github.com/hyperledger-labs/cc-tools/mock` in your test files

## Basic Usage

### 1. Derived Query Mode (Automatic Evaluation)

The mock stub automatically evaluates CouchDB Mango selectors against its in-memory state:

```go
import (
    "encoding/json"
    "testing"
    "github.com/hyperledger-labs/cc-tools/mock"
)

func TestDerivedQuery(t *testing.T) {
    stub := mock.NewMockStub("test", nil)
    
    // Add test data
    data := map[string]interface{}{"@assetType": "person", "name": "Alice", "age": 30}
    jsonData, _ := json.Marshal(data)
    stub.PutState("alice", jsonData)
    
    // Query with selector - automatically evaluated
    query := `{"selector": {"@assetType": "person", "age": {"$gte": 25}}}`
    iter, err := stub.GetQueryResult(query)
    if err != nil {
        t.Fatal(err)
    }
    defer iter.Close()
    
    // Iterate results
    for iter.HasNext() {
        kv, err := iter.Next()
        if err != nil {
            t.Fatal(err)
        }
        t.Logf("Key: %s, Value: %s", kv.Key, string(kv.Value))
    }
}
```

### 2. Registry Override Mode (Deterministic Results)

For deterministic testing or when you need exact control over results:

```go
func TestRegistryOverride(t *testing.T) {
    stub := mock.NewMockStub("test", nil)
    
    // Add test data
    stub.PutState("alice", []byte(`{"name": "Alice"}`))
    stub.PutState("bob", []byte(`{"name": "Bob"}`))
    
    // Register exact results for a query
    query := `{"selector": {"@assetType": "person"}}`
    stub.RegisterQueryResult(query, []string{"alice"}) // Only return alice
    
    iter, err := stub.GetQueryResult(query)
    // Will return only alice, ignoring bob
}
```

### 3. Pagination Support

```go
func TestPagination(t *testing.T) {
    stub := mock.NewMockStub("test", nil)
    
    // Add multiple records
    for i := 0; i < 10; i++ {
        key := fmt.Sprintf("key%d", i)
        data := map[string]interface{}{"@assetType": "item", "id": i}
        jsonData, _ := json.Marshal(data)
        stub.PutState(key, jsonData)
    }
    
    // First page: pageSize=5, bookmark=""
    query := `{"selector": {"@assetType": "item"}}`
    iter, meta, err := stub.GetQueryResultWithPagination(query, 5, "")
    if err != nil {
        t.Fatal(err)
    }
    defer iter.Close()
    
    // Process first page
    count := 0
    for iter.HasNext() {
        _, err := iter.Next()
        if err != nil {
            t.Fatal(err)
        }
        count++
    }
    
    // Second page: use bookmark from first page
    iter2, meta2, err := stub.GetQueryResultWithPagination(query, 5, meta.Bookmark)
    // ... process remaining results
}
```

## Supported Query Features

### Operators

- **Combination**: `$and`, `$or`, `$not`, `$nor`
- **Comparison**: `$eq`, `$ne`, `$lt`, `$lte`, `$gt`, `$gte`, `$in`, `$nin`
- **Existence**: `$exists`, `$type`
- **Array**: `$all`, `$elemMatch`, `$allMatch`, `$keyMapMatch`, `$size`
- **String**: `$regex`, `$beginsWith`
- **Arithmetic**: `$mod`

### Query Options

- `selector`: Required - the query selector object
- `sort`: Optional - array of sort specifications `[{"field": "asc|desc"}]`
- `limit`: Optional - maximum results (default 25)
- `skip`: Optional - skip first N results
- `fields`: Optional - array of fields to project `["field1", "field2"]`

### Example Complex Query

```go
query := `{
  "selector": {
    "$and": [
      {"@assetType": "person"},
      {"age": {"$gte": 18, "$lt": 65}},
      {"city": {"$in": ["NYC", "LA", "SF"]}}
    ]
  },
  "sort": [{"age": "asc"}, {"name": "asc"}],
  "limit": 10,
  "skip": 5,
  "fields": ["name", "age"]
}`
```

## Run Tests

- Run the full suite:
  ```bash
  go test ./...
  ```

- Run just the rich query integration test:
  ```bash
  go test -v ./test -run "TestRichQueryIntegration"
  ```

## What to Look For

Integration tests validate:
- ✓ Basic selector evaluation against mock state
- ✓ Comparison operators ($gte, $lt, etc.)
- ✓ Registry override for deterministic results
- ✓ Pagination with bookmarks
- ✓ Error handling for invalid JSON and unsupported features

## Notes

- This feature is confined to the mock implementation used in tests
- No non-standard dependencies introduced (only stdlib and existing Fabric deps)
- Thread-safe: Multiple goroutines can safely register/query concurrently
- `$text` operator not supported - use registry override instead
- Dot notation supported for nested fields: `"user.address.city"`
