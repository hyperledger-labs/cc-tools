# Feature Specification: Mock Rich Query Results

**Feature Branch**: `[001-mock-query-results]`  
**Created**: 2026-02-07  
**Status**: Draft  
**Input**: User description: "Extend MockStub to support GetQueryResult with rich-query mocking: allow registering mock results per query string; otherwise return deterministic results derived from prior mocked puts."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Rich queries work in unit tests (Priority: P1)

As a chaincode test author, I want rich query calls made during transactions (e.g., “search”) to work against the mock ledger so I can validate query-driven behavior without an external database.

**Why this priority**: Query-driven transactions are currently blocked in unit tests because the mock returns “not implemented”.

**Independent Test**: Can be fully tested by writing a few JSON documents into the mock state and running a query that should return a subset of them.

**Acceptance Scenarios**:

1. **Given** the mock state contains multiple JSON documents across different asset types, **When** a rich query is executed with a selector matching one asset type, **Then** only matching documents are returned.
2. **Given** the mock state contains matching and non-matching JSON documents, **When** the rich query is executed, **Then** the result ordering is deterministic.

---

### User Story 2 - Register deterministic results for a specific query (Priority: P2)

As a chaincode test author, I want to register a fixed result set for a specific query string so I can cover complex or unsupported query operators and keep tests focused.

**Why this priority**: The full CouchDB query language is broad; tests need an escape hatch without requiring a full query engine.

**Independent Test**: Can be fully tested by registering a result set for an exact query string and verifying `GetQueryResult` returns those keys.

**Acceptance Scenarios**:

1. **Given** a query string has a registered mock result set, **When** `GetQueryResult` is called with that exact query string, **Then** the iterator returns the registered keys in a deterministic order.
2. **Given** a registered result contains a key that is not present in mock state, **When** iterating results, **Then** the missing key is ignored (not returned).

---

### User Story 3 - Invalid queries fail clearly (Priority: P3)

As a maintainer, I want invalid query inputs to fail with clear errors so test failures are actionable.

**Why this priority**: Ambiguous “not implemented” errors slow down debugging and can mask real issues.

**Independent Test**: Can be tested by calling `GetQueryResult` with malformed JSON or missing required fields.

**Acceptance Scenarios**:

1. **Given** a malformed query string, **When** `GetQueryResult` is called, **Then** it returns an error indicating the query is invalid.
2. **Given** a well-formed query missing the required selector, **When** `GetQueryResult` is called, **Then** it returns an error indicating the selector is missing.

---

### Edge Cases

- What happens when the mock state contains non-JSON values for some keys? (They MUST be ignored by selector-based matching.)
- What happens when the selector references a field that is missing from a document? (That document MUST NOT match.)
- What happens when an empty result set is expected? (Iterator MUST be empty and must not error.)
- What happens when registered results are provided for a query but the corresponding documents were deleted/never written? (Missing keys MUST be ignored.)
- What happens when the same query is executed repeatedly in the same test? (Results MUST be deterministic across runs.)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The mock ledger MUST support rich query calls by allowing `GetQueryResult(query)` to return an iterator over a simulated result set.
- **FR-002**: Tests MUST be able to register a mock result set for a specific (exact) query string.
- **FR-003**: If a mock result set is registered for a query string, `GetQueryResult` MUST return those results (override behavior).
- **FR-004**: For registered results, the iterator MUST return key/value pairs where values reflect the current mock state at iteration time.
- **FR-005**: If a query string does not have a registered result set, the result set MUST be derived from the mock state content created via mocked put operations.
- **FR-006**: Derived results MUST be computed by evaluating a selector against stored JSON documents; documents that are not valid JSON objects MUST be ignored.
- **FR-007**: Derived results MUST support selector equality matching for one or more fields (e.g., select documents where `@assetType` equals a value).
- **FR-008**: Result ordering MUST be deterministic:
  - Registered results MUST return in the same order they were registered.
  - Derived results MUST return in a stable order.
- **FR-009**: Invalid query inputs MUST return a clear error (at minimum: invalid JSON; missing/invalid selector).
- **FR-010**: The feature MUST be test-covered (unit tests) for both derived-query behavior and registered override behavior, consistent with the project constitution.

### Assumptions

- Matching is performed against the JSON documents stored in the mock state.
- Only a minimal subset of the CouchDB query language is required for derived results; registering mock results is the supported escape hatch for complex queries.

### Key Entities *(include if feature involves data)*

- **Mock World State**: The set of key/value pairs written via mocked put operations.
- **Rich Query**: A query string that includes a selector and is executed via `GetQueryResult`.
- **Selector**: The match criteria used to decide which stored JSON documents are included in the result set.
- **Registered Mock Result Set**: A deterministic list of keys associated with an exact query string for override behavior.
- **Result Set Iterator**: The iterable view of query results as key/value pairs.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A unit test can execute a query-driven flow (e.g., search) without any external database dependency.
- **SC-002**: For a fixed mock state and query, repeated runs return the same ordered set of keys.
- **SC-003**: Registering a mock result set for a query produces the exact intended key sequence.
- **SC-004**: Malformed queries fail with an explicit, human-readable error describing what is wrong.
