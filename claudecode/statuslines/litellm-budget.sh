#!/usr/bin/env bash
# Fetches LiteLLM budget info for Claude Code status bar, caching for 60s.
#
# Modes:
#   --status  Compact one-liner for status bar (default)
#   --full    Multi-line breakdown
#
# Requires: ANTHROPIC_BASE_URL and ANTHROPIC_AUTH_TOKEN env vars (Claude Code sets these)
set -euo pipefail

CACHE_TTL="${LITELLM_BUDGET_CACHE_TTL:-60}"
MODE="--status"

while (( $# > 0 )); do
  case "$1" in
    --status|--full) MODE="$1"; shift ;;
    -h|--help)
      echo "usage: $0 [--status|--full]"
      echo "  --status  Compact one-liner for status bar (default)"
      echo "  --full    Multi-line breakdown"
      exit 0 ;;
    *) echo "usage: $0 [--status|--full]" >&2; exit 2 ;;
  esac
done

# Check Claude Code env vars
if [[ -z "${ANTHROPIC_BASE_URL:-}" || -z "${ANTHROPIC_AUTH_TOKEN:-}" ]]; then
  echo "budget: ANTHROPIC_BASE_URL / ANTHROPIC_AUTH_TOKEN not set"
  exit 0
fi

CACHE_FILE="${TMPDIR:-/tmp}/litellm-budget-claude.json"

mtime() {
  if stat -f %m "$1" >/dev/null 2>&1; then stat -f %m "$1"
  else stat -c %Y "$1"
  fi
}

cache_fresh() {
  [[ -f "$CACHE_FILE" ]] || return 1
  local age=$(( $(date +%s) - $(mtime "$CACHE_FILE") ))
  (( age < CACHE_TTL ))
}

fetch() {
  curl -fsS --max-time 3 \
    -H "Authorization: Bearer ${ANTHROPIC_AUTH_TOKEN}" \
    "${ANTHROPIC_BASE_URL%/}/key/info" \
    -o "$CACHE_FILE.tmp" && mv "$CACHE_FILE.tmp" "$CACHE_FILE"
}

if ! cache_fresh; then
  if ! fetch 2>/dev/null && [[ ! -f "$CACHE_FILE" ]]; then
    echo "budget: unavailable"
    exit 0
  fi
fi

spend=$(jq -r '.info.spend // 0' "$CACHE_FILE")
budget=$(jq -r '.info.max_budget // 0' "$CACHE_FILE")
reset=$(jq -r '.info.budget_reset_at // "n/a"' "$CACHE_FILE")
duration=$(jq -r '.info.budget_duration // "n/a"' "$CACHE_FILE")

if awk "BEGIN{exit !($budget > 0)}"; then
  pct=$(awk "BEGIN{printf \"%.0f\", ($spend/$budget)*100}")
else
  pct=""
fi

case "$MODE" in
  --status)
    if [[ -z "$pct" ]]; then
      printf '%.2f' "$spend"
    else
      if   (( pct >= 90 )); then color=$'\033[38;2;160;0;0m'    # dark red
      elif (( pct >= 75 )); then color=$'\033[38;2;160;120;0m'  # dark yellow
      else                       color=$'\033[38;2;0;120;0m'  # dark green
      fi
      reset_color=$'\033[0m'
      printf '%s$%.2f / $%.2f (%d%%)%s' "$color" "$spend" "$budget" "$pct" "$reset_color"
    fi
    ;;
  --full)
    printf 'LiteLLM budget (claude)\n'
    printf '  Spent:    $%.2f\n' "$spend"
    if [[ -n "$pct" ]]; then
      printf '  Budget:   $%.2f (%d%% used)\n' "$budget" "$pct"
    else
      printf '  Budget:   (no cap set)\n'
    fi
    printf '  Window:   %s\n' "$duration"
    printf '  Resets:   %s\n' "$reset"
    printf '  Endpoint: %s\n' "$ANTHROPIC_BASE_URL"
    ;;
esac
