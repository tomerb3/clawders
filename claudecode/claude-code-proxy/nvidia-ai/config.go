package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ----------------------
// Configuration
// ----------------------

type fileConfig struct {
	NvidiaURL string `json:"nvidia_url"`
	NvidiaKey string `json:"nvidia_key"`
}

type serverConfig struct {
	addr           string
	upstreamURL    string
	providerAPIKey string
	serverAPIKey   string
	timeout        time.Duration
	logBodyMax     int
	logStreamPreviewMax int
}

func loadConfig() (*serverConfig, error) {
	fc, err := loadFileConfig(strings.TrimSpace(envOr("CONFIG_PATH", "config.json")))
	if err != nil {
		return nil, err
	}

	addr := strings.TrimSpace(envOr("ADDR", ":3001"))
	upstreamURL := strings.TrimSpace(envOr("UPSTREAM_URL", fc.NvidiaURL))
	providerAPIKey := strings.TrimSpace(envOr("PROVIDER_API_KEY", fc.NvidiaKey))
	serverAPIKey := strings.TrimSpace(envOr("SERVER_API_KEY", ""))

	timeout := 5 * time.Minute
	if raw := strings.TrimSpace(envOr("UPSTREAM_TIMEOUT_SECONDS", "")); raw != "" {
		seconds, err := strconv.Atoi(raw)
		if err != nil || seconds <= 0 {
			return nil, fmt.Errorf("invalid UPSTREAM_TIMEOUT_SECONDS: %q", raw)
		}
		timeout = time.Duration(seconds) * time.Second
	}

	logBodyMax := 4096
	if raw := strings.TrimSpace(envOr("LOG_BODY_MAX_CHARS", "")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("invalid LOG_BODY_MAX_CHARS: %q", raw)
		}
		logBodyMax = n
	}

	logStreamPreviewMax := 256
	if raw := strings.TrimSpace(envOr("LOG_STREAM_TEXT_PREVIEW_CHARS", "")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("invalid LOG_STREAM_TEXT_PREVIEW_CHARS: %q", raw)
		}
		logStreamPreviewMax = n
	}

	if upstreamURL == "" {
		return nil, errors.New("missing nvidia_url in config.json (or UPSTREAM_URL)")
	}
	if providerAPIKey == "" {
		return nil, errors.New("missing nvidia_key in config.json (or PROVIDER_API_KEY)")
	}
	return &serverConfig{
		addr:           addr,
		upstreamURL:    upstreamURL,
		providerAPIKey: providerAPIKey,
		serverAPIKey:   serverAPIKey,
		timeout:        timeout,
		logBodyMax:     logBodyMax,
		logStreamPreviewMax: logStreamPreviewMax,
	}, nil
}

func loadFileConfig(path string) (*fileConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var fc fileConfig
	if err := json.Unmarshal(b, &fc); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &fc, nil
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
