package llm

import (
	"context"
	"encoding/json"
	"fmt"
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
	Provider string
	APIKey   string
	Model    string
}

// NewProvider creates appropriate LLM provider
func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case "openai":
		return NewOpenAIProvider(cfg.APIKey, cfg.Model)
	case "anthropic":
		return NewAnthropicProvider(cfg.APIKey, cfg.Model)
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

// Helper to parse function arguments
func ParseFunctionArguments(jsonStr string) (map[string]interface{}, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &args); err != nil {
		return nil, fmt.Errorf("failed to parse function arguments: %w", err)
	}
	return args, nil
}
