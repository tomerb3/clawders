package main

import "encoding/json"

// ----------------------
// Anthropic request types
// ----------------------

type anthropicMessageRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature *float64        `json:"temperature,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
	System      json.RawMessage `json:"system,omitempty"`
	Messages    []anthropicMsg  `json:"messages"`
	Tools       []anthropicTool `json:"tools,omitempty"`
	ToolChoice  any             `json:"tool_choice,omitempty"`
	Thinking    any             `json:"thinking,omitempty"`
}

type anthropicMsg struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type anthropicTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
}

type anthropicContentBlock struct {
	Type string `json:"type"`

	// text
	Text string `json:"text,omitempty"`

	// image
	Source *anthropicImageSource `json:"source,omitempty"`

	// tool_use
	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`

	// tool_result
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
}

type anthropicImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data,omitempty"`
	URL       string `json:"url,omitempty"`
}

// ----------------------
// OpenAI request types
// ----------------------

type openaiChatCompletionRequest struct {
	Model       string `json:"model"`
	Messages    []any  `json:"messages"`
	MaxTokens   int    `json:"max_tokens,omitempty"`
	Temperature any    `json:"temperature,omitempty"`
	Stream      bool   `json:"stream,omitempty"`
	Tools       []any  `json:"tools,omitempty"`
	ToolChoice  any    `json:"tool_choice,omitempty"`
}

type openaiChatCompletionResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role      string `json:"role"`
			Content   *string `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments any    `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		PromptTokensDetails *struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details,omitempty"`
	} `json:"usage,omitempty"`
}

type anthropicMessageResponse struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Role         string `json:"role"`
	Model        string `json:"model"`
	Content      []any  `json:"content"`
	StopReason   string `json:"stop_reason"`
	StopSequence any    `json:"stop_sequence"`
	Usage        any    `json:"usage"`
}

// ----------------------
// Streaming chunk types
// ----------------------

type openaiChatCompletionChunk struct {
	Model  string `json:"model,omitempty"`
	Choices []struct {
		Delta struct {
			Content *string `json:"content,omitempty"`
			ToolCalls []struct {
				Index int `json:"index,omitempty"`
				ID    string `json:"id,omitempty"`
				Type  string `json:"type,omitempty"`
				Function struct {
					Name      string `json:"name,omitempty"`
					Arguments string `json:"arguments,omitempty"`
				} `json:"function,omitempty"`
			} `json:"tool_calls,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason,omitempty"`
	} `json:"choices"`
}
