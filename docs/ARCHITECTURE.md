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

#### Current Architecture

The `pkg/skills` package provides a filesystem-based skills loader following the Anthropic Agent Skills Spec:

```go
type Skill struct {
    // Metadata from YAML frontmatter
    Name        string
    Description string
    License     string
    AllowedTools []string
    Metadata    map[string]string
    
    // Content from markdown body
    Instructions string
    
    // Path information
    SkillPath string
}
```

**Reference**: [Anthropic Skills Specification](https://github.com/anthropics/skills)

#### Implementation Strategy

Skills are folder-based with a `SKILL.md` file containing YAML frontmatter and markdown instructions:

1. **Skill Structure**:
   ```
   skill-name/
   ├── SKILL.md          # Required: Frontmatter + instructions
   ├── scripts/          # Optional: Executable scripts
   │   └── helper.py
   └── resources/        # Optional: Additional files
       └── template.txt
   ```

2. **SKILL.md Format**:
   ```markdown
   ---
   name: skill-name
   description: Clear description of what this skill does and when to use it
   license: MIT
   metadata:
     version: "1.0"
   ---
   
   # Skill Name
   
   [Instructions that the agent will follow when this skill is active]
   
   ## Usage
   - Example 1
   - Example 2
   ```

3. **Configuration**: Already supported:
   ```yaml
   extensions:
     skills_enabled: true
     skills_paths:
       - "./skills"
       - "~/.workwise/skills"
   ```

4. **Loading Skills**:
   ```go
   loader := skills.NewLoader([]string{"./skills", "~/.workwise/skills"})
   err := loader.LoadAll()
   
   // Get a specific skill
   skill, err := loader.Get("skill-name")
   
   // Access skill instructions
   instructions := skill.Instructions
   
   // Check for scripts
   if skill.HasScript("helper.py") {
       scriptPath := skill.GetScriptPath("helper.py")
   }
   ```

#### Example Skills

Skills provide instructions and workflows rather than executable code:

**Document Creation Skill**:
- Instructions for creating structured documents
- Templates and formatting guidelines
- References to helper scripts for PDF/DOCX generation

**Data Analysis Skill**:
- Step-by-step data processing workflows
- Python scripts for analysis tasks
- Guidelines for visualization

**Code Review Skill**:
- Checklist for code quality review
- Security scanning procedures
- Best practices references

#### Integration with Agent

The agent should:
1. Load all skills from configured paths on startup
2. Use skill descriptions to determine which skill(s) are relevant to the user's request
3. When a skill is selected, include its instructions in the context
4. Follow the instructions provided in the skill's markdown content
5. Execute any scripts referenced by the skill if needed
6. Return results based on the skill's workflow

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
1. ✅ Implement Anthropic-style skills loader
2. Create example skills following SKILL.md format
3. Add skill discovery mechanism to agent
4. Integrate skill instructions into agent context

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

Create a new skill directory with a SKILL.md file:

```bash
mkdir -p ~/.workwise/skills/calculator-skill
```

Create `~/.workwise/skills/calculator-skill/SKILL.md`:

```markdown
---
name: calculator-skill
description: Perform basic arithmetic operations. Use this skill when the user needs to calculate mathematical expressions or perform arithmetic.
license: MIT
---

# Calculator Skill

This skill provides instructions for performing arithmetic calculations.

## When to Use

Use this skill when the user requests:
- Basic arithmetic (addition, subtraction, multiplication, division)
- Mathematical calculations
- Numeric operations

## Instructions

When performing calculations:

1. Parse the mathematical expression from the user's request
2. Identify the operation and operands
3. Perform the calculation
4. Return the result with proper formatting

## Examples

### Example 1: Simple Addition
```
User: "What is 25 plus 17?"
Response: "25 + 17 = 42"
```

### Example 2: Division
```
User: "Divide 144 by 12"
Response: "144 ÷ 12 = 12"
```

## Guidelines

- Handle division by zero gracefully
- Format results appropriately (decimals, fractions, etc.)
- Show the calculation steps if helpful
- Support multiple operations in sequence
```

The skill will be automatically discovered and loaded by WorkWise when the skills system is enabled in the configuration.

## Conclusion

WorkWise's architecture is designed for extensibility from the ground up. The foundation is in place for:
- MCP integration with Anthropic
- Custom skills development
- Desktop UI integration
- Multi-provider LLM support

Each extension point has clear interfaces and examples, making it straightforward to add new capabilities without modifying core functionality.
