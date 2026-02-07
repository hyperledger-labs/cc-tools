# Research: Mock Rich Query Results

## Decision 1: Thread-safe query registry

- Decision: Store per-query mock results in `MockStub` using `map[string][]string` protected by `sync.RWMutex`.
- Rationale: Concurrency safety for parallel test execution; minimal change surface; stdlib-only.
- Alternatives considered:
  - `sync.Map`: would avoid explicit locks but makes ordered slices and deep-copy semantics less explicit.
  - Global registry: rejected; would leak between tests and break isolation.

## Decision 2: Exact query-string matching for overrides

- Decision: Registered results override by exact match on the query string passed into `GetQueryResult`.
- Rationale: Simple and deterministic; avoids needing canonical JSON normalization.
- Alternatives considered:
  - Canonicalize JSON before lookup: adds complexity and edge cases; not required by spec.

## Decision 3: Derived results from existing mock state

- Decision: When no override exists, derive results by evaluating CouchDB/Mango (`_find`) selector semantics (including `$and/$or` and the documented operator set) against JSON documents stored in the mock state, and apply `sort` when present.
- Rationale: Matches the main cc-tools use case (`transactions/search.go`) and keeps tests representative of CouchDB behavior.
- Alternatives considered:
  - Only support a minimal subset: rejected; main use case needs operators + sorting.
  - Only support override mode: rejected; tests should validate real selector behavior.

## Decision 4: Deterministic ordering

- Decision: For derived results, iterate `MockStub.Keys` (already kept lexicographically ordered by `PutState`) and emit matches in that order.
- Rationale: Stable results across runs; mirrors existing mock patterns.
- Alternatives considered:
  - Sort per-query result set: redundant work and risks diverging from mock’s established ordering.

## Decision 5: Iterator semantics

- Decision: Implement a `shim.StateQueryIteratorInterface` that yields one record per `Next()` and reads the latest value from state at iteration time.
- Rationale: Simulates streaming; avoids building large in-memory `[]*queryresult.KV` unless required.
- Alternatives considered:
  - Pre-materialize all results into `[]KV`: simpler but less “streaming”; acceptable but not preferred given spec emphasis.
