# Claude NVIDIA Proxy

A proxy service that allows Claude Code to use NVIDIA's AI API for completions.

## Prerequisites

- Go 1.24+ (for local development)
- Docker 20.10+ (for containerized deployment)
- NVIDIA API key (get one at https://ngc.nvidia.com/)

## Quick Start with Docker

```bash
# 1. Clone and navigate to the project
cd claude-code-proxy/nvidia-ai

# 2. Create your config file
cp .env.example .env
# Edit .env and add your NVIDIA API key

# 3. Build the Docker image
./dockerbuild

# 4. Run the container
./dockerrun
```

## Docker Compose (Recommended)

```bash
# Set your API key
export NVIDIA_API_KEY=your-nvidia-api-key

# Start the service
docker-compose up -d

# Check logs
docker-compose logs -f
```

## Configuration

The proxy can be configured via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `ADDR` | Server listen address | `:3001` |
| `CONFIG_PATH` | Path to config.json | `config.json` |
| `UPSTREAM_URL` | NVIDIA API endpoint | (from config.json) |
| `PROVIDER_API_KEY` | NVIDIA API key | (from config.json) |
| `SERVER_API_KEY` | Inbound auth key (optional) | (empty) |
| `UPSTREAM_TIMEOUT_SECONDS` | Request timeout | `300` |

## Kubernetes Deployment

```bash
kubectl apply -f k8s/deployment.yaml
```

Update the `nvidia-proxy-secrets` secret with your actual API key.

## Local Development

```bash
# Install dependencies
go mod download

# Run directly
CONFIG_PATH=config.json PROVIDER_API_KEY=your-key go run .
```

## CI/CD

GitHub Actions workflows are provided for:
- Go build and test
- Docker image build and push to GHCR
- Security scanning with Trivy

See `.github/workflows/ci.yml` for details.
