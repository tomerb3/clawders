#!/usr/bin/env bash

set -euo pipefail 2>/dev/null || set -eu  # Conditional pipefail for POSIX compatibility

# ============================================================
# CONFIGURATION
# ============================================================
readonly BAR_WIDTH=15
readonly BAR_FILLED="█"
readonly BAR_EMPTY="░"

readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly BLUE='\033[0;34m'
readonly MAGENTA='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly ORANGE='\033[0;33m'
readonly GRAY='\033[0;90m'
readonly NC='\033[0m'

# 256-color palette
readonly GIT_PURPLE='\033[38;5;99m'
readonly TOKENS_ORANGE='\033[38;5;208m'
readonly CLOCK_PURPLE='\033[38;5;105m'


# Derived constants
readonly SEPARATOR="${GRAY}|${NC}"
readonly NULL_VALUE="null"

# Icons
readonly MODEL_ICON="🚀"
readonly CONTEXT_ICON="🔥"
readonly DIR_ICON="📂"
readonly GIT_ICON=$'\ue0a0'

# Git state constants
readonly STATE_NOT_REPO="not_repo"
readonly STATE_CLEAN="clean"
readonly STATE_DIRTY="dirty"

# ============================================================
# UTILITY FUNCTIONS
# ============================================================

# String utilities
get_dirname() { echo "${1##*/}"; }
sep() { echo -n " ${SEPARATOR} "; }

# Calculate visible width of a string (strips ANSI escape codes, accounts for wide chars)
visible_width() {
  local str="$1"
  # Remove ANSI escape codes
  local stripped
  stripped=$(echo -e "${str}" | sed 's/\x1b\[[0-9;]*m//g')
  # Count base characters
  local base_width="${#stripped}"
  # Count wide characters (emojis) that take 2 columns - add 1 for each
  # Common emoji ranges: most emojis display as 2 columns wide
  local wide_chars
  wide_chars=$(echo -n "${stripped}" | grep -o '[📂💵✏️🚀🔥🎋⏱]' 2>/dev/null | wc -l | tr -d ' ')
  echo $((base_width + wide_chars))
}

# Pad a string to a target visible width
pad_to_width() {
  local str="$1"
  local target_width="$2"
  local current_width
  current_width=$(visible_width "${str}")
  local padding=$((target_width - current_width))
  if [[ "${padding}" -gt 0 ]]; then
    printf "%s%${padding}s" "${str}" ""
  else
    echo -n "${str}"
  fi
}

# Conditional append helper (DRY pattern)
append_if() {
  local value="$1"
  local text="$2"
  if [[ "${value}" != "0" ]] 2>/dev/null && [[ -n "${value}" ]] && [[ "${value}" != "${NULL_VALUE}" ]]; then
    echo -n " ${text}"
  fi
}

# Format numbers with K/M suffixes for readability
# Examples: 543 -> "543", 1500 -> "1.5K", 54000 -> "54K", 1200000 -> "1.2M"
format_number() {
  local num="$1"

  if [[ "${num}" -lt 1000 ]]; then
    echo "${num}"
  elif [[ "${num}" -lt 1000000 ]]; then
    # Thousands
    local k=$((num / 1000))
    local remainder=$((num % 1000))
    if [[ "${k}" -lt 10 ]]; then
      # Show decimal for < 10K
      local decimal=$((remainder / 100))
      echo "${k}.${decimal}K"
    else
      echo "${k}K"
    fi
  else
    # Millions
    local m=$((num / 1000000))
    local remainder=$((num % 1000000))
    if [[ "${m}" -lt 10 ]]; then
      # Show decimal for < 10M
      local decimal=$((remainder / 100000))
      echo "${m}.${decimal}M"
    else
      echo "${m}M"
    fi
  fi
}

# Check git version for porcelain v2 support (requires git 2.11+)
# Cache result for performance
check_git_version() {
  # Return cached result if available
  [[ -n "${GIT_VERSION_CHECKED:-}" ]] && return "${GIT_VERSION_OK:-1}"

  GIT_VERSION_CHECKED=1
  command -v git >/dev/null 2>&1 || { GIT_VERSION_OK=1; return 1; }

  local version
  version=$(git --version 2>/dev/null | awk '{print $3}')
  [[ -z "${version}" ]] && { GIT_VERSION_OK=1; return 1; }

  # Semantic version comparison: >= 2.11
  local major minor
  IFS='.' read -r major minor _ << EOF
${version}
EOF

  if [[ "${major}" -gt 2 ]] || { [[ "${major}" -eq 2 ]] && [[ "${minor}" -ge 11 ]]; }; then
    GIT_VERSION_OK=0
    return 0
  else
    GIT_VERSION_OK=1
    return 1
  fi
}

# ============================================================
# FUNCTIONS
# ============================================================

parse_claude_input() {
  local input="$1"

  local parsed
  parsed=$(echo "${input}" | jq -r '
    .model.display_name,
    .workspace.current_dir,
    (.context_window.context_window_size // 200000),
    (
      (.context_window.current_usage.input_tokens // 0) +
      (.context_window.current_usage.cache_creation_input_tokens // 0) +
      (.context_window.current_usage.cache_read_input_tokens // 0)
    ),
    (.cost.total_cost_usd // 0),
    (.cost.total_lines_added // 0),
    (.cost.total_lines_removed // 0),
    (.context_window.total_input_tokens // 0),
    (.context_window.total_output_tokens // 0),
    (.cost.total_duration_ms // 0),
    (.effort.level // "")
  ' 2>/dev/null) || {
    echo "Error: Failed to parse JSON input" >&2
    return 1
  }

  echo "${parsed}"
}

build_progress_bar() {
  local percent="$1"
  local filled=$((percent * BAR_WIDTH / 100))
  local empty=$((BAR_WIDTH - filled))

  # Determine bar color based on percentage
  local bar_color
  if [[ "${percent}" -le 20 ]]; then
    bar_color="${GREEN}"
  elif [[ "${percent}" -le 40 ]]; then
    bar_color="${CYAN}"
  elif [[ "${percent}" -le 60 ]]; then
    bar_color="${ORANGE}"
  elif [[ "${percent}" -le 80 ]]; then
    bar_color="${ORANGE}"
  else
    bar_color="${RED}"
  fi

  # Build colored filled portion and gray empty portion
  echo -n "${bar_color}"
  printf "%${filled}s" | tr ' ' "${BAR_FILLED}"
  echo -n "${NC}${GRAY}"
  printf "%${empty}s" | tr ' ' "${BAR_EMPTY}"
  echo -n "${NC}"
}

# ============================================================
# GIT OPERATIONS (Optimized - 7 calls reduced to 2)
# ============================================================

get_git_info() {
  local current_dir="$1"
  local git_opts=()

  [[ -n "${current_dir}" ]] && [[ "${current_dir}" != "${NULL_VALUE}" ]] && git_opts=(-C "${current_dir}")

  # Check if git repo
  git "${git_opts[@]}" rev-parse --is-inside-work-tree >/dev/null 2>&1 || {
    echo "${STATE_NOT_REPO}"
    return 0
  }

  # Single git status call with all info (replaces 5 separate calls)
  # Requires Git 2.11+ (Dec 2016) for --porcelain=v2
  local status_output
  status_output=$(git "${git_opts[@]}" status --porcelain=v2 --branch --untracked-files=all 2>/dev/null) || {
    echo "${STATE_NOT_REPO}"
    return 0
  }

  # Parse porcelain v2 output
  local branch ahead behind
  while IFS= read -r line; do
    case "${line}" in
      "# branch.head "*)
        branch="${line#\# branch.head }"
        ;;
      "# branch.ab "*)
        local ab="${line#\# branch.ab }"
        ahead="${ab%% *}"
        ahead="${ahead#+}"
        behind="${ab##* }"
        behind="${behind#-}"
        ;;
      *)
        # Ignore other porcelain output lines
        ;;
    esac
  done << EOF
${status_output}
EOF

  # Default values
  branch="${branch:-(detached HEAD)}"
  ahead="${ahead:-0}"
  behind="${behind:-0}"

  # Count modified files (lines not starting with #)
  local file_lines
  file_lines=$(echo "${status_output}" | grep -v '^#')
  local total_files=0
  [[ -n "${file_lines}" ]] && total_files=$(echo "${file_lines}" | wc -l | tr -d ' ')

  # Clean state if no files
  if [[ "${total_files}" -eq 0 ]]; then
    echo "${STATE_CLEAN}|${branch}|${ahead}|${behind}"
    return 0
  fi

  # Get line changes (single diff HEAD call replaces 2 separate cached + unstaged calls)
  local added removed
  read -r added removed << EOF
$(git "${git_opts[@]}" diff HEAD --numstat 2>/dev/null | awk '{a+=$1; r+=$2} END {print a+0, r+0}' || true)
EOF

  echo "${STATE_DIRTY}|${branch}|${total_files}|${added}|${removed}|${ahead}|${behind}"
}

# ============================================================
# FORMATTING FUNCTIONS (SOLID - Single Responsibility)
# ============================================================

format_ahead_behind() {
  local ahead="$1"
  local behind="$2"
  local output=""

  [[ "${ahead}" -gt 0 ]] 2>/dev/null && output+=" ${GREEN}↑${ahead}${NC}"
  [[ "${behind}" -gt 0 ]] 2>/dev/null && output+=" ${RED}↓${behind}${NC}"

  [[ -n "${output}" ]] && echo "${GRAY}|${NC}${output}"
}

format_git_not_repo() {
  echo " ${ORANGE}(not a git repository)${NC}"
}

format_git_clean() {
  local branch="$1" ahead="$2" behind="$3"

  # Simple format: branch + ahead/behind (no parentheses)
  local output="${GIT_PURPLE}${branch}${NC}"
  local ahead_behind
  ahead_behind=$(format_ahead_behind "${ahead}" "${behind}")
  [[ -n "${ahead_behind}" ]] && output+="${ahead_behind}"

  echo " ${output}"
}

format_git_dirty() {
  local branch="$1" files="$2" added="$3" removed="$4" ahead="$5" behind="$6"

  # Simple branch + ahead/behind (no file count, no line changes)
  local output="${GIT_PURPLE}${branch}${NC}"
  local ahead_behind
  ahead_behind=$(format_ahead_behind "${ahead}" "${behind}")
  [[ -n "${ahead_behind}" ]] && output+="${ahead_behind}"

  # Return git info and file count separately: "git_output|file_count"
  echo " ${output}|${files}"
}

format_git_info() {
  local git_data="$1"

  # Parse state
  local state
  IFS='|' read -r state _ << EOF
${git_data}
EOF

  case "${state}" in
    "${STATE_NOT_REPO}")
      format_git_not_repo
      echo ""  # No file count
      ;;
    "${STATE_CLEAN}")
      local branch ahead behind
      IFS='|' read -r _ branch ahead behind << EOF
${git_data}
EOF
      format_git_clean "${branch}" "${ahead}" "${behind}"
      echo ""  # No file count for clean repo
      ;;
    "${STATE_DIRTY}")
      local branch files added removed ahead behind
      IFS='|' read -r _ branch files added removed ahead behind << EOF
${git_data}
EOF
      # Returns "git_output|file_count"
      format_git_dirty "${branch}" "${files}" "${added}" "${removed}" "${ahead}" "${behind}"
      ;;
    *)
      # Unknown state - show error
      echo " ${ORANGE}(unknown git state)${NC}"
      echo ""  # No file count
      ;;
  esac
}

# ============================================================
# COMPONENT BUILDERS (Open/Closed Principle)
# ============================================================

build_model_component() {
  local model_name="$1"
  local effort_level="${2:-}"
  # Strip "Claude " prefix and any parenthetical suffix e.g. " (1M context)"
  local short_name
  short_name=$(echo "${model_name}" | sed 's/^Claude //; s/ (.*//')
  local output="${CYAN}${short_name}${NC}"
  [[ -n "${effort_level}" ]] && output+=" ${GRAY}(${effort_level})${NC}"
  echo "${output}"
}

build_context_component() {
  local context_size="$1"
  local current_usage="$2"

  local context_percent=0
  if [[ "${current_usage}" != "0" && "${context_size}" -gt 0 ]]; then
    context_percent=$((current_usage * 100 / context_size))
  fi

  # Get colored progress bar
  local bar
  bar=$(build_progress_bar "${context_percent}")

  # Format usage numbers (e.g., "54K/200K")
  local usage_formatted
  usage_formatted=$(format_number "${current_usage}")
  local size_formatted
  size_formatted=$(format_number "${context_size}")

  # Output with brackets, colored bar, formatted numbers (no emoji, no redundant percentage)
  echo "${GRAY}[${NC}${bar}${GRAY}]${NC} ${context_percent}% ${usage_formatted}/${size_formatted}"
}

build_directory_component() {
  local current_dir="$1"

  local dir_name
  if [[ -n "${current_dir}" ]] && [[ "${current_dir}" != "${NULL_VALUE}" ]]; then
    dir_name=$(get_dirname "${current_dir}")
  else
    dir_name=$(get_dirname "${PWD}")
  fi

  echo "${DIR_ICON} ${BLUE}${dir_name}${NC}"
}

build_git_component() {
  local current_dir="$1"
  local git_data

  git_data=$(get_git_info "${current_dir}")

  # format_git_info returns two lines: git_output and file_count
  local formatted git_line file_line
  formatted=$(format_git_info "${git_data}")
  git_line=$(echo "${formatted}" | sed -n '1p')
  file_line=$(echo "${formatted}" | sed -n '2p')

  # Extract state to determine emoji placement
  local state
  IFS='|' read -r state _ << EOF
${git_data}
EOF

  # Return git info and file count separately: "git_display|file_count"
  if [[ "${state}" = "${STATE_NOT_REPO}" ]]; then
    echo "${git_line}|"
  else
    echo " ${GIT_PURPLE}${GIT_ICON}${NC}${git_line}|${file_line}"
  fi
}

build_files_component() {
  local file_count="$1"

  # Only show if there are modified files
  if [[ -n "${file_count}" && "${file_count}" != "0" ]]; then
    echo "${GRAY}${file_count} files${NC}"
  fi
}

build_cost_component() {
  local cost_usd="$1"

  if [[ -n "${cost_usd}" && "${cost_usd}" != "0" && "${cost_usd}" != "${NULL_VALUE}" ]]; then
    echo "💵 ${MAGENTA}\$$(printf "%.2f" "${cost_usd}")${NC}"
  fi
}

build_litellm_component() {
  local script_path="$(dirname "$0")/litellm-budget.sh"
  [[ -x "$script_path" ]] || return 0
  
  # Run with --status for compact one-liner, suppress errors
  local lightllm_budget
  lightllm_budget=$("$script_path" --status 2>/dev/null) || return 0
  [[ -n "$lightllm_budget" && "$lightllm_budget" != budget:* ]] || return 0
  lightllm_budget=$(echo "$lightllm_budget" | sed 's/^\$//; s/ \/ /\//')

  echo "💰 ${lightllm_budget}"
}

build_lines_component() {
  local lines_added="$1"
  local lines_removed="$2"

  if [[ -n "${lines_added}" && -n "${lines_removed}" ]] && \
     [[ "${lines_added}" != "0" || "${lines_removed}" != "0" ]] && \
     [[ "${lines_added}" != "${NULL_VALUE}" && "${lines_removed}" != "${NULL_VALUE}" ]]; then
    echo "${GREEN}+${lines_added}${NC}/${RED}-${lines_removed}${NC}"
  fi
}

build_tokens_component() {
  local total_input="$1"
  local total_output="$2"

  if [[ -n "${total_input}" && -n "${total_output}" ]] && \
     [[ "${total_input}" != "0" || "${total_output}" != "0" ]] && \
     [[ "${total_input}" != "${NULL_VALUE}" && "${total_output}" != "${NULL_VALUE}" ]]; then
    local input_formatted output_formatted
    input_formatted=$(format_number "${total_input}")
    output_formatted=$(format_number "${total_output}")
    echo "${TOKENS_ORANGE}↓${input_formatted}${NC} ${TOKENS_ORANGE}↑${output_formatted}${NC}"
  fi
}

build_duration_component() {
  local duration_ms="$1"

  if [[ -n "${duration_ms}" && "${duration_ms}" != "0" && "${duration_ms}" != "${NULL_VALUE}" ]]; then
    # Convert ms to human readable format
    local total_seconds=$((duration_ms / 1000))
    local hours=$((total_seconds / 3600))
    local minutes=$(( (total_seconds % 3600) / 60 ))
    local seconds=$((total_seconds % 60))

    local formatted
    if [[ "${hours}" -gt 0 ]]; then
      formatted="${hours}h ${minutes}m"
    elif [[ "${minutes}" -gt 0 ]]; then
      formatted="${minutes}m ${seconds}s"
    else
      formatted="${seconds}s"
    fi
    echo "${CLOCK_PURPLE}⏱ ${formatted}${NC}"
  fi
}

# ============================================================
# ASSEMBLY (KISS - Simple orchestration)
# ============================================================

assemble_statusline() {
  local model_part="$1"
  local context_part="$2"
  local dir_part="$3"
  local git_part="$4"
  local files_part="$5"
  local cost_part="$6"
  local lightllm_budget_part="$7"
  local lines_part="$8"
  local tokens_part="$9"
  local duration_part="${10}"

  local separator
  separator=$(sep)

  # Calculate widths for alignment across all corresponding columns
  # Column 1: model vs directory
  local model_width dir_width col1_width
  model_width=$(visible_width "${model_part}")
  dir_width=$(visible_width "${dir_part}")
  col1_width=$((model_width > dir_width ? model_width : dir_width))

  # Column 2: context vs git
  local context_width git_width col2_width
  context_width=$(visible_width "${context_part}")
  git_width=$(visible_width "${git_part}")
  col2_width=$((context_width > git_width ? context_width : git_width))

  # Column 3: tokens vs files
  local tokens_width files_width col3_width
  tokens_width=$(visible_width "${tokens_part}")
  files_width=$(visible_width "${files_part}")
  col3_width=$((tokens_width > files_width ? tokens_width : files_width))

  # Column 4: duration vs lines
  local duration_width lines_width col4_width
  duration_width=$(visible_width "${duration_part}")
  lines_width=$(visible_width "${lines_part}")
  col4_width=$((duration_width > lines_width ? duration_width : lines_width))

  # Build line 1 with padding: model | context | tokens | duration
  local line1
  line1="$(pad_to_width "${model_part}" "${col1_width}")${separator}"
  line1+="$(pad_to_width "${context_part}" "${col2_width}")"
  [[ -n "${tokens_part}" || -n "${files_part}" ]] && line1+="${separator}$(pad_to_width "${tokens_part}" "${col3_width}")"
  [[ -n "${duration_part}" || -n "${lines_part}" ]] && line1+="${separator}$(pad_to_width "${duration_part}" "${col4_width}")"

  # Build line 2 with padding: directory | git | files | lines | cost
  local line2
  line2="$(pad_to_width "${dir_part}" "${col1_width}")${separator}"
  line2+="$(pad_to_width "${git_part}" "${col2_width}")"
  [[ -n "${tokens_part}" || -n "${files_part}" ]] && line2+="${separator}$(pad_to_width "${files_part}" "${col3_width}")"
  [[ -n "${duration_part}" || -n "${lines_part}" ]] && line2+="${separator}$(pad_to_width "${lines_part}" "${col4_width}")"
  [[ -n "${cost_part}" ]] && line2+="${separator}${cost_part}"
  [[ -n "${lightllm_budget_part}" ]] && line2+="${separator}${lightllm_budget_part}"

  echo -e "${line1}"
  echo -e "${line2}"
}

# ============================================================
# MAIN (Simplified orchestration only)
# ============================================================

main() {
  # Check dependencies
  command -v jq >/dev/null 2>&1 || {
    echo "Error: jq required" >&2
    exit 1
  }

  # Read input (POSIX-compatible: cat instead of < /dev/stdin)
  local input
  input=$(cat) || {
    echo "Error: Failed to read stdin" >&2
    exit 1
  }

  # Parse JSON
  local parsed
  parsed=$(parse_claude_input "${input}")
  if [[ -z "${parsed}" ]]; then
    exit 1
  fi

  # Extract fields
  local model_name="" current_dir="" context_size="" current_usage="" cost_usd="" lines_added="" lines_removed=""
  local total_input_tokens="" total_output_tokens="" total_duration_ms="" effort_level=""
  {
    read -r model_name
    read -r current_dir
    read -r context_size
    read -r current_usage
    read -r cost_usd
    read -r lines_added
    read -r lines_removed
    read -r total_input_tokens
    read -r total_output_tokens
    read -r total_duration_ms
    read -r effort_level || true
  } << EOF
${parsed}
EOF

  # Build components
  local model_part context_part dir_part git_part cost_part files_part lines_part tokens_part duration_part
  model_part=$(build_model_component "${model_name}" "${effort_level}")
  context_part=$(build_context_component "${context_size}" "${current_usage}")
  dir_part=$(build_directory_component "${current_dir}")

  # Git component returns "git_display|file_count"
  local git_with_files file_count
  git_with_files=$(build_git_component "${current_dir}")
  IFS='|' read -r git_part file_count <<< "${git_with_files}"

  files_part=$(build_files_component "${file_count}")
  cost_part=$(build_cost_component "${cost_usd}")
  litellm_budget_part=$(build_litellm_component)
  lines_part=$(build_lines_component "${lines_added}" "${lines_removed}")
  tokens_part=$(build_tokens_component "${total_input_tokens}" "${total_output_tokens}")
  duration_part=$(build_duration_component "${total_duration_ms}")

  # Assemble and output (2 lines)
  assemble_statusline "${model_part}" "${context_part}" "${dir_part}" "${git_part}" "${files_part}" "${cost_part}" "${litellm_budget_part}" "${lines_part}" "${tokens_part}" "${duration_part}"
}

main "$@"
