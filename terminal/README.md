# 💻 Terminal Setup







A guide to building a powerful, visually appealing terminal workflow for use with Claude Code and AI agents.

---

## Recommended Tools

| Tool | Purpose | Link |
|------|---------|------|
| **WezTerm** | GPU-accelerated cross-platform terminal emulator | [wezterm.org](https://wezterm.org) |
     inside WezTerm I use Herdr ( AI Tmux for Agents)  https://herdr.dev/ https://www.youtube.com/watch?v=XoitaexiCi0
| **Starship** | Minimal, blazing-fast shell prompt | [starship.rs](https://starship.rs) |
| **Nerd Fonts** | Icon fonts for powerline-style prompts | [nerdfonts.com](https://www.nerdfonts.com) |
| **eza** | Modern, colorful `ls` replacement | [eza.rocks](https://eza.rocks) |

---

## Quick Install

```bash
# WezTerm (macOS/Linux)
brew install --cask wezterm
# or download from https://wezterm.org

# Starship
curl -fsSL https://starship.rs/install.sh | sh

# eza
brew install eza

# Nerd Fonts — download and install from https://www.nerdfonts.com/
```

---

## Configuration

### WezTerm + Starship

📺 [Level up your terminal with WezTerm, Starship, and eza](https://www.youtube.com/watch?v=LnZdaNfQ86o)
📝 [Blog: macOS Terminal Setup Guide](https://marcopeluso.com/posts/level-up-your-macos-terminal-with-wezterm-starship-and-eza/)

### Quick Config Snippet

Add to your `~/.zshrc`:

```bash
# Starship prompt
eval "$(starship init zsh)"

# eza aliases
alias ls="eza --icons"
alias ll="eza -l --icons"
alias la="eza -la --icons"
alias tree="eza --tree --icons"
```

---

## Why These Tools?

- **WezTerm** — GPU-rendered text, cross-platform, configurable Lua keybindings, and a navigation pane for files
- **Starship** — zero-configuration prompt that shows git info, node version, and more out of the box
- **eza** — fast, modern replacement for `ls` with icons, git status, and color support
- **Nerd Fonts** — fixes the icon rendering issue when using powerline fonts

This stack pairs perfectly with Claude Code to give you a fast, readable, information-rich terminal experience.
