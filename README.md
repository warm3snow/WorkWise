# WorkWise

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

WorkWise is a flagship-level AIPC (AI-Powered Computer) application designed to be your intelligent desktop assistant. It brings the power of AI agents directly to your workspace, helping you accomplish tasks more efficiently through natural language interaction.

## 🌟 Features

- **🤖 AI-Powered Assistant**: Leverages advanced language models to understand and assist with your tasks
- **💬 Interactive CLI**: Command-line interface with conversational interaction
- **🔧 Extensible Architecture**: Built with extensibility in mind for future enhancements
- **🎯 Agent Framework**: Uses CloudWeGo Eino for robust agent orchestration
- **🔌 LLM Integration**: Integrates with various LLM providers via CloudWeGo Eino-Ext
- **📝 Conversation History**: Maintains context across conversations
- **🎨 Cross-Platform Ready**: Designed for Windows and macOS support

## 🚀 Future Roadmap

- **MCP Support**: Integration with Anthropic's Model Context Protocol
- **Skills Framework**: ✅ Implemented - Anthropic-style skills loader for extending agent capabilities
- **Desktop Integration**: Hotkey-activated overlay window on Windows/Mac
- **Multi-Provider Support**: Additional LLM providers (Anthropic Claude, etc.)

## 📋 Prerequisites

- Go 1.21 or higher
- API key for your chosen LLM provider (OpenAI, etc.)

## 🛠️ Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/warm3snow/WorkWise.git
cd WorkWise

# Build the application
go build -o workwise ./cmd/workwise

# Optional: Install to your system
go install ./cmd/workwise
```

### Using Go Install

```bash
go install github.com/warm3snow/WorkWise/cmd/workwise@latest
```

## ⚙️ Configuration

WorkWise can be configured through environment variables or a configuration file.

### Environment Variables

```bash
export WORKWISE_API_KEY="your-api-key-here"
export WORKWISE_PROVIDER="openai"
export WORKWISE_MODEL="gpt-4"
```

For Ollama (local models):
```bash
export WORKWISE_PROVIDER="ollama"
export WORKWISE_MODEL="llama3"
# Optional: customize Ollama base URL (defaults to http://localhost:11434)
# export WORKWISE_BASE_URL="http://localhost:11434"
```

### Configuration File

Create a configuration file at `~/.workwise/config.yaml`:

For OpenAI:
```yaml
ai:
  provider: openai
  api_key: your-api-key-here
  model: gpt-4
  base_url: ""  # Optional: for compatible APIs
  agent:
    max_iterations: 10
    temperature: 0.7
    system_prompt: "You are WorkWise, an intelligent desktop assistant. Help users with their tasks efficiently and professionally."
    history_enabled: true
    max_history: 50

cli:
  interactive: true
  prompt: "WorkWise> "
  history_file: ~/.workwise_history

extensions:
  mcp_enabled: false
  mcp_servers: []
  skills_enabled: false
  skills_paths: []
  desktop_enabled: false
  desktop_hotkey: ""
  desktop_position: ""
```

For Ollama (local models):
```yaml
ai:
  provider: ollama
  model: llama3  # or mistral, codellama, etc.
  base_url: ""  # Optional: defaults to http://localhost:11434
  agent:
    max_iterations: 10
    temperature: 0.7
    system_prompt: "You are WorkWise, an intelligent desktop assistant. Help users with their tasks efficiently and professionally."
    history_enabled: true
    max_history: 50

cli:
  interactive: true
  prompt: "WorkWise> "
  history_file: ~/.workwise_history

extensions:
  mcp_enabled: false
  mcp_servers: []
  skills_enabled: false
  skills_paths: []
  desktop_enabled: false
  desktop_hotkey: ""
  desktop_position: ""
```

Generate a default configuration file:

```bash
workwise config init
```

## 📖 Usage

### Interactive Mode

Start an interactive chat session:

```bash
workwise
# or
workwise chat
```

Example session:
```
WorkWise - Intelligent Desktop Assistant
Version: dev
Type 'help' for commands, 'exit' or 'quit' to leave
---
WorkWise> What can you help me with?
I can assist you with various tasks such as...
```

### Single Question Mode

Ask a single question and get a response:

```bash
workwise ask "How do I create a directory in Linux?"
# or simply
workwise "How do I create a directory in Linux?"
```

### Configuration Management

```bash
# Show current configuration
workwise config show

# Initialize configuration file
workwise config init

# Show version information
workwise version
```

### Available Commands in Interactive Mode

- `help` or `?` - Show help message
- `clear` or `cls` - Clear the screen
- `exit`, `quit`, or `q` - Exit the application

### Using Skills

WorkWise supports Anthropic-style skills for extending agent capabilities:

1. **Enable skills** in your configuration:
   ```yaml
   extensions:
     skills_enabled: true
     skills_paths:
       - "~/.workwise/skills"
       - "./examples/skills"
   ```

2. **Create custom skills** following the [Anthropic Skills Spec](https://github.com/anthropics/skills):
   ```bash
   mkdir -p ~/.workwise/skills/my-skill
   # Create SKILL.md with YAML frontmatter + instructions
   ```

3. **Use example skills** from `examples/skills/` directory

Skills are automatically discovered and loaded, providing specialized instructions and workflows to the agent.

For more details, see [examples/skills/README.md](examples/skills/README.md) and the [Architecture documentation](docs/ARCHITECTURE.md).

## 🏗️ Architecture

WorkWise is built with a modular architecture:

```
WorkWise/
├── cmd/
│   └── workwise/          # Main application entry point
├── internal/
│   ├── agent/             # AI agent implementation using Eino
│   ├── cli/               # Command-line interface
│   ├── config/            # Configuration management
│   ├── llm/               # LLM client using Eino-Ext
│   └── extensions/        # Future extensions
├── pkg/
│   ├── mcp/               # Model Context Protocol support (future)
│   └── skills/            # Skills framework (Anthropic-style)
├── examples/
│   └── skills/            # Example skills
└── README.md
```

### Key Components

- **Agent**: Uses CloudWeGo Eino framework for agent orchestration
- **LLM Client**: Integrates with LLM providers via CloudWeGo Eino-Ext
- **CLI**: User-friendly command-line interface with interactive mode
- **Skills**: Anthropic-style skills loader for extending agent capabilities with specialized instructions
- **Extensions**: Pluggable architecture for MCP, skills, and desktop integration

## 🔌 LLM Providers

Currently supported providers:

- **OpenAI**: GPT-4, GPT-3.5-turbo, and compatible APIs
- **Ollama**: Local LLM models (Llama, Mistral, etc.)

Future providers:
- Anthropic Claude (via Eino-Ext)
- Other providers as they become available in Eino-Ext

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [CloudWeGo Eino](https://github.com/cloudwego/eino) - AI agent framework
- [CloudWeGo Eino-Ext](https://github.com/cloudwego/eino-ext) - LLM integrations
- All contributors and users of WorkWise

## 📞 Support

For issues, questions, or suggestions, please open an issue on GitHub.

---

Made with ❤️ by the WorkWise team