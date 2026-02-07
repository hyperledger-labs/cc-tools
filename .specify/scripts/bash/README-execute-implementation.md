# Execute Implementation Plan - Documentation

## Overview

The `execute-implementation.sh` script automates the validation and preparation workflow for implementing features defined through the Spec-Driven Development (SDD) process. It ensures all prerequisites are met before beginning implementation work.

## Purpose

This script serves as a **pre-implementation validation tool** that:

1. ✅ Verifies all required specification documents exist
2. ✅ Validates checklist completion status
3. ✅ Loads implementation context from spec documents
4. ✅ Verifies project setup (ignore files, etc.)
5. ✅ Parses task structure from tasks.md
6. ⚠️ Provides framework for task execution (manual or AI-assisted)

## Usage

### Basic Usage

```bash
# From repository root
./.specify/scripts/bash/execute-implementation.sh

# Using environment variable for feature selection
SPECIFY_FEATURE="001-feature-name" ./.specify/scripts/bash/execute-implementation.sh
```

### Command Line Options

| Option | Description |
|--------|-------------|
| `--help`, `-h` | Show help message and exit |
| `--dry-run` | Show execution plan without running tasks |
| `--yes` | Skip confirmation prompts (auto-proceed) |
| `-- ARGUMENTS` | Pass user input arguments (after `--`) |

### Examples

```bash
# Interactive mode (default) - prompts if checklists incomplete
./execute-implementation.sh

# Auto-proceed mode - skips prompts
./execute-implementation.sh --yes

# Dry run - validate without executing
./execute-implementation.sh --dry-run

# With user arguments
./execute-implementation.sh -- "Focus on test coverage"
```

## Workflow Steps

### Step 1: Prerequisites Check

Runs `.specify/scripts/bash/check-prerequisites.sh` to verify:
- Feature directory exists (specs/NNN-feature-name/)
- plan.md exists
- tasks.md exists
- Lists available optional documents

**Output:**
```
[INFO] Step 1: Checking prerequisites...
[SUCCESS] Feature directory: /path/to/specs/001-feature
[INFO] Available documents: research.md data-model.md contracts/ quickstart.md tasks.md
```

### Step 2: Checklist Validation

Scans `checklists/` directory and validates completion status:

**Status Table Format:**
```
| Checklist | Total | Completed | Incomplete | Status |
|-----------|-------|-----------|------------|--------|
| requirements.md | 16  | 16  | 0  | ✓ PASS |
| design.md       | 12  | 8   | 4  | ✗ FAIL |
```

**Behavior:**
- ✅ **All complete**: Proceeds automatically
- ⚠️ **Some incomplete**: Prompts user to confirm (unless `--yes`)
  - `yes/y/proceed/continue` → Continue
  - `no/n/wait/stop` → Halt execution

### Step 3: Context Loading

Loads implementation context from spec documents:

**Required:**
- `tasks.md` - Task breakdown
- `plan.md` - Technical plan

**Optional (if exists):**
- `data-model.md` - Entities and relationships
- `contracts/` - API specifications
- `research.md` - Technical decisions
- `quickstart.md` - Integration scenarios

### Step 4: Project Setup Verification

Verifies ignore files for detected technologies:

**Git Repository:**
- Checks `.gitignore` exists
- Verifies essential patterns present
- Warns about missing critical patterns

**Technology Detection:**
Based on `plan.md` tech stack and project files:
- Go: vendor/, *.exe, *.test, *.out
- Node.js: node_modules/, dist/, build/
- Python: __pycache__/, *.pyc, venv/
- etc.

### Step 5: Task Structure Parsing

Extracts from tasks.md:
- Number of phases (## Phase N)
- Number of incomplete tasks (- [ ])
- Number of completed tasks (- [X] or - [x])
- Total task count

**Output:**
```
[INFO] Found 6 phases
[INFO] Found 50 incomplete tasks
[INFO] Found 0 completed tasks
[SUCCESS] Total tasks in plan: 50
```

### Step 6-7: Task Execution

**Current Status:** Framework only

The script provides the validation framework but does not execute actual implementation tasks. Task execution requires:

- AI/LLM integration for code generation
- Test framework automation
- Build system integration
- Code analysis and validation

**For Manual Implementation:**
Refer to the loaded documents:
- Task list: `specs/NNN-feature/tasks.md`
- Implementation plan: `specs/NNN-feature/plan.md`

**For AI-Assisted Implementation:**
Use GitHub Copilot or similar tools with the context loaded by this script.

### Step 8: Completion Validation

**Current Status:** Framework only

Would verify:
- ✓ All required tasks completed
- ✓ Implementation matches specification
- ✓ Tests pass and coverage meets requirements
- ✓ Implementation follows technical plan

## Feature Branch Requirements

The script works with feature branches following the naming convention:

```
NNN-feature-name
```

Where `NNN` is a 3-digit number (e.g., `001-mock-query-results`).

**Alternative:** Set `SPECIFY_FEATURE` environment variable:

```bash
SPECIFY_FEATURE="001-feature-name" ./execute-implementation.sh
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success or user-requested halt |
| 1 | Error (prerequisites failed, invalid input, etc.) |

## Integration with Spec-Driven Development

This script is part of the SDD workflow:

1. **/speckit.specify** - Create feature specification
2. **/speckit.plan** - Generate implementation plan
3. **/speckit.tasks** - Break down into tasks
4. **execute-implementation.sh** ← This script (validate & prepare)
5. **Manual/AI implementation** - Execute tasks
6. **Code review & testing** - Validate implementation

## Limitations

### Current Limitations

1. **No Actual Task Execution**: The script validates and prepares but doesn't execute code changes
2. **Simple Ignore File Handling**: Basic pattern checking, not comprehensive
3. **Bash-Based**: Limited to what bash scripts can accomplish

### Why Not Full Automation?

Full task execution requires:
- Understanding code semantics
- Making architecture decisions
- Writing context-aware code
- Running and debugging tests
- Handling edge cases

These capabilities require AI/LLM integration or manual development.

## Future Enhancements

Potential improvements:

- [ ] Integration with AI coding assistants
- [ ] Automated test execution and validation
- [ ] Progress tracking across sessions
- [ ] Rollback capabilities
- [ ] Parallel task execution orchestration
- [ ] Build system integration
- [ ] Automated PR creation

## Troubleshooting

### "Not on a feature branch"

**Problem:** Current branch doesn't match `NNN-feature-name` pattern

**Solutions:**
```bash
# Option 1: Use environment variable
SPECIFY_FEATURE="001-feature-name" ./execute-implementation.sh

# Option 2: Checkout proper feature branch
git checkout 001-feature-name
```

### "Feature directory not found"

**Problem:** specs/NNN-feature/ doesn't exist

**Solution:**
```bash
# Create feature specification first
# (or manually create the directory structure)
```

### "plan.md not found"

**Problem:** Missing implementation plan

**Solution:**
Run `/speckit.plan` to generate the plan (if using the SDD workflow)

### "tasks.md not found"

**Problem:** Missing task breakdown

**Solution:**
Run `/speckit.tasks` to generate tasks (if using the SDD workflow)

## Examples

### Complete Workflow Example

```bash
# 1. Ensure you're on a feature branch
git checkout 001-new-feature

# 2. Run dry-run to see what will be validated
./.specify/scripts/bash/execute-implementation.sh --dry-run

# 3. Run actual validation
./.specify/scripts/bash/execute-implementation.sh

# 4. If checklists incomplete, you'll be prompted:
#    "Do you want to proceed with implementation anyway? (yes/no):"

# 5. If all passes, you'll see:
#    [SUCCESS] All checklists passed!
#    [INFO] Task execution would proceed...
```

### Using with Different Features

```bash
# Feature 001
SPECIFY_FEATURE="001-mock-query" ./execute-implementation.sh

# Feature 002
SPECIFY_FEATURE="002-pagination" ./execute-implementation.sh --yes

# Feature 003 with arguments
SPECIFY_FEATURE="003-auth" ./execute-implementation.sh -- "Use JWT tokens"
```

## Related Scripts

- `check-prerequisites.sh` - Validates feature setup prerequisites
- `common.sh` - Shared utility functions
- `create-new-feature.sh` - Creates new feature structure
- `setup-plan.sh` - Sets up implementation planning
- `update-agent-context.sh` - Updates agent context files

## Support

For issues or questions:
1. Check this documentation
2. Review the script source code
3. Check the SDD workflow documentation
4. Open an issue in the repository

## License

Same as the cc-tools project license.
