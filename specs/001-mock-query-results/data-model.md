# Data Model: Mock Rich Query Results

## Entities

### 1) `MockStub` (existing)

New fields (additive, no breaking changes):

- `queryResultsMu` (mutex)
  - Type: `sync.RWMutex`
  - Purpose: guard query result registry.
- `queryResults`
  - Type: `map[string][]string`
  - Key: exact query string
  - Value: ordered list of state keys to return for that query

Validation rules:

- Registry values must be copied on set/get to avoid caller mutation.

### 2) `StateQueryIterator` (new helper type)

Fields:

- `stub` reference (the `MockStub`)
- `keys` (snapshot of keys to iterate; derived or registered)
- `pos` (current index)
- `closed` (bool)

Behavior:

- `HasNext()` returns `false` when closed or `pos >= len(keys)`.
- `Next()` returns a `queryresult.KV` for the next available key.
  - If a key has no current value in state, iterator skips it and continues.
- `Close()` marks closed and releases references.

## Relationships

- `MockStub.GetQueryResult(query)` returns a `StateQueryIterator`.
- Derived key lists are computed from `MockStub.Keys` and `MockStub.State`.

## State transitions

- Register override: query string → key list is updated (overwritten).
- Query execution: if override exists → use override; else derive.
