# Quickstart: Mock Rich Query Results

## Goal

Run unit tests that depend on `GetQueryResult` without a CouchDB backing store.

## Prereqs

- Go toolchain installed
- From repo root

## Run tests

- Run the full suite:
  - `go test ./...`

## What to look for

- New tests should validate:
  - Derived results from mock state based on a minimal selector.
  - Registered override results for exact query strings.
  - Deterministic ordering.
  - Clear errors for invalid query JSON and missing selector.

## Notes

- This feature is confined to the mock implementation used in tests.
- No new non-standard dependencies should be introduced.
