# 🦞 Clawders

**The community guide to running OpenClaw securely, affordably, and efficiently**

[![OpenClaw](https://img.shields.io/badge/OpenClaw-2026.1.30+-orange)](https://openclaw.ai)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/tomerb3/clawders/pulls)

[![Star History Chart](https://api.star-history.com/svg?repos=tomerb3/clawders&type=date&legend=top-left)](https://www.star-history.com/#tomerb3/clawders&type=date&legend=top-left)

---

## 🎯 What is Clawders ?

**Clawders** is a community-maintained guide for installing and running [OpenClaw](https://openclaw.ai) — securely, affordably, and with minimal friction.

OpenClaw is a powerful AI agent framework, but its attack surface is large and its security posture requires careful attention. This repo aggregates real-world community experience, official documentation, and security research so you don't have to learn things the hard way.

This guide helps you:

- 🔒 **Minimize risk** — proper isolation, burner accounts, and hardened configurations
- 👤 **Get started easily** — step-by-step instructions from basic to advanced
- 💰 **Save money** — free model providers and token-efficient settings

---

## ⚠️ Critical Security Warning

> **OpenClaw security vulnerabilities are by design.** The attack surface is every input. There is currently no way to fully secure its usage.

These guides help **minimize the blast radius** if something goes wrong. The golden rule:

> **Treat OpenClaw like an untrusted contractor with access to whatever you give it.**

---

## 📁 Repository Contents

| Path | Description |
|------|-------------|
| [`openclaw/`](openclaw/) | **OpenClaw guides** — installation, security hardening, and tips |
| [`claudecode/`](claudecode/) | **Claude Code guides** — configuration, plugins, skills, and workflows |
| [`claw-code/`](claw-code/) | **Claw Code** — an open agent harness built on OpenClaw |
| [`terminal/`](terminal/) | **Terminal setup** — Wezterm, Starship, and CLI tools |

---

## 🚀 Quick Start: Install OpenClaw (Free, in 30 seconds)

The fastest way to get started with OpenClaw using the **MiniMax M2.1** model — completely free:

```bash
curl -fsSL skyler-agent.github.io/oclaw/i.sh | bash
```

This one-command installer:
- ✅ Automatically configures **MiniMax M2.1** (free tier)
- ✅ No manual API key setup required
- ✅ Includes optimized "7-day Coding Plan" presets
- ✅ Works out of the box on Linux/macOS

> 📺 [See it in action (Twitter/X)](https://x.com/SkylerMiao7/status/2017789329138212986) — 73K+ views

---

## 🔒 Essential Security Checklist

**Before you install OpenClaw, you MUST do the following:**

| Step | Why It Matters |
|------|----------------|
| ✅ Use a **dedicated machine** | If compromised, only that machine is affected |
| ✅ Create a **burner email** | Don't expose your real inbox |
| ✅ Create a **new GitHub account** | Use PATs with limited scope |
| ✅ Use a **burner phone/SIM** | For WhatsApp/Telegram integration |
| ❌ Never connect **primary email** | Full inbox access = full compromise |
| ❌ Never connect **banking/financial** | No exceptions, ever |
| ❌ Never connect **password managers** | Would expose all your credentials |

### The Freelancer Test™

> Before connecting **any** service, ask yourself: *"Would I give this access to a random freelancer I just hired online?"*
>
> **If the answer is NO → Don't give it to OpenClaw.**

---

## 💰 Free & Token-Efficient Model Options

A community member once spent **$0.60 just saying "hi"** with default Anthropic settings. Don't be that person.

### Free Options

| Provider | Model | Notes |
|----------|-------|-------|
| **MiniMax M2.1** | MiniMax M2.1 | Use the one-command installer above |
| **NVIDIA NIM** | MiniMax, Kimi K2.5 | Free tier at [build.nvidia.com](https://build.nvidia.com) |
| **Ollama** | qwen2.5, llama3, etc. | 100% local — [ollama.com](https://ollama.com) |

### Token Usage Limit

Set a hard limit on context tokens to prevent runaway costs:

```bash
openclaw config set agents.defaults.contextTokens 25000
```

---

## 🛡️ Security Tools & Resources

- **[openclaw-shield](https://github.com/knostic/openclaw-shield)** — Blocks secret/PII leakage and destructive commands
- **[awesome-openclaw](https://github.com/thewh1teagle/awesome-openclaw)** — Curated list of tools, guides, and integrations

---

## 🤝 Contributing

Found a tip that saved you hours? A security practice worth sharing? PRs are welcome!

---

## 📜 Credits

This guide was compiled from community wisdom, official OpenClaw docs, and security research from JFrog, Composio, DigitalOcean, and VentureBeat.

Special thanks to:
- **@steipete** — OpenClaw founder
- **@SkylerMiao7** — one-command installer author
- **Knostic team** — openclaw-shield security tooling

https://www.youtube.com/watch?v=TTgQV21X0SQ
  https://www.josean.com/posts/how-to-setup-wezterm-terminal


---

<p align="center">
  <strong>Stay secure. Stay efficient. Stay clawed. 🦞</strong>
</p>
