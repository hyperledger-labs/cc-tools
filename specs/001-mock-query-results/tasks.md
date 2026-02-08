---

description: "Task list for Mock Rich Query Results"

---

# Tasks: Mock Rich Query Results

**Input**: Design documents in specs/001-mock-query-results/ (plan.md, spec.md, research.md, data-model.md, contracts/)

**Tests**: REQUIRED for behavior changes (per constitution). All user story phases below start with tests.

## Phase 1: Setup (Shared)

- [X] T001 Review existing MockStub and iterator patterns in mock/mockstub.go
- [X] T002 Identify existing rich-query callers and expectations in assets/search.go (including pagination) and stubwrapper/stubWrapper.go
- [X] T003 Capture current baseline by running `go test ./...` (see specs/001-mock-query-results/quickstart.md)

---

## Phase 2: Foundational (Blocking Prerequisites)

- [X] T004 Add `_find` request parsing types + helpers in mock/richquery_request.go
- [X] T005 Add JSON document decode + dotted-field-path access helpers in mock/richquery_doc.go
- [X] T006 Add deterministic collation compare helpers for sorting in mock/richquery_collation.go
- [X] T007 Add selector evaluation dispatcher (operator routing skeleton) in mock/richquery_selector.go
- [X] T008 Add a streaming iterator implementing `shim.StateQueryIteratorInterface` in mock/query_iterator.go
- [X] T009 Add pagination helpers for key-based bookmarks + metadata assembly in mock/richquery_pagination.go

**Checkpoint**: Foundational helpers compile; user story work can begin.

---

## Phase 3: User Story 1 - Rich queries work in unit tests (Priority: P1) 🎯 MVP

**Goal**: `GetQueryResult(query)` works against mock state by evaluating Mango selectors and returning deterministic results.

**Independent Test**: Write several JSON docs via `PutState`, call `GetQueryResult` with a selector, and assert the returned keys/values match expectations deterministically.

### Tests for User Story 1 (write first; must fail before implementation)

- [ ] T010 [P] [US1] Add derived-query tests for basic equality + deterministic key order in test/mock_richquery_derived_basic_test.go
- [ ] T011 [P] [US1] Add derived-query tests for logical operators ($and/$or/$not/$nor) in test/mock_richquery_derived_logic_test.go
- [ ] T012 [P] [US1] Add derived-query tests for comparison/inclusion ($lt/$lte/$gt/$gte/$ne/$in/$nin) in test/mock_richquery_derived_compare_test.go
- [ ] T013 [P] [US1] Add derived-query tests for array/map operators ($all/$size/$elemMatch/$allMatch/$keyMapMatch) in test/mock_richquery_derived_array_test.go
- [ ] T014 [P] [US1] Add derived-query tests for sort + skip/limit + fields projection in test/mock_richquery_derived_sort_test.go
- [ ] T015 [P] [US1] Add derived-query tests for $regex/$mod/$beginsWith and $exists/$type in test/mock_richquery_derived_misc_test.go
- [ ] T016 [P] [US1] Add pagination tests for `GetQueryResultWithPagination` (pageSize + bookmark across pages; includes sort) in test/mock_richquery_pagination_test.go

### Implementation for User Story 1

- [X] T017 [US1] Implement derived-mode `GetQueryResult` entrypoint (parse → evaluate → iterator) in mock/mockstub.go
- [X] T018 [P] [US1] Implement field matching semantics (implicit $eq, multi-field AND, dot notation) in mock/richquery_selector.go
- [X] T019 [P] [US1] Implement combination operators ($and/$or/$not/$nor) in mock/richquery_ops_combination.go
- [X] T020 [P] [US1] Implement condition operators ($lt/$lte/$eq/$ne/$gte/$gt/$in/$nin) in mock/richquery_ops_condition.go
- [X] T021 [P] [US1] Implement existence/type operators ($exists/$type) in mock/richquery_ops_existence.go
- [X] T022 [P] [US1] Implement array/map operators ($all/$elemMatch/$allMatch/$keyMapMatch/$size) in mock/richquery_ops_array.go
- [X] T023 [P] [US1] Implement misc operators ($mod/$regex/$beginsWith) in mock/richquery_ops_misc.go
- [X] T024 [US1] Implement sort parsing + deterministic multi-key sorting in mock/richquery_sort.go
- [X] T025 [US1] Implement skip/limit and fields projection (dot notation) in mock/richquery_projection.go
- [X] T026 [US1] Wire the derived query engine pipeline (match → sort → skip/limit → projection) in mock/richquery_engine.go
- [X] T027 [US1] Implement mock `GetQueryResultWithPagination` using key-based bookmarks + metadata in mock/mockstub.go

**Checkpoint**: US1 tests pass; derived-mode rich queries work without overrides.

---

## Phase 4: User Story 2 - Register deterministic results for a specific query (Priority: P2)

**Goal**: Tests can register exact query-string overrides that take precedence over derived evaluation.

**Independent Test**: Register keys for a query, call `GetQueryResult` with the same query, and assert the iterator returns exactly those keys (skipping missing keys) and reads latest state values at `Next()`.

### Tests for User Story 2 (write first; must fail before implementation)

- [ ] T028 [P] [US2] Add override-registry precedence tests (override beats derived) in test/mock_richquery_override_test.go
- [ ] T029 [P] [US2] Add override-registry iterator semantics tests (missing keys skipped; latest value read at Next) in test/mock_richquery_override_state_test.go
- [ ] T030 [P] [US2] Add override pagination tests for `GetQueryResultWithPagination` (registered order preserved across pages) in test/mock_richquery_override_pagination_test.go

### Implementation for User Story 2

- [X] T031 [US2] Add thread-safe registry fields + `RegisterQueryResult` helper on MockStub in mock/mockstub.go
- [X] T032 [US2] Update `GetQueryResult` to check registry first and snapshot keys for iteration in mock/mockstub.go
- [X] T033 [US2] Enforce deep-copy semantics on registry set/get to prevent caller mutation in mock/mockstub.go
- [X] T034 [US2] Ensure `GetQueryResultWithPagination` checks registry first and paginates registered keys deterministically in mock/mockstub.go

**Checkpoint**: US2 tests pass; overrides are deterministic and isolated.

---

## Phase 5: User Story 3 - Invalid queries fail clearly (Priority: P3)

**Goal**: Malformed queries and unsupported/invalid operators fail with clear, actionable errors.

**Independent Test**: Call `GetQueryResult` with malformed JSON, missing selector, unknown operator, and invalid sort; assert the error message includes the root cause.

### Tests for User Story 3 (write first; must fail before implementation)

- [ ] T035 [P] [US3] Add tests for invalid JSON and missing/invalid selector errors in test/mock_richquery_errors_test.go
- [ ] T036 [P] [US3] Add tests for unknown operator and unsupported $text errors in test/mock_richquery_unsupported_test.go
- [ ] T037 [P] [US3] Add tests for invalid sort shape/rules errors in test/mock_richquery_sort_errors_test.go
- [ ] T038 [P] [US3] Add tests for invalid pagination inputs (unknown bookmark, invalid pageSize) in test/mock_richquery_pagination_errors_test.go

### Implementation for User Story 3

- [X] T039 [US3] Implement request validation errors (invalid JSON, selector required/object) in mock/richquery_request.go
- [X] T040 [US3] Implement unknown operator detection + error plumbing in mock/richquery_selector.go
- [X] T041 [US3] Implement `$text` as unsupported-in-derived-mode with a clear error in mock/richquery_ops_misc.go
- [X] T042 [US3] Implement sort validation errors (single-key objects, no mixed directions, etc.) in mock/richquery_sort.go
- [X] T043 [US3] Implement pagination validation errors (unknown bookmark, invalid pageSize) in mock/richquery_pagination.go

**Checkpoint**: US3 tests pass; failures are clear and consistent.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T044 [P] Update quickstart expectations/examples to include rich-query behavior in specs/001-mock-query-results/quickstart.md
- [X] T045 [P] Finalize registry helper name and document pagination semantics in specs/001-mock-query-results/contracts/mock-rich-query.md
- [X] T046 [P] Sync specs/001-mock-query-results/spec.md scope language with the plan (remove "minimal subset" assumption or explicitly constrain scope)
- [X] T047 Run `gofmt` over modified files in mock/ (touchpoints include mock/mockstub.go)
- [X] T048 Run `go vet ./...` and address any findings in the touched packages (module file: go.mod)
- [X] T049 Run `go test ./...` and ensure new tests pass (see specs/001-mock-query-results/quickstart.md)
- [ ] T050 [P] Add/extend an end-to-end regression test that exercises the assets search pagination path via existing tests in test/tx_search_test.go (or a new focused test if needed)

---

## Dependencies & Execution Order

- Phase 1 (Setup) → Phase 2 (Foundational) → User stories (US1 → US2 → US3) → Phase 6 (Polish)
- US1 is the MVP and blocks US2/US3 only insofar as they build on a working `GetQueryResult`.
- US2 depends on US1 having a working iterator return path.
- US3 depends on US1 parsing/evaluation paths existing so errors can be validated.

## Parallel Opportunities

- US1 tests can be authored in parallel: T009–T014 (separate files under test/).
- US1 operator implementations can be authored in parallel: T017–T021 (separate files under mock/).
- US2 tests (T025–T026) can be authored in parallel.
- US3 tests (T030–T032) can be authored in parallel.

## Parallel Example: User Story 1

- Run in parallel (tests): T009 + T010 + T011 + T012 + T013 + T014 (all under test/)
- Run in parallel (implementation): T017 + T018 + T019 + T020 + T021 (all under mock/)

## Implementation Strategy

- Deliver MVP first: complete Phase 1–3 (US1) and validate derived-mode rich queries.
- Add override registry next (US2) to unblock complex-query tests.
- Add robust error surfacing last (US3), ensuring diagnostics are actionable.
