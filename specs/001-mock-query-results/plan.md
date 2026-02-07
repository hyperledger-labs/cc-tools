# Implementation Plan: Mock Rich Query Results

**Branch**: `[001-mock-query-results]` | **Date**: 2026-02-07 | **Spec**: specs/001-mock-query-results/spec.md
**Input**: Feature specification from `specs/001-mock-query-results/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.github/agents/speckit.plan.agent.md` for the execution workflow.

## Summary

Enable rich queries in unit tests by extending the existing mock stub so
`GetQueryResult(query)` returns a deterministic, CouchDB-like result iterator.

Derived results evaluate the CouchDB/Mango selector language against the JSON
documents previously written to the mock world state via mocked put operations,
including logical operators (e.g. `$and`, `$or`) and `sort`.

To support the main cc-tools caller (`assets.Search`), the mock also needs to
support Fabric’s pagination API (`GetQueryResultWithPagination`) in addition to
non-paginated `GetQueryResult`.

Tests can still optionally register explicit mock results for an exact query
string as an escape hatch.

## Technical Context

**Language/Version**: Go 1.21  
**Primary Dependencies**: Hyperledger Fabric chaincode shim (`fabric-chaincode-go`) and protos  
**Storage**: In-memory mock world state (`MockStub.State` + ordered `MockStub.Keys`)  
**Testing**: `go test ./...` (existing `test/` suite)  
**Target Platform**: Hyperledger Fabric chaincode unit tests (no external CouchDB)  
**Project Type**: single (Go library + test suite)  
**Performance Goals**: N/A (unit-test focused)  
**Constraints**:
- Must not break existing mock behaviors (`PutState`, `GetState`, iterators, private data mocks).
- Rich query result ordering must be deterministic.
- Implementation must use a thread-safe map for query result registration.
- Do not add new non-standard libraries (stdlib-only additions are acceptable).
**Scale/Scope**: Full Mango selector evaluation (with documented mock limitations), deterministic sorting/projection/paging behavior, and an exact-match override registry.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `gofmt` run on changed Go files
- `go test ./...` passes (include new/updated tests for behavior changes)
- `go vet ./...` passes (or PR includes explicit justification)
- Public API changes documented + tested (migration notes if needed)
- PR scope is focused; complexity justified if introduced

## Project Structure

### Documentation (this feature)

```text
specs/001-mock-query-results/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
```text
mock/
└── mockstub.go           # Extend rich query support here

stubwrapper/
└── stubWrapper.go        # Already calls Stub.GetQueryResult

assets/
└── search.go             # Uses GetQueryResult (and pagination)

transactions/
└── search.go             # Exposes search transaction

test/
└── (add new tests here to validate rich query behavior)
```

**Structure Decision**: Single Go module; implementation changes are scoped to
`mock/mockstub.go` and new unit tests under `test/` (or `mock/` package tests if
kept close to the mock).

## Design Overview

### 1) Query result registration (override)

Add a query result registry on `MockStub`:

- Thread-safe map: `map[string][]string` guarded by `sync.RWMutex` (or equivalent).
- Registration API: a helper method to register keys for an exact query string.
- Behavior: if a query string is registered, it takes precedence over derived evaluation.
- Ordering: return keys in the order registered.

Rationale: supports complex selectors without building a full CouchDB engine.

### 2) Derived result evaluation (default)

If no explicit mock exists for the query, derive results by:

- Parsing the query JSON as a CouchDB `_find` request.
- Parsing the query JSON as a CouchDB Mango (`_find`) request shape.
  - Required: `selector` (object)
  - Optional: `limit` (number, default 25), `skip` (number, default 0)
  - Optional: `sort` (array), `fields` (array)
  - Optional but ignored by the mock query engine: `use_index`, `allow_fallback`,
    `conflicts`, `r`, `bookmark` (CouchDB request-body field), `update`, `stable`,
    `stale`, `execution_stats`
- Iterating over the existing ordered key list (`MockStub.Keys`) to keep results stable.
- For each key, reading state via existing `GetState` and attempting JSON decode.
- Matching documents by evaluating Mango selector semantics.
  - Documents that are not JSON objects are ignored.
  - Missing selector or malformed query JSON returns a clear error.
  - Unknown/unsupported selector operator returns a clear error (do not silently ignore).

Additional selector validity rules to match CouchDB behavior:

- Empty field names are invalid.
- Field names starting with `$` are treated as operators unless escaped (e.g. `"\\$foo"`).

#### Selector operator coverage (derived mode)

Implement the full Mango selector operator set as documented in CouchDB Mango
Selectors (referenced by `transactions/search.go`).

Combination operators:

- `$and` (array of selectors)
- `$or` (array of selectors)
- `$not` (selector)
- `$nor` (array of selectors)
- `$all` (array; array contains all elements)
- `$elemMatch` (selector; any array element matches)
- `$allMatch` (selector; all array elements match)
- `$keyMapMatch` (selector; any map key matches)
- `$text` (string)

Condition operators:

- `$lt`, `$lte`, `$eq`, `$ne`, `$gte`, `$gt`
- `$exists` (boolean)
- `$type` (string: `null|boolean|number|string|array|object`)
- `$in`, `$nin` (array)
- `$size` (integer)
- `$mod` (two-integer array: `[divisor, remainder]`)
- `$regex` (string pattern; evaluated with Go regexp)
- `$beginsWith` (string prefix)

Notes / mock-specific constraints:

- `$text` in CouchDB requires a search/nouveau index and Lucene syntax; the mock
  evaluator will treat `$text` as unsupported in derived mode and return a clear
  error directing tests to use the override registry for such queries.
- Strict type matching is enforced for condition operators (per Mango).

Field matching rules:

- A plain value selector (e.g. `{"a": 1}`) is treated as `$eq`.
- Multiple fields in the same selector object imply AND.
- Nested fields are supported via dot notation (e.g. `"owner.id"`).
- Missing fields do not match, except where CouchDB semantics dictate otherwise
  (e.g. `$exists: false`).

#### Sorting

If `sort` is provided, sort the matched documents before iteration:

- Support Mango sort forms:
  - `[{"field": "asc"}, {"other": "desc"}]`
  - `["field", "other"]` (defaults to ascending)
- Apply multi-key sort in order, stable across runs.
- Validate sort array rules:
  - Each object element must have a single key (otherwise error).
  - CouchDB does not support mixed sort directions across fields; if mixed,
    return an error.
  - Support typed sort field form like `"field:string"` (parse and ignore the
    explicit type unless it is needed for deterministic ordering).
- Implement deterministic JSON collation for sorting.
  - Type order follows CouchDB collation (documented): `null < false < true < number < string < array < object`.
  - Numbers compared numerically.
  - Strings: mock uses deterministic Unicode codepoint ordering (note: CouchDB
    uses ICU collation for some comparisons; tests depending on ICU-specific
    ordering should use override registration).
  - Arrays compared element-by-element.
  - Objects compared by sorted key order.

If `sort` is omitted, preserve deterministic iteration order from `MockStub.Keys`.

#### Limiting, skipping, and field projection

- `skip`: after sorting (if any), drop the first N matched documents.
- `limit`: after skip, take up to N results; default is 25 (CouchDB default).
- `fields`: if provided, return only those fields in each document (no automatic
  inclusion of other metadata fields). Field names use dotted notation.

Implementation detail: apply `fields` projection on the value bytes yielded by
the iterator (so consumers like `assets.Search` can unmarshal the projected
result).

### 2b) Pagination API support (`GetQueryResultWithPagination`)

The cc-tools `assets.Search` code path uses Fabric pagination when the request
includes `limit`, calling `GetQueryResultWithPagination(query, pageSize, bookmark)`.
In this path, `limit` and `bookmark` are removed from the JSON query body before
marshaling.

Design goals:

- Provide deterministic paging over the same derived/override result set.
- Return `QueryResponseMetadata` with a deterministic `Bookmark` suitable for
  requesting the next page.

Behavior:

- Evaluate the query exactly as in `GetQueryResult` (selector + sort + fields +
  request-body `skip`), producing the full ordered list of matching keys.
- Apply pagination arguments:
  - `bookmark == ""`: start from the beginning.
  - `bookmark != ""`: resume *after* the bookmarked key.
    - If the bookmark key is not found in the result set, return a clear error.
  - `pageSize`: acts as the effective limit for the page.
- Returned metadata:
  - `FetchedRecordsCount`: number of records returned in this page.
  - `Bookmark`: empty if there are no more results; otherwise the last key
    returned in this page.

Notes:

- This mock bookmark format is intentionally deterministic (key-based) and is
  not intended to match CouchDB’s opaque bookmark encoding.
- Override mode should also support pagination using the same bookmark rules.

### 3) Streaming iterator simulation

Implement a `shim.StateQueryIteratorInterface` that:

- Streams results one-at-a-time (like CouchDB pagination/streaming).
- Returns `queryresult.KV` where `Value` is read from the current state at `Next()` time.
- Skips keys that are missing in state (do not return nil-value records).

Compatibility constraints:

- Must not change semantics of existing range iterator.
- Must not change `PutState`/`GetState` behaviors or key ordering.

## Testing Strategy

- Add unit tests for:
  - Derived selector behavior with logical operators (`$and`, `$or`, `$not`, `$nor`).
  - Derived selector behavior with comparison operators (`$gt/$gte/$lt/$lte/$in/$nin/$ne`).
  - Derived selector behavior with existence/type (`$exists`, `$type`).
  - Derived selector behavior with array operators (`$all`, `$size`, `$elemMatch`).
  - Derived selector behavior with `$regex` and `$mod`.
  - Sorting behavior (single and multi-field, asc/desc) and deterministic collation.
  - Pagination behavior via `GetQueryResultWithPagination` (page size + bookmark across pages).
  - Override behavior for registered exact query strings.
  - Invalid query JSON and missing selector error paths.
  - Unknown operator error paths (explicit error message).
  - Deterministic ordering.

All behavior changes must be test-covered per constitution.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

None.
