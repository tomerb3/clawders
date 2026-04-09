# 🤖 Claude Code Guide

A curated collection of tips, configurations, and workflows for getting the most out of [Claude Code](https://docs.anthropic.com/en/docs/claude-code) — Anthropic's official CLI agent.

---

## Table of Contents

- [Quick Start](#quick-start)
- [Basic Configuration](#basic-configuration)
- [Plugins & Skills](#plugins--skills)
- [Workflows & Productivity](#workflows--productivity)
- [Status Line & Terminal](#status-line--terminal)
- [Advanced Features](#advanced-features)
- [Remote Access](#remote-access)
- [Useful CLI Tools to Pair With Claude Code](#useful-cli-tools-to-pair-with-claude-code)
- [External Resources](#external-resources)

---

## Quick Start

### First-Time Setup

Add this alias to your `~/.zshrc` or `~/.bashrc`:

```bash
alias cc="claude --dangerously-skip-permissions $@"
```

> ⚠️ **If you are new to Claude Code, do NOT use the `--dangerously-skip-permissions` flag.** Start with the default permission prompts until you understand how Claude Code works.

### Starting Claude Code

```bash
cc
```

Once running, you can switch modes:
- `plan` — review actions before executing
- `accept edits on` — allow direct file modifications
- `bypass permissions on` — skip permission prompts (advanced only)

---

## Basic Configuration

### Essential Settings (`settings.json`)

A recommended base configuration:

```json
{
  "env": {
    "ANTHROPIC_MODEL": "MiniMax-M2.7",
    "ANTHROPIC_SMALL_FAST_MODEL": "MiniMax-M2.7",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "MiniMax-M2.7",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "MiniMax-M2.7",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "MiniMax-M2.7",
    "ANTHROPIC_BASE_URL": "https://api.minimax.io/anthropic",
    "ANTHROPIC_AUTH_TOKEN": "your-token-here",
    "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
    "DISABLE_TELEMETRY": "1",
    "DISABLE_ERROR_REPORTING": "1",
    "CLAUDE_CODE_DISABLE_FEEDBACK_SERVEY": "1",
    "CLAUDE_CODE_FILE_READ_MAX_OUTPUT_TOKENS": "100000",
    "CLAUDE_AUTOCOMPACT_PCT_OVERRIDE": "75"
  },
  "skipDangerousModePermissionPrompt": true
}
```

### Environment Variables

| Variable | Purpose |
|---------|---------|
| `ANTHROPIC_AUTH_TOKEN` | API token for your model provider |
| `ANTHROPIC_BASE_URL` | Override the API endpoint (e.g., for MiniMax proxy) |
| `ANTHROPIC_MODEL` | Default model to use |
| `API_TIMEOUT_MS` | HTTP request timeout in milliseconds |
| `DISABLE_TELEMETRY` | Disable anonymous usage telemetry (`1` to disable) |
| `LOG_LEVEL` | Set log verbosity (`debug`, `info`, `warn`, `error`) |
| `BASH_MAX_OUTPUT_LENGTH` | Limit shell output capture length |

### Reading Large Files

Claude Code can struggle with files over 2000 lines. Two strategies:

**1. Use offset + limit parameters:**

```
Read the file in chunks using offset=0, limit=500 then offset=500, limit=500, etc.
```

**2. Stash prompt trick:** start writing your prompt, then press `Ctrl+S` to stash it, then give another prompt — after pressing Enter, the stashed prompt returns.

---

## Plugins & Skills

### Marketplace Plugins

Install plugins via the marketplace:

```bash
/plugin marketplace add obra/superpowers-developing-for-claude-code
```

#### Recommended Plugins

| Plugin | Description |
|--------|-------------|
| **superpowers** | Smart planning, brainstorming, and context management |
| **code-simplifier** | Simplifies complex code automatically |
| **context7** | Provides up-to-date context about APIs, libraries, and services |
| **commit-commands** | Smart commit message generation |
| **code-review** | Automated code review |

#### Enabling Agent Teams (Experimental)

In `settings.json`:

```json
"CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
```

Then create teams of agents that can work together autonomously on complex tasks.

### Skills

Skills extend Claude Code's capabilities in specific domains.

#### Official Skills

| Resource | Link |
|----------|------|
| Claude Code Official Skills | `npx skills add @anthropics/skills` |
| Everything Claude Code | [affaan-m/everything-claude-code](https://github.com/affaan-m/everything-claude-code) |
| Composio Skills | [ComposioHQ/awesome-claude-skills](https://github.com/ComposioHQ/awesome-claude-skills) |
| Superpowers | [obra/superpowers](https://github.com/obra/superpowers) |

#### Remotion (Motion Graphics)

Generate videos from prompts using Remotion:

```bash
npx skills add remotion-dev/skills
mkdir v1 && cd v1
bun create video    # creates a Remotion project
# then prompt Claude Code to work with the Remotion project
```

📺 [Remotion x Claude Code Tutorial](https://www.youtube.com/watch?v=7OR-L0AySn8)
   https://www.youtube.com/shorts/sWQMAl8uf90
   https://www.youtube.com/watch?v=vc8rY3IJ5h4
---

## Workflows & Productivity

### Get Shit Done (GSD) — Meta-prompting Framework

A structured workflow for getting Claude Code to complete tasks efficiently:

```bash
npx get-shit-done-cc@latest
claude --dangerously-skip-permissions
```

📺 [GSD Framework](https://github.com/gsd-build/get-shit-done)

### Built-in Slash Commands

| Command | What it does |
|---------|-------------|
| `/insights` | Generates a usage report at `~/.claude/usage-data/reports.html` |
| `/effort` | Sets token budget: `auto`, `low`, `medium`, `high` |
| `/simplify` | Rewrites code in simpler, more readable form |

### Connecting CLI Tools

Claude Code can be chained with powerful CLI tools for browser automation and more:

📺 [Video: CLI tools that make Claude Code unstoppable](https://www.youtube.com/watch?v=uULvhQrKB_c)

| Tool | Purpose |
|------|---------|
| [Playwright CLI](https://github.com/microsoft/playwright-cli) | Browser automation |
| [Vercel Agent Browser](https://www.youtube.com/watch?v=P7JrP57AxR0) | AI-driven browser control |

---

## Status Line & Terminal

### Setting Up the Status Line

The status line shows live info at the bottom of your terminal. Install it with:

```bash
npx -y ccstatusline@latest
```

**Recommended 2-line configuration:**

**Line 1:**
```
Model | Context % | Session Cost | Session Clock
```

**Line 2:**
```
Git Branch | Git Worktree
```

### Terminal Recommendations

| Tool | Purpose |
|------|---------|
| **[WezTerm](https://wezterm.org)** | GPU-accelerated cross-platform terminal with navigation pane |
| **[Warp](https://warp.dev)** | Modern terminal with file navigation sidebar |
| **[Starship](https://starship.rs)** | Minimal, fast shell prompt |
| **[Nerd Fonts](https://www.nerdfonts.com)** | Icon fonts for powerline-style prompts |

📺 [Level up your terminal with WezTerm, Starship, and eza](https://marcopelxeluso.com/posts/level-up-your-macos-terminal-with-wezterm-starship-and-eza/)

---

## Advanced Features

### 21 Hidden Settings

Many Claude Code behaviors can be tuned. Check the video below for 21 hidden tricks:

📺 [21 Hidden Claude Code Settings](https://www.youtube.com/watch?v=pDoBe4qbFPE)

### Claude Code Self-Improvement

Claude Code can analyze its own past sessions and improve its behavior over time:

📺 [Video: Claude Code Self-Improve](https://www.youtube.com/watch?v=wQ0duoTeAAU)

### History Explorer

Review past Claude Code sessions in a visual GUI:

```bash
# Install
curl -fsSL https://raw.githubusercontent.com/jhlee0409/claude-code-history-viewer/main/install-server.sh | sh

# Dependencies (Ubuntu/Debian)
sudo apt-get install -y libwebkit2gtk-4.1-0

# Launch
cchv-server --serve
```

Then open `http://localhost:3727?token=<your-token>` in your browser.

---

## Remote Access

### Continue Work from Your Phone

Use [Happy Engineering](https://happy.engineering) to access your Claude Code session remotely from your mobile device.

---

## External Resources

### YouTube Playlists & Videos

| Topic | Link |
|-------|------|
| Claude Code Complete Guide | [Video](https://www.youtube.com/watch?v=TiNpzxoBPz0) |
| Claude Code Tips & Tricks | [Video](https://www.youtube.com/watch?v=-O6MEtleOdA) |
| Claude Code Workflows | [Video](https://www.youtube.com/watch?v=rVEoyx349Hk) |
| Claude Code Telegram Plugin | [Video](https://www.youtube.com/watch?v=W9igiY2JdHA) |
| Claude Code + OpenClaw | [LinkedIn Post](https://www.linkedin.com/feed/update/urn:li:activity:7446945416140619777/) |

### Community Projects

| Project | Description |
|---------|-------------|
| [serena](https://github.com/oraios/serena) | Search and indexing tool |
| [skillsmp.com](https://skillsmp.com) | Skill marketplace |
| [85danf/agent-skills](https://github.com/85danf/agent-skills) | Extensive CLAUDE.md collection |
| [claw-code](https://github.com/ultraworkers/claw-code) | Open agent harness |
| [oh-my-claudecode](https://github.com/Yeachan-Heo/oh-my-claudecode) | Claude Code enhancements |

---

<p align="center">
  <strong>Happy coding! 🚀</strong>
</p>
