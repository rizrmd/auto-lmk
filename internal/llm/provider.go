package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Provider interface for LLM implementations
type Provider interface {
	Chat(ctx context.Context, messages []Message, functions []Function) (*Response, error)
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// Function represents a callable function for LLM
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Response from LLM
type Response struct {
	Content      string
	FunctionCall *FunctionCall
}

// FunctionCall represents a function call from LLM
type FunctionCall struct {
	Name      string
	Arguments map[string]interface{}
}

// Config for LLM provider
type Config struct {
	Provider    string
	APIKey      string
	Model       string
	ZAIEndpoint string
}

// NewProvider creates appropriate LLM provider
func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case "openai":
		return NewOpenAIProvider(cfg.APIKey, cfg.Model)
	case "anthropic":
		return NewAnthropicProvider(cfg.APIKey, cfg.Model)
	case "zai":
		return NewZAIProvider(cfg.APIKey, cfg.Model, cfg.ZAIEndpoint)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}

// OpenAIProvider implements OpenAI
type OpenAIProvider struct {
	apiKey string
	model  string
}

func NewOpenAIProvider(apiKey, model string) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key required")
	}
	if model == "" {
		model = "gpt-4o-mini" // Default to cheaper model
	}
	return &OpenAIProvider{apiKey: apiKey, model: model}, nil
}

func (p *OpenAIProvider) Chat(ctx context.Context, messages []Message, functions []Function) (*Response, error) {
	// TODO: Implement OpenAI API call
	// 1. Build request with messages and functions
	// 2. Call OpenAI API
	// 3. Parse response
	// 4. Handle function calls if present
	// 5. Return response

	return nil, fmt.Errorf("OpenAI integration not yet implemented")
}

// AnthropicProvider implements Anthropic Claude
type AnthropicProvider struct {
	apiKey string
	model  string
}

func NewAnthropicProvider(apiKey, model string) (*AnthropicProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Anthropic API key required")
	}
	if model == "" {
		model = "claude-3-5-haiku-20241022" // Default to Haiku
	}
	return &AnthropicProvider{apiKey: apiKey, model: model}, nil
}

func (p *AnthropicProvider) Chat(ctx context.Context, messages []Message, functions []Function) (*Response, error) {
	// TODO: Implement Anthropic API call
	// 1. Build request with messages and tools
	// 2. Call Anthropic API
	// 3. Parse response
	// 4. Handle tool use if present
	// 5. Return response

	return nil, fmt.Errorf("Anthropic integration not yet implemented")
}

// ZAIProvider implements Z.AI API
type ZAIProvider struct {
	apiKey   string
	model    string
	endpoint string
	client   *http.Client
}

func NewZAIProvider(apiKey, model, endpoint string) (*ZAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Z.AI API key required")
	}
	if model == "" {
		model = "glm-4-flash" // Default to GLM-4-Flash
	}
	if endpoint == "" {
		endpoint = "https://api.z.ai/api/coding/paas/v4"
	}
	return &ZAIProvider{
		apiKey:   apiKey,
		model:    model,
		endpoint: endpoint,
		client:   &http.Client{},
	}, nil
}

// Z.AI API request/response structures
type zaiRequest struct {
	Model    string       `json:"model"`
	Messages []Message    `json:"messages"`
	Tools    []zaiTool    `json:"tools,omitempty"`
	Stream   bool         `json:"stream"`
}

type zaiTool struct {
	Type     string      `json:"type"`
	Function zaiFunction `json:"function"`
}

type zaiFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type zaiResponse struct {
	Choices []struct {
		Message struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func (p *ZAIProvider) Chat(ctx context.Context, messages []Message, functions []Function) (*Response, error) {
	// Build request
	req := zaiRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   false,
	}

	// Convert functions to tools format
	if len(functions) > 0 {
		req.Tools = make([]zaiTool, len(functions))
		for i, fn := range functions {
			req.Tools[i] = zaiTool{
				Type: "function",
				Function: zaiFunction{
					Name:        fn.Name,
					Description: fn.Description,
					Parameters:  fn.Parameters,
				},
			}
		}
	}

	// Marshal request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	// Send request
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Z.AI API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var zaiResp zaiResponse
	if err := json.Unmarshal(body, &zaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(zaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := zaiResp.Choices[0]
	response := &Response{
		Content: choice.Message.Content,
	}

	// Handle function calls
	if len(choice.Message.ToolCalls) > 0 {
		toolCall := choice.Message.ToolCalls[0]
		args, err := ParseFunctionArguments(toolCall.Function.Arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse function arguments: %w", err)
		}

		response.FunctionCall = &FunctionCall{
			Name:      toolCall.Function.Name,
			Arguments: args,
		}
	}

	return response, nil
}

// Helper to parse function arguments
func ParseFunctionArguments(jsonStr string) (map[string]interface{}, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &args); err != nil {
		return nil, fmt.Errorf("failed to parse function arguments: %w", err)
	}
	return args, nil
}
