<!--
Sync Impact Report

- Version change: N/A (template) → 0.1.0
- Modified principles: N/A (template placeholders replaced)
- Added sections: Core Principles (filled), Quality Gates, Development Workflow, Governance (filled)
- Removed sections: N/A
- Templates requiring updates:
	- ✅ updated: .specify/templates/plan-template.md
	- ✅ updated: .specify/templates/tasks-template.md
	- ✅ no change: .specify/templates/spec-template.md
	- ✅ no change: .specify/templates/checklist-template.md
	- ✅ no change: .specify/templates/agent-file-template.md
- Follow-up TODOs:
	- TODO(RATIFICATION_DATE): original ratification date is unknown
-->

# Hyperledger Labs CC Tools Constitution

## Core Principles

### I. Clarity Over Cleverness
Code MUST be easy to read, review, and maintain.

- Prefer simple, idiomatic Go over clever abstractions.
- Choose descriptive names (types, funcs, vars) that convey intent.
- Keep functions small and focused; avoid hidden side effects.
- Public-facing behavior MUST be discoverable via docs and tests.

Rationale: cc-tools is a shared foundation; readability reduces defects and
lowers the cost of contribution.

### II. Consistent Go Quality Gates
Every change MUST keep the codebase in a “green” state under standard Go tools.

- Code MUST be formatted with `gofmt`.
- `go test ./...` MUST pass.
- `go vet ./...` MUST pass (or any exceptions MUST be justified in the PR).
- New exported identifiers MUST include GoDoc comments.

Rationale: these gates are low-friction, widely understood, and prevent
avoidable regressions.

### III. Testing Is Not Optional (NON-NEGOTIABLE)
Any change that alters behavior MUST be covered by tests.

- Bug fixes MUST include a test that fails before the fix and passes after.
- New features MUST include unit tests for the core logic.
- When behavior depends on Fabric/stub interactions, prefer using the existing
	mocks and add focused tests around the boundary.
- If tests are intentionally omitted (e.g., pure comment typo), the PR MUST
	explicitly state why.

Rationale: chaincode changes are high-risk; tests are the cheapest safety net.

### IV. Public APIs Are Contracts
Changes to exported packages/types/functions MUST be treated as contract
changes.

- Breaking changes MUST be avoided unless there is a clear migration path.
- Any public API change MUST update documentation and tests that assert the
	contract.
- Error messages and returned structures that callers may depend on MUST remain
	stable or be versioned.

Rationale: consumers build chaincode on top of cc-tools; stability is part of
the product.

### V. Small, Reviewable Changes
Prefer incremental delivery that is easy to review.

- PRs MUST stay focused; unrelated refactors belong in separate PRs.
- Refactors MUST preserve behavior unless explicitly stated and tested.
- Complexity MUST be justified when simpler options exist.

Rationale: smaller diffs reduce review time and regression probability.

## Quality Gates

- All CI checks MUST pass.
- Code review MUST verify conformance to Core Principles.
- New logic MUST have tests at an appropriate level (unit by default).
- Changes touching ledger/state interactions MUST include edge-case coverage
	(missing keys, bad input, empty results, permission checks).

## Development Workflow

- Start from a clear problem statement (issue/PR description) and define
	acceptance criteria.
- Keep API changes explicit: document, test, and call out migration notes.
- Prefer improving existing helpers/utilities over duplicating patterns.
- If a change increases complexity, include a short justification and
	alternatives considered.

## Governance
<!-- Example: Constitution supersedes all other practices; Amendments require documentation, approval, migration plan -->

- This constitution supersedes other contributing norms.
- Every PR review MUST include an explicit constitution compliance check.
- Amendments MUST be made via PR that:
	- states the rationale,
	- updates dependent templates under `.specify/templates/`, and
	- bumps the constitution version following semantic versioning:
		- MAJOR: incompatible governance/principle removals or redefinitions
		- MINOR: new principle/section or materially expanded guidance
		- PATCH: clarifications/typos without semantic change
- If a PR cannot comply with a MUST requirement, it MUST be labeled as a
	constitution exception and include a migration/mitigation plan.

**Version**: 0.1.0 | **Ratified**: TODO(RATIFICATION_DATE): original ratification date is unknown | **Last Amended**: 2026-02-07
