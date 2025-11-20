package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// AI configuration
	AI AIConfig `yaml:"ai"`

	// CLI configuration
	CLI CLIConfig `yaml:"cli"`

	// Future extensions
	Extensions ExtensionsConfig `yaml:"extensions"`
}

// AIConfig contains AI-related configuration
type AIConfig struct {
	// LLM provider configuration (using eino-ext)
	Provider string `yaml:"provider"` // e.g., "openai", "anthropic", etc.
	APIKey   string `yaml:"api_key"`
	Model    string `yaml:"model"`
	BaseURL  string `yaml:"base_url,omitempty"`

	// Agent configuration (using eino framework)
	Agent AgentConfig `yaml:"agent"`
}

// AgentConfig contains agent framework configuration
type AgentConfig struct {
	MaxIterations  int     `yaml:"max_iterations"`
	Temperature    float64 `yaml:"temperature"`
	SystemPrompt   string  `yaml:"system_prompt"`
	HistoryEnabled bool    `yaml:"history_enabled"`
	MaxHistory     int     `yaml:"max_history"`
}

// CLIConfig contains CLI-related configuration
type CLIConfig struct {
	Interactive bool   `yaml:"interactive"`
	Prompt      string `yaml:"prompt"`
	HistoryFile string `yaml:"history_file"`
}

// ExtensionsConfig contains configuration for future extensions
type ExtensionsConfig struct {
	// MCP (Model Context Protocol) support for Anthropic
	MCPEnabled bool     `yaml:"mcp_enabled"`
	MCPServers []string `yaml:"mcp_servers,omitempty"`

	// Skills framework support
	SkillsEnabled bool     `yaml:"skills_enabled"`
	SkillsPaths   []string `yaml:"skills_paths,omitempty"`

	// Desktop integration for Windows/Mac
	DesktopEnabled  bool   `yaml:"desktop_enabled"`
	DesktopHotkey   string `yaml:"desktop_hotkey,omitempty"`
	DesktopPosition string `yaml:"desktop_position,omitempty"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		AI: AIConfig{
			Provider: "ollama",
			Model:    "deepseek-r1:8b",
			BaseURL:  "http://localhost:11434",
			Agent: AgentConfig{
				MaxIterations:  10,
				Temperature:    0.7,
				SystemPrompt:   "You are WorkWise, an intelligent desktop assistant. Help users with their tasks efficiently and professionally.",
				HistoryEnabled: true,
				MaxHistory:     50,
			},
		},
		CLI: CLIConfig{
			Interactive: true,
			Prompt:      "WorkWise> ",
			HistoryFile: filepath.Join(homeDir, ".workwise_history"),
		},
		Extensions: ExtensionsConfig{
			MCPEnabled:     false,
			SkillsEnabled:  false,
			DesktopEnabled: false,
		},
	}
}

// Load loads configuration from file or returns default
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Try to load from config file
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables if present
	if apiKey := os.Getenv("WORKWISE_API_KEY"); apiKey != "" {
		cfg.AI.APIKey = apiKey
	}
	if provider := os.Getenv("WORKWISE_PROVIDER"); provider != "" {
		cfg.AI.Provider = provider
	}
	if model := os.Getenv("WORKWISE_MODEL"); model != "" {
		cfg.AI.Model = model
	}
	if baseURL := os.Getenv("WORKWISE_BASE_URL"); baseURL != "" {
		cfg.AI.BaseURL = baseURL
	}

	return cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	// Check if custom config path is set
	if configPath := os.Getenv("WORKWISE_CONFIG"); configPath != "" {
		return configPath
	}

	// Use default location in user's home directory
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".workwise", "config.yaml")
}
