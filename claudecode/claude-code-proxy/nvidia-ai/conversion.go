package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// ----------------------
// Conversion functions
// ----------------------

func convertAnthropicToOpenAI(req *anthropicMessageRequest) (openaiChatCompletionRequest, error) {
	var messages []any

	if sys := strings.TrimSpace(extractSystemText(req.System)); sys != "" {
		messages = append(messages, map[string]any{
			"role":    "system",
			"content": sys,
		})
	}

	for _, m := range req.Messages {
		role := strings.TrimSpace(m.Role)
		if role == "" {
			continue
		}

		// content can be string or array of blocks
		var asString string
		if err := json.Unmarshal(m.Content, &asString); err == nil {
			messages = append(messages, map[string]any{
				"role":    role,
				"content": asString,
			})
			continue
		}

		var blocks []anthropicContentBlock
		if err := json.Unmarshal(m.Content, &blocks); err != nil {
			return openaiChatCompletionRequest{}, fmt.Errorf("invalid message content for role %q", role)
		}

		switch role {
		case "user":
			userMsgs, err := convertAnthropicUserBlocksToOpenAIMessages(blocks)
			if err != nil {
				return openaiChatCompletionRequest{}, err
			}
			messages = append(messages, userMsgs...)
		case "assistant":
			assistantMsg, err := convertAnthropicAssistantBlocksToOpenAIMessage(blocks)
			if err != nil {
				return openaiChatCompletionRequest{}, err
			}
			messages = append(messages, assistantMsg)
		default:
			// pass through unknown roles as string if possible
			text := joinTextBlocks(blocks)
			messages = append(messages, map[string]any{
				"role":    role,
				"content": text,
			})
		}
	}

	out := openaiChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      req.Stream,
	}

	if len(req.Tools) > 0 {
		out.Tools = make([]any, 0, len(req.Tools))
		for _, t := range req.Tools {
			var params any
			if len(t.InputSchema) > 0 {
				_ = json.Unmarshal(t.InputSchema, &params)
			}
			out.Tools = append(out.Tools, map[string]any{
				"type": "function",
				"function": map[string]any{
					"name":        t.Name,
					"description": t.Description,
					"parameters":  params,
				},
			})
		}
	}

	if req.ToolChoice != nil {
		out.ToolChoice = convertToolChoice(req.ToolChoice)
	}

	return out, nil
}

func extractSystemText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var blocks []anthropicContentBlock
	if err := json.Unmarshal(raw, &blocks); err == nil {
		return joinTextBlocks(blocks)
	}
	return ""
}

func joinTextBlocks(blocks []anthropicContentBlock) string {
	var b strings.Builder
	for _, blk := range blocks {
		if blk.Type == "text" && blk.Text != "" {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString(blk.Text)
		}
	}
	return b.String()
}

func convertAnthropicUserBlocksToOpenAIMessages(blocks []anthropicContentBlock) ([]any, error) {
	var out []any

	// tool_result blocks become separate OpenAI "tool" messages.
	for _, blk := range blocks {
		if blk.Type != "tool_result" || strings.TrimSpace(blk.ToolUseID) == "" {
			continue
		}
		contentStr := ""
		if len(blk.Content) > 0 {
			var s string
			if err := json.Unmarshal(blk.Content, &s); err == nil {
				contentStr = s
			} else {
				contentStr = string(blk.Content)
			}
		}
		out = append(out, map[string]any{
			"role":         "tool",
			"tool_call_id": blk.ToolUseID,
			"content":      contentStr,
		})
	}

	// remaining text/image blocks become a user message
	var parts []any
	allText := true
	var allTextJoined strings.Builder
	for _, blk := range blocks {
		switch blk.Type {
		case "text":
			if blk.Text != "" {
				parts = append(parts, map[string]any{"type": "text", "text": blk.Text})
				if allTextJoined.Len() > 0 {
					allTextJoined.WriteString("\n")
				}
				allTextJoined.WriteString(blk.Text)
			}
		case "image":
			allText = false
			if blk.Source == nil {
				continue
			}
			url := ""
			switch blk.Source.Type {
			case "base64":
				if blk.Source.MediaType == "" || blk.Source.Data == "" {
					continue
				}
				// Validate base64 to avoid obviously invalid payloads.
				if _, err := base64.StdEncoding.DecodeString(blk.Source.Data); err != nil {
					continue
				}
				url = "data:" + blk.Source.MediaType + ";base64," + blk.Source.Data
			case "url":
				url = blk.Source.URL
			default:
				continue
			}
			if url != "" {
				parts = append(parts, map[string]any{
					"type": "image_url",
					"image_url": map[string]any{
						"url": url,
					},
				})
			}
		}
	}

	if len(parts) == 0 {
		out = append(out, map[string]any{"role": "user", "content": ""})
		return out, nil
	}
	if len(parts) == 1 {
		if p, ok := parts[0].(map[string]any); ok && p["type"] == "text" {
			if t, ok := p["text"].(string); ok {
				out = append(out, map[string]any{"role": "user", "content": t})
				return out, nil
			}
		}
	}

	// Many OpenAI-compatible providers (incl. some NVIDIA endpoints) require messages[].content to be a string.
	// If we only have text blocks, collapse them into a single string.
	if allText {
		out = append(out, map[string]any{
			"role":    "user",
			"content": allTextJoined.String(),
		})
		return out, nil
	}

	out = append(out, map[string]any{
		"role":    "user",
		"content": parts,
	})
	return out, nil
}

func convertAnthropicAssistantBlocksToOpenAIMessage(blocks []anthropicContentBlock) (any, error) {
	text := joinTextBlocks(blocks)

	var toolCalls []any
	for _, blk := range blocks {
		if blk.Type != "tool_use" || strings.TrimSpace(blk.ID) == "" || strings.TrimSpace(blk.Name) == "" {
			continue
		}
		args := "{}"
		if len(blk.Input) > 0 {
			args = string(blk.Input)
		}
		toolCalls = append(toolCalls, map[string]any{
			"id":   blk.ID,
			"type": "function",
			"function": map[string]any{
				"name":      blk.Name,
				"arguments": args,
			},
		})
	}

	msg := map[string]any{
		"role": "assistant",
	}
	if text != "" {
		msg["content"] = text
	} else {
		msg["content"] = nil
	}
	if len(toolCalls) > 0 {
		msg["tool_calls"] = toolCalls
	}
	return msg, nil
}

func convertToolChoice(v any) any {
	// Anthropic forms:
	// - {"type":"auto"}
	// - {"type":"tool","name":"my_tool"}
	// - string values (rare)
	m, ok := v.(map[string]any)
	if !ok {
		return v
	}
	typ, _ := m["type"].(string)
	switch typ {
	case "auto", "none", "required":
		return typ
	case "tool":
		name, _ := m["name"].(string)
		if name == "" {
			return "auto"
		}
		return map[string]any{
			"type": "function",
			"function": map[string]any{
				"name": name,
			},
		}
	default:
		return v
	}
}

// ----------------------
// OpenAI response conversion
// ----------------------

func convertOpenAIToAnthropic(resp openaiChatCompletionResponse) anthropicMessageResponse {
	content := make([]any, 0, 4)

	var finishReason string
	if len(resp.Choices) > 0 {
		ch := resp.Choices[0]
		finishReason = ch.FinishReason
		if ch.Message.Content != nil && *ch.Message.Content != "" {
			content = append(content, map[string]any{
				"type": "text",
				"text": *ch.Message.Content,
			})
		}
		if len(ch.Message.ToolCalls) > 0 {
			for _, tc := range ch.Message.ToolCalls {
				input := map[string]any{}
				switch v := tc.Function.Arguments.(type) {
				case string:
					_ = json.Unmarshal([]byte(v), &input)
				case map[string]any:
					input = v
				default:
					input = map[string]any{"text": fmt.Sprintf("%v", v)}
				}
				content = append(content, map[string]any{
					"type":  "tool_use",
					"id":    tc.ID,
					"name":  tc.Function.Name,
					"input": input,
				})
			}
		}
	}

	inputTokens := 0
	outputTokens := 0
	cacheRead := 0
	if resp.Usage != nil {
		cacheRead = 0
		if resp.Usage.PromptTokensDetails != nil {
			cacheRead = resp.Usage.PromptTokensDetails.CachedTokens
		}
		inputTokens = resp.Usage.PromptTokens - cacheRead
		outputTokens = resp.Usage.CompletionTokens
	}

	return anthropicMessageResponse{
		ID:    resp.ID,
		Type:  "message",
		Role:  "assistant",
		Model: resp.Model,
		Content: content,
		StopReason: mapFinishReason(finishReason),
		StopSequence: nil,
		Usage: map[string]any{
			"input_tokens":           inputTokens,
			"output_tokens":          outputTokens,
			"cache_read_input_tokens": cacheRead,
		},
	}
}

func mapFinishReason(finish string) string {
	switch finish {
	case "stop":
		return "end_turn"
	case "length":
		return "max_tokens"
	case "tool_calls":
		return "tool_use"
	case "content_filter":
		return "stop_sequence"
	default:
		if finish == "" {
			return "end_turn"
		}
		return "end_turn"
	}
}

// ----------------------
// Model helpers
// ----------------------

func modelSupportsToolUse(model string) bool {
	m := strings.ToLower(strings.TrimSpace(model))
	if m == "" {
		return true
	}
	// NVIDIA free Qwen endpoint commonly rejects tool calling.
	if m == "qwen/qwen2.5-coder-32b-instruct" {
		return false
	}
	return true
}

func clampMaxTokens(model string, requested int) int {
	if requested <= 0 {
		return requested
	}
	m := strings.ToLower(strings.TrimSpace(model))
	if m == "qwen/qwen2.5-coder-32b-instruct" {
		cap := 8192
		if requested > cap {
			return cap
		}
	}
	return requested
}
