#!/usr/bin/env bash

# Execute Implementation Plan
#
# This script executes the implementation plan by processing and executing all tasks
# defined in tasks.md following the spec-driven development workflow.
#
# Usage: ./execute-implementation.sh [OPTIONS] [-- ARGUMENTS]
#
# OPTIONS:
#   --help, -h          Show help message
#   --dry-run           Show what would be executed without running
#   --yes               Skip confirmation prompts (proceed automatically)
#
# ARGUMENTS:
#   User input arguments to consider before proceeding
#
# WORKFLOW:
#   1. Check prerequisites (FEATURE_DIR, AVAILABLE_DOCS)
#   2. Validate checklists (halt if incomplete unless --yes)
#   3. Load implementation context
#   4. Verify project setup (ignore files)
#   5. Parse tasks.md
#   6. Execute tasks phase by phase
#   7. Track progress and handle errors
#   8. Validate completion

set -e

# Parse command line arguments
DRY_RUN=false
AUTO_YES=false
USER_ARGS=()

while [[ $# -gt 0 ]]; do
    case "$1" in
        --help|-h)
            cat << 'EOF'
Usage: execute-implementation.sh [OPTIONS] [-- ARGUMENTS]

Execute the implementation plan by processing tasks defined in tasks.md.

OPTIONS:
  --help, -h          Show this help message
  --dry-run           Show what would be executed without running
  --yes               Skip confirmation prompts (proceed automatically)

ARGUMENTS:
  User input arguments to consider before proceeding (after --)

EXAMPLES:
  # Execute implementation with user confirmation
  ./execute-implementation.sh
  
  # Execute with automatic yes to all prompts
  ./execute-implementation.sh --yes
  
  # Dry run to see execution plan
  ./execute-implementation.sh --dry-run
  
  # Execute with user input
  ./execute-implementation.sh -- "Implement with high test coverage"

WORKFLOW:
  1. Check prerequisites (FEATURE_DIR, AVAILABLE_DOCS)
  2. Validate checklists (halt if incomplete unless --yes)
  3. Load implementation context
  4. Verify project setup (ignore files)
  5. Parse tasks.md structure
  6. Execute tasks phase by phase
  7. Track progress and handle errors
  8. Validate completion

EOF
            exit 0
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --yes)
            AUTO_YES=true
            shift
            ;;
        --)
            shift
            USER_ARGS=("$@")
            break
            ;;
        *)
            echo "ERROR: Unknown option '$1'. Use --help for usage information." >&2
            exit 1
            ;;
    esac
done

# Source common functions
SCRIPT_DIR="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/common.sh"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
echo_success() { echo -e "${GREEN}[SUCCESS]${NC} $*"; }
echo_warning() { echo -e "${YELLOW}[WARNING]${NC} $*"; }
echo_error() { echo -e "${RED}[ERROR]${NC} $*" >&2; }

# ============================================================================
# STEP 1: Run prerequisites check
# ============================================================================

echo_info "Step 1: Checking prerequisites..."

# Run check-prerequisites.sh with required flags
if ! prereq_output=$("$SCRIPT_DIR/check-prerequisites.sh" --json --require-tasks --include-tasks 2>&1); then
    echo_error "Prerequisites check failed:"
    echo "$prereq_output"
    exit 1
fi

# Parse JSON output
FEATURE_DIR=$(echo "$prereq_output" | jq -r '.FEATURE_DIR')
AVAILABLE_DOCS=$(echo "$prereq_output" | jq -r '.AVAILABLE_DOCS[]' 2>/dev/null || echo "")

if [[ -z "$FEATURE_DIR" ]] || [[ "$FEATURE_DIR" == "null" ]]; then
    echo_error "Failed to determine FEATURE_DIR from prerequisites"
    exit 1
fi

echo_success "Feature directory: $FEATURE_DIR"
echo_info "Available documents: $(echo "$AVAILABLE_DOCS" | tr '\n' ' ')"

# ============================================================================
# STEP 2: Check checklists status
# ============================================================================

echo ""
echo_info "Step 2: Checking checklists status..."

CHECKLISTS_DIR="$FEATURE_DIR/checklists"
CHECKLIST_FAILED=false

if [[ -d "$CHECKLISTS_DIR" ]]; then
    echo ""
    echo "| Checklist | Total | Completed | Incomplete | Status |"
    echo "|-----------|-------|-----------|------------|--------|"
    
    overall_incomplete=0
    
    for checklist_file in "$CHECKLISTS_DIR"/*.md; do
        if [[ -f "$checklist_file" ]]; then
            filename=$(basename "$checklist_file")
            
            # Count completed items (- [X] or - [x])
            completed=$(grep -E '^[[:space:]]*-[[:space:]]*\[[Xx]\]' "$checklist_file" | wc -l)
            
            # Count incomplete items (- [ ])
            incomplete=$(grep -E '^[[:space:]]*-[[:space:]]*\[ \]' "$checklist_file" | wc -l)
            
            # Calculate total
            total=$((completed + incomplete))
            
            # Determine status
            if [[ $incomplete -eq 0 ]]; then
                status="✓ PASS"
            else
                status="✗ FAIL"
                CHECKLIST_FAILED=true
                overall_incomplete=$((overall_incomplete + incomplete))
            fi
            
            printf "| %-9s | %-5s | %-9s | %-10s | %-6s |\n" \
                "$filename" "$total" "$completed" "$incomplete" "$status"
        fi
    done
    
    echo ""
    
    # If any checklist is incomplete, ask user
    if $CHECKLIST_FAILED; then
        echo_warning "Some checklists are incomplete ($overall_incomplete items remaining)."
        
        if ! $AUTO_YES; then
            echo ""
            read -p "Do you want to proceed with implementation anyway? (yes/no): " response
            
            case "$response" in
                yes|y|Yes|YES|proceed|continue)
                    echo_info "Proceeding with implementation..."
                    ;;
                no|n|No|NO|wait|stop)
                    echo_info "Halting execution as requested."
                    exit 0
                    ;;
                *)
                    echo_error "Invalid response. Halting execution."
                    exit 1
                    ;;
            esac
        else
            echo_info "Auto-yes enabled: Proceeding despite incomplete checklists..."
        fi
    else
        echo_success "All checklists passed!"
    fi
else
    echo_info "No checklists directory found. Skipping checklist validation."
fi

# ============================================================================
# STEP 3: Load implementation context
# ============================================================================

echo ""
echo_info "Step 3: Loading implementation context..."

TASKS_FILE="$FEATURE_DIR/tasks.md"
PLAN_FILE="$FEATURE_DIR/plan.md"
DATA_MODEL_FILE="$FEATURE_DIR/data-model.md"
RESEARCH_FILE="$FEATURE_DIR/research.md"
QUICKSTART_FILE="$FEATURE_DIR/quickstart.md"
CONTRACTS_DIR="$FEATURE_DIR/contracts"

# Required files
if [[ ! -f "$TASKS_FILE" ]]; then
    echo_error "Required file not found: $TASKS_FILE"
    exit 1
fi

if [[ ! -f "$PLAN_FILE" ]]; then
    echo_error "Required file not found: $PLAN_FILE"
    exit 1
fi

echo_success "Loaded tasks.md"
echo_success "Loaded plan.md"

# Optional files
[[ -f "$DATA_MODEL_FILE" ]] && echo_info "Loaded data-model.md" || echo_info "data-model.md not found (optional)"
[[ -f "$RESEARCH_FILE" ]] && echo_info "Loaded research.md" || echo_info "research.md not found (optional)"
[[ -d "$CONTRACTS_DIR" ]] && echo_info "Loaded contracts/" || echo_info "contracts/ not found (optional)"
[[ -f "$QUICKSTART_FILE" ]] && echo_info "Loaded quickstart.md" || echo_info "quickstart.md not found (optional)"

# ============================================================================
# STEP 4: Project Setup Verification (Ignore Files)
# ============================================================================

echo ""
echo_info "Step 4: Verifying project setup (ignore files)..."

# Get repository root
eval $(get_feature_paths)

# Check if git repository
if git rev-parse --git-dir >/dev/null 2>&1; then
    echo_info "Git repository detected. Checking .gitignore..."
    
    GITIGNORE_FILE="$REPO_ROOT/.gitignore"
    
    # Essential patterns for Go project
    ESSENTIAL_PATTERNS=(
        "# Dependencies"
        "vendor/"
        "# Build outputs"
        "*.exe"
        "*.test"
        "*.out"
        "# IDE"
        ".vscode/"
        ".idea/"
        "# OS"
        ".DS_Store"
        "Thumbs.db"
        "# Temporary files"
        "*.tmp"
        "*.swp"
        "*.log"
        ".env*"
    )
    
    if [[ -f "$GITIGNORE_FILE" ]]; then
        echo_info ".gitignore already exists. Verifying essential patterns..."
        
        # Check for essential patterns (simplified check)
        missing_patterns=()
        for pattern in "vendor/" "*.exe" "*.test" ".DS_Store"; do
            if ! grep -q "$pattern" "$GITIGNORE_FILE"; then
                missing_patterns+=("$pattern")
            fi
        done
        
        if [[ ${#missing_patterns[@]} -gt 0 ]]; then
            echo_warning "Missing patterns in .gitignore: ${missing_patterns[*]}"
            echo_info "Consider adding these patterns manually."
        else
            echo_success ".gitignore contains essential patterns"
        fi
    else
        echo_info ".gitignore not found. This is acceptable."
    fi
else
    echo_info "Not a git repository. Skipping .gitignore check."
fi

# Note: For a minimal implementation, we're not creating new ignore files
# as this would require more context about the project structure

# ============================================================================
# STEP 5: Parse tasks.md structure
# ============================================================================

echo ""
echo_info "Step 5: Parsing tasks.md structure..."

# Extract task phases and task count
# Note: Using POSIX-compliant [[:space:]] instead of \s for portability
phase_count=$(grep -c '^## Phase' "$TASKS_FILE" 2>/dev/null) || phase_count=0
task_count=$(grep -E -c '^[[:space:]]*-[[:space:]]*\[ \]' "$TASKS_FILE" 2>/dev/null) || task_count=0
completed_count=$(grep -E -c '^[[:space:]]*-[[:space:]]*\[[Xx]\]' "$TASKS_FILE" 2>/dev/null) || completed_count=0

echo_info "Found $phase_count phases"
echo_info "Found $task_count incomplete tasks"  
echo_info "Found $completed_count completed tasks"

total_tasks=$((task_count + completed_count))
echo_success "Total tasks in plan: $total_tasks"

# ============================================================================
# STEP 6-7: Task Execution
# ============================================================================

echo ""
echo_info "Step 6-7: Task Execution"
echo ""

if $DRY_RUN; then
    echo_warning "DRY RUN MODE: Would execute tasks but not actually running them"
    echo_info "Task execution would proceed phase by phase following tasks.md"
    echo_info "Each phase would be completed before moving to the next"
    echo_info "Parallel tasks marked [P] would be executed concurrently"
    exit 0
fi

echo_warning "==========================================="
echo_warning "IMPORTANT: TASK EXECUTION NOT IMPLEMENTED"
echo_warning "==========================================="
echo ""
echo_info "This script currently validates prerequisites and context,"
echo_info "but does not execute the actual implementation tasks."
echo ""
echo_info "Actual task execution would require:"
echo_info "  - AI/LLM integration for code generation"
echo_info "  - Test framework integration"
echo_info "  - Build system integration"
echo_info "  - Code analysis and validation"
echo ""
echo_info "For now, this script serves as a framework for the execution workflow."
echo ""
echo_info "To implement tasks manually, refer to:"
echo_info "  - Task list: $TASKS_FILE"
echo_info "  - Implementation plan: $PLAN_FILE"
echo ""

# ============================================================================
# STEP 8: Completion Validation
# ============================================================================

echo_info "Step 8: Completion Validation"
echo ""
echo_info "Would verify:"
echo_info "  ✓ All required tasks completed"
echo_info "  ✓ Implementation matches specification"  
echo_info "  ✓ Tests pass and coverage meets requirements"
echo_info "  ✓ Implementation follows technical plan"
echo ""

echo_success "Script completed successfully!"
echo_info "Next steps: Implement actual task execution logic or execute tasks manually"

exit 0
