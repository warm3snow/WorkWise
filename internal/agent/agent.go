package agent

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/warm3snow/WorkWise/internal/config"
	"github.com/warm3snow/WorkWise/internal/llm"
)

// Agent represents the AI agent using eino framework
type Agent struct {
	config  *config.Config
	llm     model.ChatModel
	history []*schema.Message
}

// NewAgent creates a new AI agent
func NewAgent(cfg *config.Config) (*Agent, error) {
	// Initialize LLM client using eino-ext
	llmClient, err := llm.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM client: %w", err)
	}

	agent := &Agent{
		config:  cfg,
		llm:     llmClient,
		history: make([]*schema.Message, 0),
	}

	// Add system prompt to history if configured
	if cfg.AI.Agent.SystemPrompt != "" {
		agent.history = append(agent.history, schema.SystemMessage(cfg.AI.Agent.SystemPrompt))
	}

	return agent, nil
}

// Process processes a user message and returns the agent's response
func (a *Agent) Process(ctx context.Context, userMessage string) (string, error) {
	// Add user message to history
	a.addToHistory(schema.UserMessage(userMessage))

	// Prepare messages for the LLM
	messages := a.getRecentHistory()

	// Call the LLM
	resp, err := a.llm.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	// Extract response content
	responseContent := resp.Content

	// Add assistant response to history
	a.addToHistory(schema.AssistantMessage(responseContent, nil))

	return responseContent, nil
}

// addToHistory adds a message to the conversation history
func (a *Agent) addToHistory(msg *schema.Message) {
	if !a.config.AI.Agent.HistoryEnabled {
		return
	}

	a.history = append(a.history, msg)

	// Trim history if it exceeds max size
	maxHistory := a.config.AI.Agent.MaxHistory
	if maxHistory > 0 && len(a.history) > maxHistory {
		// Keep system message if present
		systemMsgCount := 0
		if len(a.history) > 0 && a.history[0].Role == schema.System {
			systemMsgCount = 1
		}

		// Remove oldest messages while keeping system message
		excessCount := len(a.history) - maxHistory
		if excessCount > 0 {
			a.history = append(
				a.history[:systemMsgCount],
				a.history[systemMsgCount+excessCount:]...,
			)
		}
	}
}

// getRecentHistory returns recent conversation history
func (a *Agent) getRecentHistory() []*schema.Message {
	if !a.config.AI.Agent.HistoryEnabled {
		// Return only the last user message with system prompt
		if len(a.history) > 0 {
			if a.history[0].Role == schema.System {
				if len(a.history) > 1 {
					return []*schema.Message{a.history[0], a.history[len(a.history)-1]}
				}
				return []*schema.Message{a.history[0]}
			}
			return []*schema.Message{a.history[len(a.history)-1]}
		}
		return a.history
	}

	return a.history
}

// ClearHistory clears the conversation history
func (a *Agent) ClearHistory() {
	var systemMsg *schema.Message
	hasSystem := false

	// Preserve system message if present
	if len(a.history) > 0 && a.history[0].Role == schema.System {
		systemMsg = a.history[0]
		hasSystem = true
	}

	a.history = make([]*schema.Message, 0)

	if hasSystem {
		a.history = append(a.history, systemMsg)
	}
}
