package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ----------------------
// HTTP Handlers
// ----------------------

func handleMessages(w http.ResponseWriter, r *http.Request, cfg *serverConfig) {
	reqID := fmt.Sprintf("req_%d", time.Now().UnixNano())
	if cfg.serverAPIKey != "" && !checkInboundAuth(r, cfg.serverAPIKey) {
		log.Printf("[%s] inbound unauthorized", reqID)
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var anthropicReq anthropicMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&anthropicReq); err != nil {
		log.Printf("[%s] invalid inbound json: %v", reqID, err)
		writeJSONError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if strings.TrimSpace(anthropicReq.Model) == "" {
		log.Printf("[%s] missing model", reqID)
		writeJSONError(w, http.StatusBadRequest, "missing_model")
		return
	}
	if anthropicReq.MaxTokens == 0 {
		// Anthropic requires max_tokens; NVIDIA/OpenAI also expects it. Default conservatively.
		anthropicReq.MaxTokens = 1024
	}

	openaiReq, err := convertAnthropicToOpenAI(&anthropicReq)
	if err != nil {
		log.Printf("[%s] request conversion failed: %v", reqID, err)
		writeJSONError(w, http.StatusBadRequest, "request_conversion_failed")
		return
	}

	// Some NVIDIA/OpenAI-compatible models reject tool-calling fields entirely.
	// Claude Code frequently sends tools/tool_choice even for simple prompts.
	if !modelSupportsToolUse(openaiReq.Model) {
		openaiReq.Tools = nil
		openaiReq.ToolChoice = nil
	}

	// Claude Code may request very large completions (e.g. 32000). Some models have smaller context windows.
	// We don't tokenize here, so cap conservatively for known-limited models to avoid upstream 400s.
	openaiReq.MaxTokens = clampMaxTokens(openaiReq.Model, openaiReq.MaxTokens)

	logForwardedRequest(reqID, cfg, anthropicReq, openaiReq)

	if anthropicReq.Stream {
		if err := proxyStream(w, r, cfg, reqID, openaiReq); err != nil {
			log.Printf("[%s] stream proxy error: %v", reqID, err)
		}
		return
	}

	openaiRespBody, resp, err := doUpstreamJSON(r.Context(), cfg, openaiReq)
	if err != nil {
		log.Printf("[%s] upstream request failed: %v", reqID, err)
		writeJSONError(w, http.StatusBadGateway, "upstream_request_failed")
		return
	}
	defer resp.Body.Close()
	log.Printf("[%s] upstream status=%d", reqID, resp.StatusCode)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(openaiRespBody)
		logForwardedUpstreamBody(reqID, cfg, openaiRespBody)
		return
	}

	var openaiResp openaiChatCompletionResponse
	if err := json.Unmarshal(openaiRespBody, &openaiResp); err != nil {
		log.Printf("[%s] invalid upstream json: %v", reqID, err)
		logForwardedUpstreamBody(reqID, cfg, openaiRespBody)
		writeJSONError(w, http.StatusBadGateway, "invalid_upstream_json")
		return
	}
	anthropicResp := convertOpenAIToAnthropic(openaiResp)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(anthropicResp)
}

// ----------------------
// Authentication
// ----------------------

func checkInboundAuth(r *http.Request, expected string) bool {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		got := strings.TrimSpace(auth[len("bearer "):])
		return subtleConstantTimeCompare(got, expected) == 1
	}
	if got := strings.TrimSpace(r.Header.Get("x-api-key")); got != "" {
		return subtleConstantTimeCompare(got, expected) == 1
	}
	return false
}

func subtleConstantTimeCompare(a, b string) int {
	if len(a) != len(b) {
		return 0
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	if result == 0 {
		return 1
	}
	return 0
}

// ----------------------
// Upstream proxy
// ----------------------

func doUpstreamJSON(ctx context.Context, cfg *serverConfig, openaiReq openaiChatCompletionRequest) ([]byte, *http.Response, error) {
	bodyBytes, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.upstreamURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.providerAPIKey)

	client := &http.Client{Timeout: cfg.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return nil, nil, err
	}
	_ = resp.Body.Close()
	// Re-wrap body so caller can optionally read again after status checks.
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	return respBody, resp, nil
}

// ----------------------
// Error handling
// ----------------------

func writeJSONError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"type":    "proxy_error",
			"code":    code,
			"message": code,
		},
	})
}
