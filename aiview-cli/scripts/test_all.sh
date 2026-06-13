#!/usr/bin/env bash
# scripts/test_all.sh
# End-to-end automated test script for aiview CLI
# Usage: ./scripts/test_all.sh [--skip-build]

set -euo pipefail

# ─── Colors ──────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# ─── Config ──────────────────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

BINARY="./aiview"
SKIP_BUILD=false
PASS_COUNT=0
FAIL_COUNT=0
SKIP_COUNT=0
RESULTS=()

# ─── Parse args ──────────────────────────────────────────────────────────────
for arg in "$@"; do
  case $arg in
    --skip-build) SKIP_BUILD=true ;;
    *) echo "Unknown option: $arg"; exit 1 ;;
  esac
done

# ─── Helpers ─────────────────────────────────────────────────────────────────
log_info()  { echo -e "${CYAN}[INFO]${NC}  $*"; }
log_ok()    { echo -e "${GREEN}[PASS]${NC}  $*"; }
log_fail()  { echo -e "${RED}[FAIL]${NC}  $*"; }
log_skip()  { echo -e "${YELLOW}[SKIP]${NC}  $*"; }

record_result() {
  local name="$1" status="$2" detail="${3:-}"
  RESULTS+=("$name|$status|$detail")
  if [[ "$status" == "PASS" ]]; then
    ((PASS_COUNT++))
  elif [[ "$status" == "FAIL" ]]; then
    ((FAIL_COUNT++))
  else
    ((SKIP_COUNT++))
  fi
}

# Run a command and check exit code + optional output pattern
run_test() {
  local name="$1"
  shift
  local cmd="$*"

  log_info "Running: $cmd"
  set +e
  output=$(eval "$cmd" 2>&1)
  exit_code=$?
  set -e

  if [[ $exit_code -ne 0 ]]; then
    log_fail "$name (exit code: $exit_code)"
    echo "  Output: $(echo "$output" | head -5)"
    record_result "$name" "FAIL" "exit_code=$exit_code"
    return 1
  fi

  log_ok "$name"
  record_result "$name" "PASS" ""
  return 0
}

# Run a command and check that output contains a pattern
run_test_expect() {
  local name="$1"
  local pattern="$2"
  shift 2
  local cmd="$*"

  log_info "Running: $cmd"
  set +e
  output=$(eval "$cmd" 2>&1)
  exit_code=$?
  set -e

  if [[ $exit_code -ne 0 ]]; then
    log_fail "$name (exit code: $exit_code)"
    echo "  Output: $(echo "$output" | head -5)"
    record_result "$name" "FAIL" "exit_code=$exit_code"
    return 1
  fi

  if echo "$output" | grep -qE "$pattern"; then
    log_ok "$name"
    record_result "$name" "PASS" ""
    return 0
  else
    log_fail "$name (pattern '$pattern' not found)"
    echo "  Output: $(echo "$output" | head -10)"
    record_result "$name" "FAIL" "pattern_not_found=$pattern"
    return 1
  fi
}

# ─── Build ───────────────────────────────────────────────────────────────────
build_binary() {
  if [[ "$SKIP_BUILD" == "true" ]]; then
    log_info "Skipping build (--skip-build)"
    if [[ ! -x "$BINARY" ]]; then
      log_fail "Binary not found: $BINARY"
      exit 1
    fi
    return 0
  fi

  log_info "Building aiview..."
  if ! go build -o "$BINARY" .; then
    log_fail "Build failed"
    exit 1
  fi
  log_ok "Build successful"
}

# ─── Test: Hot commands ─────────────────────────────────────────────────────
test_hot_commands() {
  echo ""
  echo "============================================"
  echo "  Platform Hot Commands"
  echo "============================================"

  run_test "bilibili hot" "$BINARY bilibili hot" || true
  run_test "bilibili hot --json" "$BINARY bilibili hot --json" || true
  run_test "douyin hot" "$BINARY douyin hot" || true
  run_test "douyin hot --json" "$BINARY douyin hot --json" || true
  run_test "xiaohongshu hot" "$BINARY xiaohongshu hot" || true
  run_test "xiaohongshu hot --json" "$BINARY xiaohongshu hot --json" || true
}

# ─── Test: Global command help ───────────────────────────────────────────────
test_global_help() {
  echo ""
  echo "============================================"
  echo "  Global Command Help"
  echo "============================================"

  run_test_expect "analyze trend --help" "Usage|usage|trend" "$BINARY analyze trend --help" || true
  run_test_expect "compare --help" "Usage|usage|compare" "$BINARY compare --help" || true
  run_test_expect "schedule --help" "Usage|usage|schedule" "$BINARY schedule --help" || true
  run_test_expect "dashboard --help" "Usage|usage|dashboard" "$BINARY dashboard --help" || true
  run_test_expect "tui --help" "Usage|usage|tui" "$BINARY tui --help" || true
  run_test_expect "export --help" "Usage|usage|export" "$BINARY export --help" || true
  run_test_expect "root --help" "Usage|usage|aiview" "$BINARY --help" || true
}

# ─── Test: Subcommand help ──────────────────────────────────────────────────
test_subcommand_help() {
  echo ""
  echo "============================================"
  echo "  Subcommand Help"
  echo "============================================"

  run_test_expect "schedule add --help" "Usage|usage|add" "$BINARY schedule add --help" || true
  run_test_expect "schedule list --help" "Usage|usage|list" "$BINARY schedule list --help" || true
  run_test_expect "schedule remove --help" "Usage|usage|remove" "$BINARY schedule remove --help" || true
}

# ─── Test: Version / unknown command ────────────────────────────────────────
test_misc() {
  echo ""
  echo "============================================"
  echo "  Miscellaneous"
  echo "============================================"

  # Unknown command should fail (non-zero exit)
  log_info "Running: $BINARY nonexistent_command_xyz"
  set +e
  output=$("$BINARY" nonexistent_command_xyz 2>&1)
  exit_code=$?
  set -e
  if [[ $exit_code -ne 0 ]]; then
    log_ok "unknown command returns non-zero exit"
    record_result "unknown command" "PASS" ""
  else
    log_fail "unknown command should return non-zero exit"
    record_result "unknown command" "FAIL" "exit_code=$exit_code"
  fi
}

# ─── Summary ─────────────────────────────────────────────────────────────────
print_summary() {
  local total=$((PASS_COUNT + FAIL_COUNT + SKIP_COUNT))
  echo ""
  echo "============================================"
  echo "  Test Summary"
  echo "============================================"
  echo -e "  Total : $total"
  echo -e "  ${GREEN}Passed${NC}: $PASS_COUNT"
  echo -e "  ${RED}Failed${NC}: $FAIL_COUNT"
  echo -e "  ${YELLOW}Skipped${NC}: $SKIP_COUNT"
  echo "============================================"

  if [[ $FAIL_COUNT -gt 0 ]]; then
    echo ""
    echo "Failed tests:"
    for entry in "${RESULTS[@]}"; do
      IFS='|' read -r name status detail <<< "$entry"
      if [[ "$status" == "FAIL" ]]; then
        echo -e "  ${RED}✗${NC} $name  ($detail)"
      fi
    done
  fi

  echo ""
  if [[ $FAIL_COUNT -gt 0 ]]; then
    echo -e "${RED}TEST SUITE FAILED${NC}"
    exit 1
  else
    echo -e "${GREEN}ALL TESTS PASSED${NC}"
    exit 0
  fi
}

# ─── Main ────────────────────────────────────────────────────────────────────
main() {
  echo "============================================"
  echo "  aiview E2E Test Suite"
  echo "  $(date '+%Y-%m-%d %H:%M:%S')"
  echo "============================================"

  build_binary
  test_hot_commands
  test_global_help
  test_subcommand_help
  test_misc
  print_summary
}

main "$@"
