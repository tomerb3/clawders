# Build stage
FROM docker.io/library/golang:1.24-alpine AS build

WORKDIR /src

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Build with explicit architecture for multi-platform support
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/claude-nvidia-proxy ./

# Runtime stage
FROM docker.io/library/alpine:3.20

LABEL org.opencontainers.image.title="Claude NVIDIA Proxy"
LABEL org.opencontainers.image.description="Proxy service for Claude Code to use NVIDIA AI"
LABEL org.opencontainers.image.source="https://github.com/zhangrr/claude-nvidia-proxy"
LABEL org.opencontainers.image.licenses="MIT"

RUN adduser -D -H -u 10001 app

WORKDIR /app

# Install CA certificates for HTTPS calls
RUN apk add --no-cache ca-certificates

COPY --from=build /out/claude-nvidia-proxy /app/claude-nvidia-proxy

USER app

EXPOSE 3001

ENV ADDR=":3001"
ENV CONFIG_PATH="/app/config.json"

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:3001/ || exit 1

ENTRYPOINT ["/app/claude-nvidia-proxy"]
