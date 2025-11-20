# WorkWise Architecture: Extensibility Design

## Overview

WorkWise is designed with extensibility as a core principle. This document explains how the architecture supports future enhancements, particularly for Anthropic's Model Context Protocol (MCP) and the Skills framework.

## Extensibility Points

### 1. Model Context Protocol (MCP) Integration

#### Current Architecture

The `pkg/mcp` package provides the foundation for MCP support:

```go
type Server interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    ListTools(ctx context.Context) ([]Tool, error)
    CallTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error)
}
```

#### Implementation Strategy

When implementing MCP support:

1. **Server Implementation**: Create concrete implementations of the `Server` interface for different MCP servers
   ```go
   type AnthropicMCPServer struct {
       url    string
       client *http.Client
   }
   ```

2. **Integration with Agent**: Modify `internal/agent/agent.go` to:
   - Register MCP servers during initialization
   - Detect when to use MCP tools based on user requests
   - Convert MCP tool calls to Eino framework tool calls

3. **Configuration**: Already supported in `config.yaml`:
   ```yaml
   extensions:
     mcp_enabled: true
     mcp_servers:
       - name: "filesystem"
         url: "http://localhost:3000"
       - name: "browser"
         url: "http://localhost:3001"
   ```

4. **Tool Discovery**: MCP servers expose tools dynamically. The agent should:
   - Discover tools from all connected MCP servers on startup
   - Convert MCP tool schemas to Eino's tool format
   - Handle tool execution by routing to appropriate MCP server

#### Example Usage Flow

```
User Request → Agent Processes → Determines Tool Needed → 
Queries MCP Manager → Calls Appropriate MCP Server → 
Returns Result → Agent Continues Processing
```

### 2. Skills Framework

#### Current Architecture

The `pkg/skills` package provides a registry-based system:

```go
type Skill interface {
    Name() string
    Description() string
    Parameters() map[string]interface{}
    Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}
```

#### Implementation Strategy

Skills can be implemented as:

1. **Built-in Skills**: Compiled into the binary
   ```go
   type FileReadSkill struct {
       *BaseSkill
   }
   
   func (s *FileReadSkill) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
       path := params["path"].(string)
       return os.ReadFile(path)
   }
   ```

2. **Plugin Skills**: Loaded dynamically from separate binaries using Go plugins
   ```go
   plugin, err := plugin.Open("skills/custom_skill.so")
   skillFunc, err := plugin.Lookup("NewSkill")
   skill := skillFunc.(func() skills.Skill)()
   ```

3. **Configuration**: Already supported:
   ```yaml
   extensions:
     skills_enabled: true
     skills_paths:
       - "./skills"
       - "~/.workwise/skills"
   ```

#### Example Built-in Skills

**File Operations Skill**:
- Read files
- Write files
- List directories
- Search files

**Web Search Skill**:
- Search the web
- Fetch webpage content
- Extract structured data

**Code Execution Skill**:
- Execute code snippets safely
- Support multiple languages
- Return output and errors

#### Integration with Agent

The agent should:
1. Load and register all available skills on startup
2. Include skill descriptions in the system prompt or tool list
3. When the LLM requests a skill, validate and execute it
4. Return skill results for LLM to process

### 3. Desktop Integration

#### Architecture Plan

For Windows and macOS desktop integration:

1. **Platform-Specific Code**:
   ```go
   // internal/desktop/desktop.go
   type DesktopManager interface {
       ShowWindow() error
       HideWindow() error
       RegisterHotkey(hotkey string) error
   }
   
   // internal/desktop/desktop_windows.go
   // +build windows
   type WindowsDesktopManager struct { ... }
   
   // internal/desktop/desktop_darwin.go
   // +build darwin
   type DarwinDesktopManager struct { ... }
   ```

2. **UI Framework Options**:
   - **Webview**: Embed a lightweight web view (lorca, webview)
   - **Native**: Use platform-specific UI (walk for Windows, cocoa for macOS)
   - **Cross-platform**: Fyne or Wails

3. **Communication**: Use channels or IPC to communicate between CLI and desktop UI

4. **Configuration**:
   ```yaml
   extensions:
     desktop_enabled: true
     desktop_hotkey: "Ctrl+Alt+W"
     desktop_position: "top-right"
     desktop_theme: "dark"
   ```

### 4. Multi-Provider LLM Support

#### Adding New Providers

To add support for Anthropic Claude or other providers:

1. **Check Eino-Ext Support**: First check if the provider is supported by eino-ext
   ```go
   import "github.com/cloudwego/eino-ext/components/model/anthropic"
   ```

2. **Implement Provider Client** in `internal/llm/client.go`:
   ```go
   func newAnthropicClient(cfg *config.Config) (model.ChatModel, error) {
       clientConfig := &anthropic.ChatModelConfig{
           APIKey: cfg.AI.APIKey,
           Model:  cfg.AI.Model,
       }
       return anthropic.NewChatModel(context.Background(), clientConfig)
   }
   ```

3. **Add to Switch Statement**:
   ```go
   switch cfg.AI.Provider {
   case "openai":
       return newOpenAIClient(cfg)
   case "anthropic":
       return newAnthropicClient(cfg)
   default:
       return nil, fmt.Errorf("unsupported provider: %s", cfg.AI.Provider)
   }
   ```

## Design Patterns

### 1. Interface-Based Design

All major components are interface-based, allowing for:
- Easy mocking in tests
- Multiple implementations
- Dependency injection
- Loose coupling

### 2. Configuration-Driven

Features can be enabled/disabled via configuration:
- No code changes required
- Easy experimentation
- User control

### 3. Plugin Architecture

Both MCP and Skills use a plugin-like approach:
- Register new capabilities at runtime
- No core code modification needed
- Easy to extend

### 4. Context Propagation

All operations use `context.Context` for:
- Cancellation support
- Timeout handling
- Tracing and logging

## Migration Path

### Phase 1: Current State (Complete)
- ✅ Basic CLI interface
- ✅ OpenAI integration
- ✅ Configuration system
- ✅ Extension placeholders

### Phase 2: Skills Implementation
1. Implement basic built-in skills
2. Add skill registry to agent
3. Create skill discovery mechanism
4. Add skill execution to agent loop

### Phase 3: MCP Integration
1. Implement MCP protocol client
2. Add MCP server manager
3. Integrate with agent for tool calling
4. Test with Anthropic's MCP servers

### Phase 4: Desktop Integration
1. Choose UI framework
2. Implement platform-specific managers
3. Add hotkey registration
4. Create overlay window UI

### Phase 5: Advanced Features
1. Multi-agent support
2. Voice interaction
3. Advanced tool use
4. Cloud synchronization

## Best Practices for Extension Development

1. **Follow Existing Patterns**: Look at current implementations
2. **Use Interfaces**: Define clear contracts
3. **Error Handling**: Always return meaningful errors
4. **Context Support**: Use context for cancellation
5. **Configuration**: Make features configurable
6. **Documentation**: Document all public APIs
7. **Testing**: Write unit tests for new functionality

## Example: Adding a Custom Skill

```go
package skills

import "context"

type CalculatorSkill struct {
    *BaseSkill
}

func NewCalculatorSkill() *CalculatorSkill {
    return &CalculatorSkill{
        BaseSkill: NewBaseSkill(
            "calculator",
            "Perform basic arithmetic operations",
            map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "operation": map[string]string{"type": "string"},
                    "a":         map[string]string{"type": "number"},
                    "b":         map[string]string{"type": "number"},
                },
                "required": []string{"operation", "a", "b"},
            },
        ),
    }
}

func (s *CalculatorSkill) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    op := params["operation"].(string)
    a := params["a"].(float64)
    b := params["b"].(float64)
    
    switch op {
    case "add":
        return a + b, nil
    case "subtract":
        return a - b, nil
    case "multiply":
        return a * b, nil
    case "divide":
        if b == 0 {
            return nil, fmt.Errorf("division by zero")
        }
        return a / b, nil
    default:
        return nil, fmt.Errorf("unknown operation: %s", op)
    }
}
```

## Conclusion

WorkWise's architecture is designed for extensibility from the ground up. The foundation is in place for:
- MCP integration with Anthropic
- Custom skills development
- Desktop UI integration
- Multi-provider LLM support

Each extension point has clear interfaces and examples, making it straightforward to add new capabilities without modifying core functionality.
