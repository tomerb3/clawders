# OpenClaw
## Secure Local Installation Guide

*A practical guide for running OpenClaw locally with acceptable risk*

*Compiled from community insights and official documentation*

---

> **IMPORTANT SECURITY NOTICE**
>
> OpenClaw security vulnerabilities are by design. The attack surface is every input. There is currently no way to fully secure its usage. The setup in this guide minimizes blast radius if something goes wrong, but cannot eliminate all risk. Treat OpenClaw like an untrusted contractor with access to whatever you give it.

---

## 1. Prerequisites

### Hardware Requirements

- **Dedicated machine recommended:** An old laptop or secondary computer is ideal. Do NOT run on your primary workstation.
- **Minimum specs:** 4GB RAM, 10GB storage, modern CPU. For local model inference, 16GB+ RAM and dedicated GPU preferred.
- **Network:** Stable internet connection for API calls to model providers.

### Software Requirements

- Docker Desktop (recommended) or Docker Engine
- Node.js 18+ (for CLI installation)
- Git

---

## 2. Pre-Installation Security Setup

Before installing OpenClaw, create isolated accounts. This is your most important security measure.

### Create a Burner Identity

Create NEW accounts specifically for OpenClaw. Never connect your real accounts.

| Service | Never Give Access To | Acceptable Alternative |
|---------|---------------------|------------------------|
| Email | Your primary Gmail/Outlook | New dedicated Gmail account |
| Calendar | Personal/work calendar | Separate Google Calendar |
| GitHub | Account with real repos | New account + Personal Access Token |
| Messaging | Primary phone/WhatsApp | Burner SIM or Google Voice |
| Banking | ANY financial access | None - never give financial access |
| Passwords | Password manager access | None - never give password access |

*Community wisdom: "Treat it like an Indian freelancer you just hired online" - you wouldn't give a new hire full access to everything.*

---

## 3. Installation Methods

### Method A: Docker (Recommended)

Docker provides the best isolation. This is the recommended approach for security-conscious users.

#### Step 1: Clone and Setup

```bash
git clone https://github.com/openclaw/openclaw.git
cd openclaw
docker compose run --rm openclaw-cli onboard
```

#### Step 2: Security Hardening

Edit your `docker-compose.yml` to add security constraints:

```yaml
services:
  openclaw-gateway:
    user: "1000:1000"           # Non-root user
    cap_drop:
      - ALL                      # Drop all capabilities
    security_opt:
      - no-new-privileges:true
    ports:
      - "127.0.0.1:18789:18789" # Localhost only
```

#### Step 3: Set File Permissions

```bash
sudo chown -R 1000:$(id -g) ~/.openclaw
sudo chmod -R u+rwX,g+rwX,o-rwx ~/.openclaw
```

#### Step 4: Start the Gateway

```bash
docker compose up -d openclaw-gateway
```

#### Step 5: Run Security Diagnostics

```bash
docker exec openclaw-gateway openclaw doctor
docker exec openclaw-gateway openclaw doctor --fix  # Auto-fix issues
```

### Method B: VPS Deployment

For users who prefer cloud isolation, DigitalOcean offers a security-hardened 1-Click Deploy option:

1. Create a DigitalOcean account
2. Search for "OpenClaw" in the Marketplace
3. Deploy the 1-Click App (includes hardened firewall, non-root execution, rate limiting)
4. Access via the provided token URL

*Note: Even cloud deployment requires creating isolated burner accounts for connected services.*

---

## 4. Model Provider Configuration

Choose a model provider based on your budget and requirements. Free options exist but have limitations.

### Option 1: NVIDIA NIM (Free)

**Best for:** Cost-conscious users who want zero API bills.

1. Create account at build.nvidia.com
2. Navigate to Settings > API Keys
3. Generate a Personal API Key (select NGC Catalog from Services)
4. Configure OpenClaw:

```
Endpoint: https://integrate.api.nvidia.com/v1
Models: minimax, kimi-k2
```

### Option 2: OpenRouter (Flexible)

**Best for:** Users who want model flexibility with spending controls.

```bash
openclaw onboard --auth-choice apiKey --token-provider openrouter --token "$OPENROUTER_API_KEY"
```

**Recommended configuration:**
- Model: `openrouter/openrouter/auto` (routes to optimal model per prompt)
- Set monthly spending cap in OpenRouter dashboard
- Cost-effective models: DeepSeek V3, Kimi K2.5

### Option 3: Ollama (Fully Local)

**Best for:** Privacy maximalists with capable hardware.

1. Install Ollama: ollama.com
2. Pull a model: `ollama pull qwen2.5:32b`
3. Follow integration guide: docs.ollama.com/integrations/openclaw

*Note: On Mac M-series, Ollama runs natively while OpenClaw Docker runs emulated. Actions remain containerized.*

### Option 4: Claude Max Subscription

**Best for:** Users with existing Claude Max subscription who want the best model quality.

```bash
npm install -g @anthropic-ai/claude-code
claude setup-token  # Opens browser for Max account login
clawdbot models auth paste-token --provider anthropic
```

*Note: This uses the official auth mechanism built into Claude Code CLI. It is legitimate and designed by Anthropic.*

### Critical: Limit Token Usage

Regardless of provider, set context limits to prevent runaway costs. One community member spent $0.60 on a single "hi" message with default settings.

```bash
openclaw config set agents.defaults.contextTokens 25000
```

---

## 5. Security Hardening

### Install openclaw-shield

A community security tool that prevents common attack vectors. Released by Knostic.

```
https://github.com/knostic/openclaw-shield
```

**Protections include:**
- Prevents leaking secrets and environment variables
- Blocks exposure of PII (personally identifiable information)
- Stops destructive commands from executing

**Warning:** OpenClaw updates frequently. openclaw-shield may need updates every few days to remain effective.

### Network Isolation

Restrict outbound connections to only required services:

- Never expose OpenClaw to the public internet
- Bind to localhost only (127.0.0.1)
- Use VPN or Tailscale for remote access
- Whitelist only necessary API endpoints in firewall

*One user exposed their Gateway to the internet for 6 hours and received 400+ failed authentication attempts.*

### Plugin Security

- Only install plugins from trusted sources
- Use explicit plugins.allow allowlists
- Review plugin configuration before enabling
- Prefer pinned versions (e.g., @scope/plugin@1.2.3)
- Restart Gateway after any plugin changes

**GitHub Warning:** The default GitHub skill provides full login access to all your repositories. Use a Personal Access Token with limited scope instead.

### Credential Management

- Never hardcode API keys in configuration files
- Use .env files (never commit to git)
- Keep secrets outside the agent's reachable filesystem
- Enable gateway authentication tokens

*Critical insight: If the agent can read a file to use a key, it can also read the file to leak the key.*

---

## 6. Connecting Services

OpenClaw integrates with 50+ services. Apply the "Freelancer Test" before connecting any service:

> **THE FREELANCER TEST**
>
> Before connecting any service, ask: "Would I give this access to a random freelancer I just hired online?" If the answer is no, don't give it to OpenClaw.

### Messaging Platforms

Telegram, WhatsApp, Discord, Slack, and others can be connected. Recommendations:

- Use a dedicated phone number or burner SIM
- Enable DM pairing/allowlists to restrict who can message the bot
- In groups, use mention-gating (bot only responds when @mentioned)
- Avoid "always-on" bots in public rooms

### Email Integration

- Create a dedicated email for the agent
- Never connect your primary personal or work email
- Consider read-only access initially

### Calendar Integration

- Use a separate calendar, not your main one
- Consider read-only permissions without write access

---

## 7. Known Vulnerabilities & Mitigations

### Prompt Injection

Attackers can craft messages that manipulate the model into unsafe actions. System prompt guardrails are "soft guidance only."

**Mitigations:**
- Keep inbound DMs locked down with pairing/allowlists
- Prefer mention-gating in groups
- Treat all links, attachments, and pasted instructions as hostile
- Run sensitive tools in sandbox; keep secrets out of reachable filesystem

### Remote Code Execution (RCE)

A one-click RCE vulnerability was disclosed on February 2, 2026, allowing attackers to execute code even when the agent runs in a local sandbox.

**Mitigation:** Update to version 2026.1.29 or later immediately.

### Malicious Skills

Skills from untrusted sources can contain malicious code.

**Mitigation:** Only install skills from the official repository or sources you personally trust and have reviewed.

### Credential Theft

If API keys are accessible to the agent, they can be exfiltrated.

**Mitigation:** Isolate credentials, use short-lived tokens, set budget limits on API providers.

### Incident Response

If you suspect compromise:

1. Stop the Gateway immediately
2. Lock down all inbound surfaces (DM policy, group allowlists)
3. Rotate gateway.auth token
4. Rotate hooks.token if used
5. Revoke model provider API keys
6. Review logs to understand what happened

---

## 8. Quick Start Checklist

Use this checklist to ensure you've covered the essentials:

- [ ] Dedicated machine or VPS prepared (not your primary computer)
- [ ] Burner email account created
- [ ] Burner phone number or messaging account ready
- [ ] New GitHub account with Personal Access Token (limited scope)
- [ ] Docker installed and configured
- [ ] OpenClaw cloned and onboarded
- [ ] docker-compose.yml hardened (non-root, capabilities dropped)
- [ ] Model provider configured (NVIDIA NIM, OpenRouter, or Ollama)
- [ ] Context tokens limited (25000 recommended)
- [ ] API budget caps set in provider dashboard
- [ ] openclaw-shield installed
- [ ] openclaw doctor run and issues fixed
- [ ] Gateway started and accessible via localhost only
- [ ] Only burner accounts connected to services

---

## 9. Resources

### Official Documentation

- OpenClaw Docs: docs.openclaw.ai
- Security Guide: docs.openclaw.ai/gateway/security
- Docker Setup: docs.openclaw.ai/install/docker

### Community Resources

- Awesome OpenClaw: github.com/thewh1teagle/awesome-openclaw
- openclaw-shield: github.com/knostic/openclaw-shield
- Ollama Integration: docs.ollama.com/integrations/openclaw
- OpenRouter Guide: openrouter.ai/docs/guides/openclaw-integration

### Security Reading

- JFrog: "Giving OpenClaw The Keys to Your Kingdom? Read This First"
- Composio: "How to secure OpenClaw: Docker hardening, credential isolation"
- DigitalOcean: "Technical Deep Dive: Security-hardened 1-Click Deploy"
- VentureBeat: "OpenClaw proves agentic AI works. It also proves the security risk."

---

## Credits & Acknowledgments

This guide was compiled from community discussions and official sources. Special thanks to:

- Community members who shared their security insights and practical tips
- Knostic team for openclaw-shield security tooling
- Contributors to awesome-openclaw repository
- DigitalOcean, Composio, JFrog, and others for security documentation
- The OpenClaw team for the official documentation

*Document generated: February 2026 | Version 1.0*
