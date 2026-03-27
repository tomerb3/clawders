package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// ----------------------
// Logging
// ----------------------

func logForwardedRequest(reqID string, cfg *serverConfig, anthropicReq anthropicMessageRequest, openaiReq openaiChatCompletionRequest) {
	inSummary := map[string]any{
		"model":      anthropicReq.Model,
		"max_tokens": anthropicReq.MaxTokens,
		"stream":     anthropicReq.Stream,
		"messages":   len(anthropicReq.Messages),
		"tools":      len(anthropicReq.Tools),
	}
	log.Printf("[%s] inbound summary=%s", reqID, mustJSONTrunc(inSummary, cfg.logBodyMax))

	out := sanitizeOpenAIRequest(openaiReq)
	log.Printf("[%s] forward url=%s", reqID, cfg.upstreamURL)
	log.Printf("[%s] forward headers=%s", reqID, mustJSONTrunc(map[string]any{
		"Content-Type":  "application/json",
		"Authorization": "Bearer <redacted>",
	}, cfg.logBodyMax))
	log.Printf("[%s] forward body=%s", reqID, mustJSONTrunc(out, cfg.logBodyMax))
}

func logForwardedUpstreamBody(reqID string, cfg *serverConfig, body []byte) {
	if cfg.logBodyMax == 0 {
		return
	}
	s := string(body)
	if len([]rune(s)) > cfg.logBodyMax {
		s = string([]rune(s)[:cfg.logBodyMax]) + "...(truncated)"
	}
	log.Printf("[%s] upstream body=%s", reqID, s)
}

func mustJSONTrunc(v any, maxChars int) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"_error":"json_marshal_failed","detail":%q}`, err.Error())
	}
	s := string(b)
	if maxChars == 0 {
		return "(disabled)"
	}
	if len([]rune(s)) > maxChars {
		return string([]rune(s)[:maxChars]) + "...(truncated)"
	}
	return s
}

// ----------------------
// Sanitization
// ----------------------

func sanitizeOpenAIRequest(req openaiChatCompletionRequest) openaiChatCompletionRequest {
	out := req
	out.Messages = sanitizeOpenAIMessages(req.Messages)
	out.Tools = sanitizeAnySlice(req.Tools)
	return out
}

func sanitizeOpenAIMessages(msgs []any) []any {
	if len(msgs) == 0 {
		return nil
	}
	out := make([]any, 0, len(msgs))
	for _, m := range msgs {
		mm, ok := m.(map[string]any)
		if !ok {
			out = append(out, m)
			continue
		}
		cp := map[string]any{}
		for k, v := range mm {
			cp[k] = v
		}
		if content, ok := cp["content"]; ok {
			cp["content"] = sanitizeMessageContent(content)
		}
		// tool_calls may carry huge arguments; keep but truncate strings.
		if tc, ok := cp["tool_calls"]; ok {
			cp["tool_calls"] = sanitizeAny(tc)
		}
		out = append(out, cp)
	}
	return out
}

func sanitizeMessageContent(content any) any {
	switch v := content.(type) {
	case string:
		return v
	case []any:
		parts := make([]any, 0, len(v))
		for _, p := range v {
			pm, ok := p.(map[string]any)
			if !ok {
				parts = append(parts, p)
				continue
			}
			cp := map[string]any{}
			for k, vv := range pm {
				cp[k] = vv
			}
			if cp["type"] == "image_url" {
				if iu, ok := cp["image_url"].(map[string]any); ok {
					if url, ok := iu["url"].(string); ok && strings.HasPrefix(url, "data:") {
						iu2 := map[string]any{}
						for k, vv := range iu {
							iu2[k] = vv
						}
						iu2["url"] = "data:<redacted>"
						cp["image_url"] = iu2
					}
				}
			}
			parts = append(parts, cp)
		}
		return parts
	default:
		return sanitizeAny(v)
	}
}

func sanitizeAnySlice(v []any) []any {
	if len(v) == 0 {
		return nil
	}
	out := make([]any, 0, len(v))
	for _, it := range v {
		out = append(out, sanitizeAny(it))
	}
	return out
}

func sanitizeAny(v any) any {
	switch t := v.(type) {
	case map[string]any:
		cp := map[string]any{}
		for k, vv := range t {
			cp[k] = sanitizeAny(vv)
		}
		return cp
	case []any:
		return sanitizeAnySlice(t)
	case string:
		// keep strings; truncation is handled at final JSON layer
		return t
	default:
		return v
	}
}

func takeFirstRunes(s string, max int) string {
	if max <= 0 || s == "" {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}
