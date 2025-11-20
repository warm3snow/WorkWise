package mcp

import (
	"context"
	"fmt"
)

// Server represents an MCP (Model Context Protocol) server
// This is a placeholder for future Anthropic MCP integration
type Server interface {
	// Connect establishes connection to the MCP server
	Connect(ctx context.Context) error

	// Disconnect closes connection to the MCP server
	Disconnect(ctx context.Context) error

	// ListTools lists available tools from the MCP server
	ListTools(ctx context.Context) ([]Tool, error)

	// CallTool invokes a tool on the MCP server
	CallTool(ctx context.Context, name string, params map[string]interface{}) (interface{}, error)
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Manager manages multiple MCP servers
type Manager struct {
	servers map[string]Server
}

// NewManager creates a new MCP manager
func NewManager() *Manager {
	return &Manager{
		servers: make(map[string]Server),
	}
}

// RegisterServer registers an MCP server
func (m *Manager) RegisterServer(name string, server Server) error {
	if _, exists := m.servers[name]; exists {
		return fmt.Errorf("server %s already registered", name)
	}
	m.servers[name] = server
	return nil
}

// GetServer retrieves a registered MCP server
func (m *Manager) GetServer(name string) (Server, error) {
	server, exists := m.servers[name]
	if !exists {
		return nil, fmt.Errorf("server %s not found", name)
	}
	return server, nil
}

// ListServers returns all registered server names
func (m *Manager) ListServers() []string {
	names := make([]string, 0, len(m.servers))
	for name := range m.servers {
		names = append(names, name)
	}
	return names
}

// ConnectAll connects to all registered servers
func (m *Manager) ConnectAll(ctx context.Context) error {
	for name, server := range m.servers {
		if err := server.Connect(ctx); err != nil {
			return fmt.Errorf("failed to connect to server %s: %w", name, err)
		}
	}
	return nil
}

// DisconnectAll disconnects from all registered servers
func (m *Manager) DisconnectAll(ctx context.Context) error {
	var lastErr error
	for _, server := range m.servers {
		if err := server.Disconnect(ctx); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Note: This is a foundational structure for future MCP support.
// Actual implementation will be added when integrating with Anthropic's MCP specification.
