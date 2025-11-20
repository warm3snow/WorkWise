# WorkWise Quick Start Guide

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/warm3snow/WorkWise.git
cd WorkWise

# Build the application
make build

# Optional: Move to your PATH
sudo cp workwise /usr/local/bin/
# or
export PATH=$PATH:$(pwd)
```

### Option 2: Using Go Install

```bash
go install github.com/warm3snow/WorkWise/cmd/workwise@latest
```

## Configuration

### Quick Setup

1. Set your API key:
```bash
export WORKWISE_API_KEY="your-api-key-here"
```

2. (Optional) Customize settings:
```bash
export WORKWISE_PROVIDER="openai"
export WORKWISE_MODEL="gpt-4"
```

### Configuration File

Create `~/.workwise/config.yaml`:

```yaml
ai:
  provider: openai
  api_key: your-api-key-here
  model: gpt-4
  agent:
    temperature: 0.7
    system_prompt: "You are WorkWise, an intelligent desktop assistant."
```

Or use:
```bash
workwise config init  # Creates default config
```

## Usage Examples

### Interactive Mode (Recommended)

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
WorkWise> What's the weather like today?
[AI response...]

WorkWise> Help me write a Python function to sort a list
[AI response...]

WorkWise> exit
Goodbye!
```

### Single Question Mode

```bash
# Quick question
workwise ask "How do I list files in Linux?"

# Or without the 'ask' command
workwise "What is the capital of France?"
```

### Configuration Management

```bash
# Show current configuration
workwise config show

# Create default configuration file
workwise config init

# Show version
workwise version
```

## Interactive Commands

When in interactive mode, you can use:

- **help** or **?** - Show help message
- **clear** or **cls** - Clear the screen
- **exit**, **quit**, or **q** - Exit the application

## Tips

1. **Long conversations**: WorkWise maintains conversation history, so you can refer to previous messages.

2. **Custom prompts**: Edit `~/.workwise/config.yaml` to customize the system prompt:
   ```yaml
   ai:
     agent:
       system_prompt: "You are a coding expert specializing in Go."
   ```

3. **OpenAI-compatible APIs**: Set a custom base URL:
   ```yaml
   ai:
     base_url: "https://your-api-endpoint.com/v1"
   ```

4. **Adjust temperature**: Control response creativity (0.0 = focused, 2.0 = creative):
   ```yaml
   ai:
     agent:
       temperature: 0.7
   ```

## Troubleshooting

### "API key is required" error
Set your API key:
```bash
export WORKWISE_API_KEY="sk-..."
```

### Connection errors
Check your internet connection and API endpoint accessibility.

### Rate limiting
Some providers have rate limits. Consider adjusting the model or upgrading your API plan.

## Next Steps

- Read the full [README](../README.md) for detailed features
- Check [DEVELOPMENT.md](DEVELOPMENT.md) for architecture and extending WorkWise
- Explore the example config at [config.example.yaml](../config.example.yaml)

## Getting Help

- Run `workwise --help` for command-line options
- Type `help` in interactive mode
- Check the documentation in the `docs/` directory
- Open an issue on GitHub for bugs or feature requests

---

Happy working with WorkWise! 🚀
