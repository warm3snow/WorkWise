# WorkWise Development Guide

## Project Structure

```
WorkWise/
├── cmd/
│   └── workwise/          # Main application entry point
│       └── main.go        # Application bootstrap
├── internal/              # Internal packages (not importable by external packages)
│   ├── agent/             # AI agent implementation using Eino framework
│   │   └── agent.go       # Agent core logic with conversation management
│   ├── cli/               # Command-line interface implementation
│   │   └── cli.go         # CLI commands and interactive mode
│   ├── config/            # Configuration management
│   │   └── config.go      # Config structure, loading, and defaults
│   └── llm/               # LLM client implementation
│       └── client.go      # LLM provider integration using Eino-Ext
├── pkg/                   # Public packages (can be imported by external packages)
│   ├── mcp/               # Model Context Protocol support (future)
│   │   └── mcp.go         # MCP server management interface
│   └── skills/            # Skills framework (future)
│       └── skills.go      # Skills registry and execution
├── .gitignore            # Git ignore patterns
├── config.example.yaml   # Example configuration file
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── Makefile              # Build automation
└── README.md             # Project documentation
```

## Architecture Overview

### Core Components

#### 1. Agent (`internal/agent`)
- Implements the AI agent using CloudWeGo Eino framework
- Manages conversation history and context
- Orchestrates LLM calls and response processing
- Features:
  - Conversation history management with configurable limits
  - System prompt configuration
  - Context-aware responses

#### 2. LLM Client (`internal/llm`)
- Abstracts LLM provider integration
- Uses CloudWeGo Eino-Ext for provider implementations
- Currently supports:
  - OpenAI (GPT-4, GPT-3.5-turbo, etc.)
  - OpenAI-compatible APIs (via custom base URL)
- Future providers: Anthropic Claude, etc.

#### 3. CLI (`internal/cli`)
- Provides command-line interface using urfave/cli
- Supports both interactive and one-shot modes
- Commands:
  - `chat`: Interactive conversation mode
  - `ask`: Single question mode
  - `config show/init`: Configuration management
  - `version`: Version information

#### 4. Configuration (`internal/config`)
- YAML-based configuration
- Environment variable override support
- Default configuration with sensible values
- Configuration file location: `~/.workwise/config.yaml`

### Extension Framework

#### MCP Support (`pkg/mcp`)
- Placeholder for Anthropic's Model Context Protocol
- Designed for extensibility
- Will enable:
  - Tool/function calling
  - External context integration
  - Resource access control

#### Skills Framework (`pkg/skills`)
- Extensible capability system
- Registry-based skill management
- Future skills:
  - File operations
  - Web search
  - Code execution
  - System commands
  - Custom user-defined skills

## Development Workflow

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install to GOPATH/bin
make install
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter (requires golangci-lint)
make lint
```

### Running

```bash
# Run in interactive mode
make run

# Or directly
./workwise

# Ask a single question
./workwise ask "How do I list files in Linux?"
```

## Configuration

### Configuration Priority
1. Environment variables (highest priority)
2. Configuration file
3. Default values (lowest priority)

### Environment Variables
- `WORKWISE_API_KEY`: LLM provider API key
- `WORKWISE_PROVIDER`: LLM provider (openai, etc.)
- `WORKWISE_MODEL`: Model name (gpt-4, etc.)
- `WORKWISE_CONFIG`: Custom config file path

### Configuration Structure

```yaml
ai:
  provider: "openai"           # LLM provider
  api_key: "sk-..."            # API key
  model: "gpt-4"               # Model name
  base_url: ""                 # Optional: for compatible APIs
  agent:
    max_iterations: 10         # Max agent iterations
    temperature: 0.7           # Sampling temperature (0.0-2.0)
    system_prompt: "..."       # System prompt for the agent
    history_enabled: true      # Enable conversation history
    max_history: 50            # Max messages to keep

cli:
  interactive: true            # Default to interactive mode
  prompt: "WorkWise> "         # CLI prompt string
  history_file: "~/.workwise_history"  # Command history file

extensions:
  mcp_enabled: false           # Enable MCP support
  mcp_servers: []              # MCP server configurations
  skills_enabled: false        # Enable skills framework
  skills_paths: []             # Paths to skill modules
  desktop_enabled: false       # Enable desktop integration
  desktop_hotkey: ""           # Hotkey for desktop window
  desktop_position: ""         # Window position
```

## Adding New Features

### Adding a New LLM Provider

1. Create provider implementation in `internal/llm/client.go`
2. Add provider case in `NewClient` function
3. Update configuration to support new provider
4. Add documentation and example

### Adding a New Skill

1. Implement `Skill` interface from `pkg/skills`
2. Register skill in skill registry
3. Document skill parameters and usage
4. Add tests for skill functionality

### Adding MCP Support

1. Implement MCP server interface in `pkg/mcp`
2. Add MCP protocol handling
3. Integrate with agent for tool calling
4. Update configuration and documentation

## Dependencies

### Direct Dependencies
- `github.com/cloudwego/eino`: AI agent framework
- `github.com/cloudwego/eino-ext`: LLM provider integrations
- `github.com/urfave/cli/v2`: CLI framework
- `gopkg.in/yaml.v3`: YAML configuration parsing

### Key Transitive Dependencies
- OpenAI Go SDK (via eino-ext)
- Various utility libraries

## Build System

The Makefile provides several targets:
- `build`: Build the binary
- `build-all`: Build for multiple platforms
- `install`: Install to GOPATH/bin
- `clean`: Clean build artifacts
- `test`: Run tests
- `test-coverage`: Run tests with coverage
- `lint`: Run linter
- `fmt`: Format code
- `run`: Build and run
- `deps`: Download dependencies
- `help`: Show available targets

## Future Enhancements

### Near-term
- [ ] Add unit tests
- [ ] Implement logging framework
- [ ] Add more LLM providers (Anthropic, etc.)
- [ ] Implement basic skills (file ops, web search)

### Medium-term
- [ ] MCP protocol implementation
- [ ] Desktop window integration (Windows/Mac)
- [ ] Plugin system for custom skills
- [ ] Multi-agent support

### Long-term
- [ ] GUI application
- [ ] Voice interaction
- [ ] Mobile companion app
- [ ] Cloud sync for preferences

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

MIT License - see LICENSE file for details
