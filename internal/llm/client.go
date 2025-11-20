package llm

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/warm3snow/WorkWise/internal/config"
)

// NewClient creates a new LLM client based on configuration
// Uses cloudwego/eino-ext for LLM provider integration
func NewClient(cfg *config.Config) (model.ChatModel, error) {
	// API key is required for OpenAI but not for Ollama
	if cfg.AI.Provider == "openai" && cfg.AI.APIKey == "" {
		return nil, fmt.Errorf("API key is required. Please set WORKWISE_API_KEY environment variable or configure it in config file")
	}

	switch cfg.AI.Provider {
	case "openai":
		return newOpenAIClient(cfg)
	case "ollama":
		return newOllamaClient(cfg)
	// Future providers can be added here
	// case "anthropic":
	//     return newAnthropicClient(cfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.AI.Provider)
	}
}

// newOpenAIClient creates an OpenAI client using eino-ext
func newOpenAIClient(cfg *config.Config) (model.ChatModel, error) {
	clientConfig := &openai.ChatModelConfig{
		APIKey: cfg.AI.APIKey,
		Model:  cfg.AI.Model,
	}

	// Set base URL if provided (for compatible APIs)
	if cfg.AI.BaseURL != "" {
		clientConfig.BaseURL = cfg.AI.BaseURL
	}

	// Set temperature
	if cfg.AI.Agent.Temperature > 0 {
		temp := float32(cfg.AI.Agent.Temperature)
		clientConfig.Temperature = &temp
	}

	client, err := openai.NewChatModel(context.Background(), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	return client, nil
}

// newOllamaClient creates an Ollama client using eino-ext
func newOllamaClient(cfg *config.Config) (model.ChatModel, error) {
	clientConfig := &ollama.ChatModelConfig{
		Model: cfg.AI.Model,
	}

	// Set base URL if provided, otherwise use default Ollama endpoint
	if cfg.AI.BaseURL != "" {
		clientConfig.BaseURL = cfg.AI.BaseURL
	} else {
		// Default Ollama base URL
		clientConfig.BaseURL = "http://localhost:11434"
	}

	// Set temperature through Options
	if cfg.AI.Agent.Temperature > 0 {
		clientConfig.Options = &ollama.Options{
			Temperature: float32(cfg.AI.Agent.Temperature),
		}
	}

	client, err := ollama.NewChatModel(context.Background(), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	return client, nil
}

// Future: Add support for other providers
// func newAnthropicClient(cfg *config.Config) (model.ChatModel, error) {
//     // Implementation for Anthropic Claude using eino-ext
//     // This will be added when eino-ext supports Anthropic
//     return nil, fmt.Errorf("anthropic provider not yet implemented")
// }
