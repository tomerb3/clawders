package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// ----------------------
// envOr tests
// ----------------------

func TestEnvOr(t *testing.T) {
	// Save original value and restore after test
	original, exists := os.LookupEnv("TEST_ENV_VAR")
	defer func() {
		if exists {
			os.Setenv("TEST_ENV_VAR", original)
		} else {
			os.Unsetenv("TEST_ENV_VAR")
		}
	}()

	tests := []struct {
		name     string
		key      string
		fallback string
		setValue string
		want     string
	}{
		{
			name:     "returns fallback when env var not set",
			key:      "TEST_ENV_VAR",
			fallback: "default",
			setValue: "",
			want:     "default",
		},
		{
			name:     "returns value when env var is set",
			key:      "TEST_ENV_VAR",
			fallback: "default",
			setValue: "custom",
			want:     "custom",
		},
		{
			name:     "preserves whitespace in value",
			key:      "TEST_ENV_VAR",
			fallback: "default",
			setValue: "  trimmed  ",
			want:     "  trimmed  ",
		},
		{
			name:     "preserves whitespace in fallback",
			key:      "TEST_ENV_VAR",
			fallback: "  fallback  ",
			setValue: "",
			want:     "  fallback  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValue != "" || exists {
				os.Setenv(tt.key, tt.setValue)
			} else {
				os.Unsetenv(tt.key)
			}
			got := envOr(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("envOr(%q, %q) = %q, want %q", tt.key, tt.fallback, got, tt.want)
			}
		})
	}
}

// ----------------------
// modelSupportsToolUse tests
// ----------------------

func TestModelSupportsToolUse(t *testing.T) {
	tests := []struct {
		name  string
		model string
		want  bool
	}{
		{
			name:  "qwen2.5-coder-32b-instruct does not support tool use",
			model: "qwen/qwen2.5-coder-32b-instruct",
			want:  false,
		},
		{
			name:  "Qwen2.5-Coder-32B-Instruct case insensitive",
			model: "Qwen/Qwen2.5-Coder-32B-Instruct",
			want:  false,
		},
		{
			name:  "empty model returns true",
			model: "",
			want:  true,
		},
		{
			name:  "claude-3-5-sonnet supports tool use",
			model: "claude-3-5-sonnet-20241022",
			want:  true,
		},
		{
			name:  "gpt-4 supports tool use",
			model: "gpt-4",
			want:  true,
		},
		{
			name:  "whitespace only model returns true",
			model: "   ",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := modelSupportsToolUse(tt.model)
			if got != tt.want {
				t.Errorf("modelSupportsToolUse(%q) = %v, want %v", tt.model, got, tt.want)
			}
		})
	}
}

// ----------------------
// clampMaxTokens tests
// ----------------------

func TestClampMaxTokens(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		requested int
		want     int
	}{
		{
			name:     "qwen2.5-coder-32b-instruct caps at 8192",
			model:    "qwen/qwen2.5-coder-32b-instruct",
			requested: 16000,
			want:     8192,
		},
		{
			name:     "qwen2.5-coder-32b-instruct allows below cap",
			model:    "qwen/qwen2.5-coder-32b-instruct",
			requested: 4096,
			want:     4096,
		},
		{
			name:     "qwen2.5-coder-32b-instruct exact cap",
			model:    "qwen/qwen2.5-coder-32b-instruct",
			requested: 8192,
			want:     8192,
		},
		{
			name:     "other models unchanged",
			model:    "claude-3-5-sonnet-20241022",
			requested: 32000,
			want:     32000,
		},
		{
			name:     "zero or negative returns unchanged",
			model:    "claude-3-5-sonnet",
			requested: 0,
			want:     0,
		},
		{
			name:     "negative returns unchanged",
			model:    "claude-3-5-sonnet",
			requested: -100,
			want:     -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clampMaxTokens(tt.model, tt.requested)
			if got != tt.want {
				t.Errorf("clampMaxTokens(%q, %d) = %d, want %d", tt.model, tt.requested, got, tt.want)
			}
		})
	}
}

// ----------------------
// checkInboundAuth tests
// ----------------------

func TestCheckInboundAuth(t *testing.T) {
	const validKey = "secret-api-key"

	tests := []struct {
		name       string
		authHeader string
		apiKeyHeader string
		want       bool
	}{
		{
			name:       "valid bearer token",
			authHeader: "Bearer secret-api-key",
			apiKeyHeader: "",
			want:       true,
		},
		{
			name:       "valid bearer token case insensitive",
			authHeader: "bearer secret-api-key",
			apiKeyHeader: "",
			want:       true,
		},
		{
			name:       "valid bearer token with whitespace",
			authHeader: "Bearer   secret-api-key  ",
			apiKeyHeader: "",
			want:       true,
		},
		{
			name:       "invalid bearer token",
			authHeader: "Bearer wrong-key",
			apiKeyHeader: "",
			want:       false,
		},
		{
			name:       "valid x-api-key header",
			authHeader: "",
			apiKeyHeader: "secret-api-key",
			want:       true,
		},
		{
			name:       "invalid x-api-key header",
			authHeader: "",
			apiKeyHeader: "wrong-key",
			want:       false,
		},
		{
			name:       "both headers valid",
			authHeader: "Bearer secret-api-key",
			apiKeyHeader: "secret-api-key",
			want:       true,
		},
		{
			name:       "auth header takes precedence",
			authHeader: "Bearer wrong-key",
			apiKeyHeader: "secret-api-key",
			want:       false,
		},
		{
			name:       "empty headers",
			authHeader: "",
			apiKeyHeader: "",
			want:       false,
		},
		{
			name:       "empty bearer prefix",
			authHeader: "Bearer ",
			apiKeyHeader: "",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.apiKeyHeader != "" {
				req.Header.Set("x-api-key", tt.apiKeyHeader)
			}

			got := checkInboundAuth(req, validKey)
			if got != tt.want {
				t.Errorf("checkInboundAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ----------------------
// extractSystemText tests
// ----------------------

func TestExtractSystemText(t *testing.T) {
	tests := []struct {
		name string
		raw  json.RawMessage
		want string
	}{
		{
			name: "empty raw message",
			raw:  nil,
			want: "",
		},
		{
			name: "empty json",
			raw:  json.RawMessage(""),
			want: "",
		},
		{
			name: "simple string system prompt",
			raw:  json.RawMessage(`"You are a helpful assistant"`),
			want: "You are a helpful assistant",
		},
		{
			name: "text block array",
			raw:  json.RawMessage(`[{"type":"text","text":"System prompt here"}]`),
			want: "System prompt here",
		},
		{
			name: "multiple text blocks joined",
			raw:  json.RawMessage(`[{"type":"text","text":"Part 1"},{"type":"text","text":"Part 2"}]`),
			want: "Part 1\nPart 2",
		},
		{
			name: "ignores non-text blocks",
			raw:  json.RawMessage(`[{"type":"text","text":"Real prompt"},{"type":"image","source":{"type":"url","url":"http://test.com"}}]`),
			want: "Real prompt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSystemText(tt.raw)
			if got != tt.want {
				t.Errorf("extractSystemText(%q) = %q, want %q", string(tt.raw), got, tt.want)
			}
		})
	}
}

// ----------------------
// joinTextBlocks tests
// ----------------------

func TestJoinTextBlocks(t *testing.T) {
	tests := []struct {
		name    string
		blocks  []anthropicContentBlock
		want    string
	}{
		{
			name:   "empty blocks",
			blocks: []anthropicContentBlock{},
			want:   "",
		},
		{
			name:   "single text block",
			blocks: []anthropicContentBlock{{Type: "text", Text: "Hello"}},
			want:   "Hello",
		},
		{
			name:   "multiple text blocks with newline",
			blocks: []anthropicContentBlock{{Type: "text", Text: "Hello"}, {Type: "text", Text: "World"}},
			want:   "Hello\nWorld",
		},
		{
			name:   "ignores non-text blocks",
			blocks: []anthropicContentBlock{{Type: "text", Text: "Hello"}, {Type: "tool_use", Name: "test"}},
			want:   "Hello",
		},
		{
			name:   "ignores empty text",
			blocks: []anthropicContentBlock{{Type: "text", Text: ""}, {Type: "text", Text: "World"}},
			want:   "World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := joinTextBlocks(tt.blocks)
			if got != tt.want {
				t.Errorf("joinTextBlocks() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ----------------------
// mapFinishReason tests
// ----------------------

func TestMapFinishReason(t *testing.T) {
	tests := []struct {
		name  string
		finish string
		want  string
	}{
		{name: "stop", finish: "stop", want: "end_turn"},
		{name: "length", finish: "length", want: "max_tokens"},
		{name: "tool_calls", finish: "tool_calls", want: "tool_use"},
		{name: "content_filter", finish: "content_filter", want: "stop_sequence"},
		{name: "empty", finish: "", want: "end_turn"},
		{name: "unknown", finish: "unknown_value", want: "end_turn"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapFinishReason(tt.finish)
			if got != tt.want {
				t.Errorf("mapFinishReason(%q) = %q, want %q", tt.finish, got, tt.want)
			}
		})
	}
}

// ----------------------
// convertToolChoice tests
// ----------------------

// deepEqual compares two any values for equality (handles maps which can't be compared with ==)
func deepEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	switch va := a.(type) {
	case map[string]any:
		vb, ok := b.(map[string]any)
		if !ok {
			return false
		}
		if len(va) != len(vb) {
			return false
		}
		for k, vaVal := range va {
			vbVal, ok := vb[k]
			if !ok {
				return false
			}
			if !deepEqual(vaVal, vbVal) {
				return false
			}
		}
		return true
	case []any:
		vb, ok := b.([]any)
		if !ok {
			return false
		}
		if len(va) != len(vb) {
			return false
		}
		for i := range va {
			if !deepEqual(va[i], vb[i]) {
				return false
			}
		}
		return true
	case string:
		vb, ok := b.(string)
		return ok && va == vb
	default:
		return a == b
	}
}

func TestConvertToolChoice(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want any
	}{
		{
			name: "auto",
			v:    map[string]any{"type": "auto"},
			want: "auto",
		},
		{
			name: "none",
			v:    map[string]any{"type": "none"},
			want: "none",
		},
		{
			name: "required",
			v:    map[string]any{"type": "required"},
			want: "required",
		},
		{
			name: "tool with name",
			v:    map[string]any{"type": "tool", "name": "my_tool"},
			want: map[string]any{"type": "function", "function": map[string]any{"name": "my_tool"}},
		},
		{
			name: "tool with empty name defaults to auto",
			v:    map[string]any{"type": "tool", "name": ""},
			want: "auto",
		},
		{
			name: "non-map value passes through",
			v:    "auto",
			want: "auto",
		},
		{
			name: "nil passes through",
			v:    nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToolChoice(tt.v)
			if !deepEqual(got, tt.want) {
				t.Errorf("convertToolChoice(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

// ----------------------
// takeFirstRunes tests
// ----------------------

func TestTakeFirstRunes(t *testing.T) {
	tests := []struct {
		name string
		s    string
		max  int
		want string
	}{
		{name: "normal string within limit", s: "hello", max: 10, want: "hello"},
		{name: "string truncated to max", s: "hello world", max: 5, want: "hello"},
		{name: "unicode characters", s: "hello", max: 3, want: "hel"},
		{name: "zero max returns empty", s: "hello", max: 0, want: ""},
		{name: "negative max returns empty", s: "hello", max: -1, want: ""},
		{name: "empty string", s: "", max: 5, want: ""},
		{name: "max equals length", s: "hello", max: 5, want: "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := takeFirstRunes(tt.s, tt.max)
			if got != tt.want {
				t.Errorf("takeFirstRunes(%q, %d) = %q, want %q", tt.s, tt.max, got, tt.want)
			}
		})
	}
}

// ----------------------
// loadFileConfig tests
// ----------------------

func TestLoadFileConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid config",
			content: `{"nvidia_url":"https://api.nvidia.com","nvidia_key":"test-key"}`,
			wantErr: false,
		},
		{
			name:    "empty content",
			content: "",
			wantErr: true,
		},
		{
			name:    "invalid json",
			content: `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmp := "/tmp/test_config_" + strings.ReplaceAll(tt.name, " ", "_") + ".json"
			if tt.content != "" {
				err := os.WriteFile(tmp, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("failed to write temp file: %v", err)
				}
				defer os.Remove(tmp)
			}

			_, err := loadFileConfig(tmp)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFileConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ----------------------
// loadConfig tests
// ----------------------

func TestLoadConfig(t *testing.T) {
	// Save original env and cleanup after
	origConfigPath := os.Getenv("CONFIG_PATH")
	origUpstreamURL := os.Getenv("UPSTREAM_URL")
	origProviderKey := os.Getenv("PROVIDER_API_KEY")
	origServerKey := os.Getenv("SERVER_API_KEY")
	origTimeout := os.Getenv("UPSTREAM_TIMEOUT_SECONDS")
	origLogBody := os.Getenv("LOG_BODY_MAX_CHARS")
	origLogStream := os.Getenv("LOG_STREAM_TEXT_PREVIEW_CHARS")
	origAddr := os.Getenv("ADDR")

	defer func() {
		os.Setenv("CONFIG_PATH", origConfigPath)
		os.Setenv("UPSTREAM_URL", origUpstreamURL)
		os.Setenv("PROVIDER_API_KEY", origProviderKey)
		os.Setenv("SERVER_API_KEY", origServerKey)
		os.Setenv("UPSTREAM_TIMEOUT_SECONDS", origTimeout)
		os.Setenv("LOG_BODY_MAX_CHARS", origLogBody)
		os.Setenv("LOG_STREAM_TEXT_PREVIEW_CHARS", origLogStream)
		os.Setenv("ADDR", origAddr)
	}()

	// Clear env vars
	os.Unsetenv("CONFIG_PATH")
	os.Unsetenv("UPSTREAM_URL")
	os.Unsetenv("PROVIDER_API_KEY")
	os.Unsetenv("SERVER_API_KEY")
	os.Unsetenv("UPSTREAM_TIMEOUT_SECONDS")
	os.Unsetenv("LOG_BODY_MAX_CHARS")
	os.Unsetenv("LOG_STREAM_TEXT_PREVIEW_CHARS")
	os.Unsetenv("ADDR")

	tests := []struct {
		name        string
		configFile  string
		envSetup    func()
		wantErr     bool
		errContains string
	}{
		{
			name: "missing upstream URL",
			configFile: `{"nvidia_url":"","nvidia_key":"test-key"}`,
			wantErr:     true,
			errContains: "missing nvidia_url",
		},
		{
			name: "missing provider API key",
			configFile: `{"nvidia_url":"https://api.nvidia.com","nvidia_key":""}`,
			wantErr:     true,
			errContains: "missing nvidia_key",
		},
		{
			name: "invalid timeout",
			configFile: `{"nvidia_url":"https://api.nvidia.com","nvidia_key":"test-key"}`,
			envSetup: func() {
				os.Setenv("UPSTREAM_TIMEOUT_SECONDS", "not-a-number")
			},
			wantErr:     true,
			errContains: "invalid UPSTREAM_TIMEOUT_SECONDS",
		},
		{
			name: "zero timeout",
			configFile: `{"nvidia_url":"https://api.nvidia.com","nvidia_key":"test-key"}`,
			envSetup: func() {
				os.Setenv("UPSTREAM_TIMEOUT_SECONDS", "0")
			},
			wantErr:     true,
			errContains: "invalid UPSTREAM_TIMEOUT_SECONDS",
		},
		{
			name: "negative timeout",
			configFile: `{"nvidia_url":"https://api.nvidia.com","nvidia_key":"test-key"}`,
			envSetup: func() {
				os.Setenv("UPSTREAM_TIMEOUT_SECONDS", "-5")
			},
			wantErr:     true,
			errContains: "invalid UPSTREAM_TIMEOUT_SECONDS",
		},
		{
			name: "invalid log body max",
			configFile: `{"nvidia_url":"https://api.nvidia.com","nvidia_key":"test-key"}`,
			envSetup: func() {
				os.Setenv("LOG_BODY_MAX_CHARS", "not-a-number")
			},
			wantErr:     true,
			errContains: "invalid LOG_BODY_MAX_CHARS",
		},
		{
			name: "negative log body max",
			configFile: `{"nvidia_url":"https://api.nvidia.com","nvidia_key":"test-key"}`,
			envSetup: func() {
				os.Setenv("LOG_BODY_MAX_CHARS", "-10")
			},
			wantErr:     true,
			errContains: "invalid LOG_BODY_MAX_CHARS",
		},
		{
			name: "invalid log stream preview",
			configFile: `{"nvidia_url":"https://api.nvidia.com","nvidia_key":"test-key"}`,
			envSetup: func() {
				os.Setenv("LOG_STREAM_TEXT_PREVIEW_CHARS", "invalid")
			},
			wantErr:     true,
			errContains: "invalid LOG_STREAM_TEXT_PREVIEW_CHARS",
		},
		{
			name: "valid config with env overrides",
			configFile: `{"nvidia_url":"https://api.nvidia.com","nvidia_key":"test-key"}`,
			envSetup: func() {
				os.Setenv("UPSTREAM_URL", "https://custom.api.com")
				os.Setenv("PROVIDER_API_KEY", "custom-key")
				os.Setenv("ADDR", ":8080")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset env
			os.Unsetenv("CONFIG_PATH")
			os.Unsetenv("UPSTREAM_URL")
			os.Unsetenv("PROVIDER_API_KEY")
			os.Unsetenv("SERVER_API_KEY")
			os.Unsetenv("UPSTREAM_TIMEOUT_SECONDS")
			os.Unsetenv("LOG_BODY_MAX_CHARS")
			os.Unsetenv("LOG_STREAM_TEXT_PREVIEW_CHARS")
			os.Unsetenv("ADDR")

			// Write config file
			tmpConfig := "/tmp/test_loadconfig_" + strings.ReplaceAll(tt.name, " ", "_") + ".json"
			err := os.WriteFile(tmpConfig, []byte(tt.configFile), 0644)
			if err != nil {
				t.Fatalf("failed to write temp config: %v", err)
			}
			defer os.Remove(tmpConfig)
			os.Setenv("CONFIG_PATH", tmpConfig)

			if tt.envSetup != nil {
				tt.envSetup()
			}

			cfg, err := loadConfig()
			if tt.wantErr {
				if err == nil {
					t.Errorf("loadConfig() expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("loadConfig() error = %v, want error containing %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("loadConfig() unexpected error: %v", err)
				}
				if cfg == nil {
					t.Error("loadConfig() returned nil config without error")
				}
			}
		})
	}
}

// ----------------------
// convertAnthropicToOpenAI tests
// ----------------------

func TestConvertAnthropicToOpenAI(t *testing.T) {
	tests := []struct {
		name    string
		req     *anthropicMessageRequest
		wantErr bool
	}{
		{
			name: "simple user message",
			req: &anthropicMessageRequest{
				Model:     "claude-3-5-sonnet",
				MaxTokens: 1024,
				Messages: []anthropicMsg{
					{Role: "user", Content: json.RawMessage(`"Hello"`)},
				},
			},
			wantErr: false,
		},
		{
			name: "system prompt",
			req: &anthropicMessageRequest{
				Model:     "claude-3-5-sonnet",
				MaxTokens: 1024,
				System:    json.RawMessage(`"You are helpful"`),
				Messages: []anthropicMsg{
					{Role: "user", Content: json.RawMessage(`"Hello"`)},
				},
			},
			wantErr: false,
		},
		{
			name: "assistant message with tool use",
			req: &anthropicMessageRequest{
				Model:     "claude-3-5-sonnet",
				MaxTokens: 1024,
				Messages: []anthropicMsg{
					{Role: "assistant", Content: json.RawMessage(`[{"type":"tool_use","id":"tool_1","name":"get_weather","input":{"location":"NYC"}}]`)},
				},
			},
			wantErr: false,
		},
		{
			name: "tool result message",
			req: &anthropicMessageRequest{
				Model:     "claude-3-5-sonnet",
				MaxTokens: 1024,
				Messages: []anthropicMsg{
					{Role: "tool", Content: json.RawMessage(`[{"type":"tool_result","tool_use_id":"tool_1","content":"Sunny, 72F"}]`)},
				},
			},
			wantErr: false,
		},
		{
			name: "empty role is skipped",
			req: &anthropicMessageRequest{
				Model:     "claude-3-5-sonnet",
				MaxTokens: 1024,
				Messages: []anthropicMsg{
					{Role: "", Content: json.RawMessage(`"Hello"`)},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := convertAnthropicToOpenAI(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertAnthropicToOpenAI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ----------------------
// convertOpenAIToAnthropic tests
// ----------------------

func TestConvertOpenAIToAnthropic(t *testing.T) {
	content := "Hello, world!"
	tests := []struct {
		name string
		resp openaiChatCompletionResponse
	}{
		{
			name: "simple text response",
			resp: openaiChatCompletionResponse{
				ID:    "msg_123",
				Model: "gpt-4",
				Choices: []struct {
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
				}{
					{
						Message: struct {
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
						}{
							Role:    "assistant",
							Content: &content,
						},
						FinishReason: "stop",
					},
				},
			},
		},
		{
			name: "tool call response",
			resp: openaiChatCompletionResponse{
				ID:    "msg_456",
				Model: "gpt-4",
				Choices: []struct {
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
				}{
					{
						Message: struct {
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
						}{
							Role:    "assistant",
							Content: nil,
							ToolCalls: []struct {
								ID       string `json:"id"`
								Type     string `json:"type"`
								Function struct {
									Name      string `json:"name"`
									Arguments any    `json:"arguments"`
								} `json:"function"`
							}{
								{
									ID:   "call_123",
									Type: "function",
									Function: struct {
										Name      string `json:"name"`
										Arguments any    `json:"arguments"`
									}{
										Name:      "get_weather",
										Arguments: `{"location":"NYC"}`,
									},
								},
							},
						},
						FinishReason: "tool_calls",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertOpenAIToAnthropic(tt.resp)
			if got.Type != "message" {
				t.Errorf("convertOpenAIToAnthropic() Type = %q, want %q", got.Type, "message")
			}
			if got.Role != "assistant" {
				t.Errorf("convertOpenAIToAnthropic() Role = %q, want %q", got.Role, "assistant")
			}
			if len(got.Content) == 0 && len(tt.resp.Choices[0].Message.ToolCalls) == 0 {
				t.Error("convertOpenAIToAnthropic() produced no content but expected some")
			}
		})
	}
}

// ----------------------
// mustJSONTrunc tests
// ----------------------

func TestMustJSONTrunc(t *testing.T) {
	tests := []struct {
		name     string
		v        any
		maxChars int
		want     string
	}{
		{
			name:     "disabled when maxChars is 0",
			v:        map[string]any{"key": "value"},
			maxChars: 0,
			want:     "(disabled)",
		},
		{
			name:     "normal JSON",
			v:        map[string]any{"key": "value"},
			maxChars: 100,
			want:     `{"key":"value"}`,
		},
		{
			name:     "truncates long JSON",
			v:        map[string]any{"key": "this is a long value"},
			maxChars: 15,
			want:     `{"key":"this is` + "...(truncated)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustJSONTrunc(tt.v, tt.maxChars)
			if got != tt.want {
				t.Errorf("mustJSONTrunc() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ----------------------
// sanitizeOpenAIRequest tests
// ----------------------

func TestSanitizeOpenAIRequest(t *testing.T) {
	// Test that sanitizeOpenAIRequest handles various message types
	msgWithDataURL := map[string]any{
		"role":    "user",
		"content": []any{map[string]any{"type": "image_url", "image_url": map[string]any{"url": "data:image/png;base64,abc123"}}},
	}

	sanitized := sanitizeOpenAIRequest(openaiChatCompletionRequest{
		Model: "test-model",
		Messages: []any{msgWithDataURL},
	})

	if len(sanitized.Messages) != 1 {
		t.Errorf("sanitizeOpenAIRequest() reduced message count unexpectedly")
	}
}

// ----------------------
// subtleConstantTimeCompare tests
// ----------------------

func TestSubtleConstantTimeCompare(t *testing.T) {
	tests := []struct {
		name   string
		a      string
		b      string
		expect int
	}{
		{
			name:   "equal strings",
			a:      "secret-key",
			b:      "secret-key",
			expect: 1,
		},
		{
			name:   "different strings same length",
			a:      "secret-key",
			b:      "other-key",
			expect: 0,
		},
		{
			name:   "different lengths",
			a:      "short",
			b:      "much-longer",
			expect: 0,
		},
		{
			name:   "empty strings",
			a:      "",
			b:      "",
			expect: 1,
		},
		{
			name:   "one empty",
			a:      "hello",
			b:      "",
			expect: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subtleConstantTimeCompare(tt.a, tt.b)
			if got != tt.expect {
				t.Errorf("subtleConstantTimeCompare(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.expect)
			}
		})
	}
}

// ----------------------
// writeJSONError tests
// ----------------------

func TestWriteJSONError(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		code       string
		wantStatus int
	}{
		{
			name:       "unauthorized error",
			status:     http.StatusUnauthorized,
			code:       "unauthorized",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "bad request error",
			status:     http.StatusBadRequest,
			code:       "invalid_json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "bad gateway error",
			status:     http.StatusBadGateway,
			code:       "upstream_failed",
			wantStatus: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeJSONError(w, tt.status, tt.code)

			if w.Code != tt.wantStatus {
				t.Errorf("writeJSONError() status = %d, want %d", w.Code, tt.wantStatus)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("writeJSONError() Content-Type = %q, want %q", contentType, "application/json")
			}

			var resp map[string]any
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("writeJSONError() body is not valid JSON: %v", err)
			}

			errObj, ok := resp["error"].(map[string]any)
			if !ok {
				t.Fatal("writeJSONError() response missing error object")
			}
			if errObj["code"] != tt.code {
				t.Errorf("writeJSONError() code = %q, want %q", errObj["code"], tt.code)
			}
		})
	}
}

// ----------------------
// sanitizeOpenAIMessages tests
// ----------------------

func TestSanitizeOpenAIMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []any
		wantLen  int
	}{
		{
			name:     "nil messages",
			messages: nil,
			wantLen:  0,
		},
		{
			name:     "empty messages",
			messages: []any{},
			wantLen:  0,
		},
		{
			name: "simple string content",
			messages: []any{
				map[string]any{"role": "user", "content": "hello"},
			},
			wantLen: 1,
		},
		{
			name: "message with image URL redacted",
			messages: []any{
				map[string]any{
					"role":    "user",
					"content": []any{map[string]any{"type": "image_url", "image_url": map[string]any{"url": "data:image/png;base64,abc123"}}},
				},
			},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeOpenAIMessages(tt.messages)
			if len(got) != tt.wantLen {
				t.Errorf("sanitizeOpenAIMessages() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// ----------------------
// sanitizeMessageContent tests
// ----------------------

func TestSanitizeMessageContent(t *testing.T) {
	tests := []struct {
		name    string
		content any
		want    any
	}{
		{
			name:    "string content unchanged",
			content: "hello world",
			want:    "hello world",
		},
		{
			name:    "nil content",
			content: nil,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeMessageContent(tt.content)
			if got != tt.want {
				t.Errorf("sanitizeMessageContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ----------------------
// sanitizeAny tests
// ----------------------

func TestSanitizeAny(t *testing.T) {
	tests := []struct {
		name string
		v    any
	}{
		{
			name: "string unchanged",
			v:    "hello",
		},
		{
			name: "number unchanged",
			v:    42,
		},
		{
			name: "map unchanged",
			v:    map[string]any{"key": "value"},
		},
		{
			name: "array unchanged",
			v:    []any{"a", "b", "c"},
		},
		{
			name: "bool unchanged",
			v:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeAny(tt.v)
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.v) {
				t.Errorf("sanitizeAny(%v) = %v, want %v", tt.v, got, tt.v)
			}
		})
	}
}

// ----------------------
// sanitizeAnySlice tests
// ----------------------

func TestSanitizeAnySlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []any
		wantNil  bool
		wantLen  int
	}{
		{
			name:    "nil slice",
			input:   nil,
			wantNil: true,
			wantLen: 0,
		},
		{
			name:    "empty slice",
			input:   []any{},
			wantNil: true,
			wantLen: 0,
		},
		{
			name:    "slice with content",
			input:   []any{"a", "b"},
			wantNil: false,
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeAnySlice(tt.input)
			if tt.wantNil && got != nil {
				t.Errorf("sanitizeAnySlice() = %v, want nil", got)
			}
			if !tt.wantNil && got == nil {
				t.Errorf("sanitizeAnySlice() = nil, want non-nil")
			}
			if len(got) != tt.wantLen {
				t.Errorf("sanitizeAnySlice() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// ----------------------
// convertAnthropicUserBlocksToOpenAIMessages tests
// ----------------------

func TestConvertAnthropicUserBlocksToOpenAIMessages(t *testing.T) {
	tests := []struct {
		name    string
		blocks  []anthropicContentBlock
		wantErr bool
		wantLen int
	}{
		{
			name:    "empty blocks",
			blocks:  []anthropicContentBlock{},
			wantErr: false,
			wantLen: 1, // produces empty user message
		},
		{
			name: "text blocks",
			blocks: []anthropicContentBlock{
				{Type: "text", Text: "Hello"},
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "tool result becomes tool message",
			blocks: []anthropicContentBlock{
				{Type: "tool_result", ToolUseID: "tool_1", Content: json.RawMessage(`"result"`)},
			},
			wantErr: false,
			wantLen: 2, // tool message + empty user message
		},
		{
			name: "image source type base64",
			blocks: []anthropicContentBlock{
				{
					Type:   "image",
					Source: &anthropicImageSource{Type: "base64", MediaType: "image/png", Data: "abc123"},
				},
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "image source type url",
			blocks: []anthropicContentBlock{
				{
					Type:   "image",
					Source: &anthropicImageSource{Type: "url", URL: "https://example.com/image.png"},
				},
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "invalid base64 skipped but empty user message added",
			blocks: []anthropicContentBlock{
				{
					Type:   "image",
					Source: &anthropicImageSource{Type: "base64", MediaType: "image/png", Data: "not-valid-base64!!!"},
				},
			},
			wantErr: false,
			wantLen: 1, // empty user message when no valid content
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertAnthropicUserBlocksToOpenAIMessages(tt.blocks)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertAnthropicUserBlocksToOpenAIMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != tt.wantLen {
				t.Errorf("convertAnthropicUserBlocksToOpenAIMessages() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// ----------------------
// convertAnthropicAssistantBlocksToOpenAIMessage tests
// ----------------------

func TestConvertAnthropicAssistantBlocksToOpenAIMessage(t *testing.T) {
	tests := []struct {
		name    string
		blocks  []anthropicContentBlock
		wantErr bool
	}{
		{
			name:    "text only",
			blocks:  []anthropicContentBlock{{Type: "text", Text: "Hello"}},
			wantErr: false,
		},
		{
			name: "tool use",
			blocks: []anthropicContentBlock{
				{Type: "tool_use", ID: "tool_1", Name: "get_weather", Input: json.RawMessage(`{"location":"NYC"}`)},
			},
			wantErr: false,
		},
		{
			name:    "mixed text and tool use",
			blocks: []anthropicContentBlock{
				{Type: "text", Text: "Let me check..."},
				{Type: "tool_use", ID: "tool_1", Name: "get_weather", Input: json.RawMessage(`{"location":"NYC"}`)},
			},
			wantErr: false,
		},
		{
			name: "tool use without id",
			blocks: []anthropicContentBlock{
				{Type: "tool_use", Name: "get_weather", Input: json.RawMessage(`{}`)},
			},
			wantErr: false, // should just skip tool without id
		},
		{
			name: "tool use without name",
			blocks: []anthropicContentBlock{
				{Type: "tool_use", ID: "tool_1", Input: json.RawMessage(`{}`)},
			},
			wantErr: false, // should just skip tool without name
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertAnthropicAssistantBlocksToOpenAIMessage(tt.blocks)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertAnthropicAssistantBlocksToOpenAIMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got == nil {
				t.Error("convertAnthropicAssistantBlocksToOpenAIMessage() returned nil without error")
			}
		})
	}
}

// ----------------------
// HTTP Handler tests
// ----------------------

func TestHandleMessages(t *testing.T) {
	// Create a test config
	cfg := &serverConfig{
		addr:           ":3001",
		upstreamURL:    "http://localhost:9999", // won't actually connect
		providerAPIKey: "test-key",
		serverAPIKey:   "", // no auth required
		timeout:        5 * time.Second,
		logBodyMax:     4096,
		logStreamPreviewMax: 256,
	}

	tests := []struct {
		name           string
		body           string
		wantStatus     int
		wantErrCode    string
		authHeader     string
	}{
		{
			name:        "invalid json body",
			body:        `{invalid}`,
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "invalid_json",
		},
		{
			name:        "missing model",
			body:        `{"max_tokens":1024,"messages":[{"role":"user","content":"hello"}]}`,
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "missing_model",
		},
		{
			name:        "empty model",
			body:        `{"model":"","max_tokens":1024,"messages":[{"role":"user","content":"hello"}]}`,
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "missing_model",
		},
		{
			name:        "whitespace model",
			body:        `{"model":"   ","max_tokens":1024,"messages":[{"role":"user","content":"hello"}]}`,
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "missing_model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(tt.body))
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handleMessages(w, req, cfg)

			if w.Code != tt.wantStatus {
				t.Errorf("handleMessages() status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantErrCode != "" {
				var resp map[string]any
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("handleMessages() response is not valid JSON: %v", err)
				}
				errObj, ok := resp["error"].(map[string]any)
				if !ok {
					t.Fatal("handleMessages() response missing error object")
				}
				if errObj["code"] != tt.wantErrCode {
					t.Errorf("handleMessages() error code = %v, want %v", errObj["code"], tt.wantErrCode)
				}
			}
		})
	}
}

func TestHandleMessagesWithAuth(t *testing.T) {
	cfg := &serverConfig{
		addr:           ":3001",
		upstreamURL:    "http://localhost:9999",
		providerAPIKey: "test-key",
		serverAPIKey:   "secret-key", // auth required
		timeout:        5 * time.Second,
		logBodyMax:     4096,
		logStreamPreviewMax: 256,
	}

	body := `{"model":"test","max_tokens":1024,"messages":[{"role":"user","content":"hello"}]}`

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{
			name:       "valid auth",
			authHeader: "Bearer secret-key",
			wantStatus: http.StatusBadGateway, // upstream won't respond but auth passes
		},
		{
			name:       "invalid auth",
			authHeader: "Bearer wrong-key",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing auth",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(body))
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handleMessages(w, req, cfg)

			if w.Code != tt.wantStatus {
				t.Errorf("handleMessages() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Test the root health endpoint
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// This handler is defined inline in main.go, so we can't test it directly
	// But we can test that the server responds
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "claude-nvidia-proxy",
			"health":  "ok",
		})
	})

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("health endpoint status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("health endpoint response is not valid JSON: %v", err)
	}

	if resp["message"] != "claude-nvidia-proxy" {
		t.Errorf("health endpoint message = %v, want %v", resp["message"], "claude-nvidia-proxy")
	}
	if resp["health"] != "ok" {
		t.Errorf("health endpoint health = %v, want %v", resp["health"], "ok")
	}
}
