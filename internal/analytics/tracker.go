package analytics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// BehaviorType represents the type of user behavior
type BehaviorType string

const (
	BehaviorQuery    BehaviorType = "query"
	BehaviorResponse BehaviorType = "response"
	BehaviorCommand  BehaviorType = "command"
	BehaviorSession  BehaviorType = "session"
	BehaviorError    BehaviorType = "error"
)

// BehaviorEvent represents a single user behavior event
type BehaviorEvent struct {
	ID        string       `json:"id"`
	Type      BehaviorType `json:"type"`
	Content   string       `json:"content,omitempty"`
	Metadata  Metadata     `json:"metadata"`
	Timestamp time.Time    `json:"timestamp"`
}

// Metadata contains additional information about the behavior
type Metadata struct {
	Duration     time.Duration `json:"duration,omitempty"`
	TokenCount   int           `json:"token_count,omitempty"`
	Model        string        `json:"model,omitempty"`
	Success      bool          `json:"success"`
	ErrorMessage string        `json:"error_message,omitempty"`
	SessionID    string        `json:"session_id,omitempty"`
	CommandName  string        `json:"command_name,omitempty"`
}

// Tracker tracks user behaviors
type Tracker struct {
	enabled      bool
	dataPath     string
	events       []BehaviorEvent
	sessionID    string
	sessionStart time.Time
	mu           sync.RWMutex
}

// NewTracker creates a new behavior tracker
func NewTracker(enabled bool, dataPath string) (*Tracker, error) {
	if !enabled {
		return &Tracker{enabled: false}, nil
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	sessionID := generateSessionID()

	tracker := &Tracker{
		enabled:      true,
		dataPath:     dataPath,
		events:       make([]BehaviorEvent, 0),
		sessionID:    sessionID,
		sessionStart: time.Now(),
	}

	// Load existing events from today's file
	if err := tracker.loadTodayEvents(); err != nil {
		// Log error but don't fail - start with empty events
		fmt.Fprintf(os.Stderr, "Warning: failed to load existing events: %v\n", err)
	}

	return tracker, nil
}

// TrackQuery tracks a user query
func (t *Tracker) TrackQuery(query string, sessionID string) {
	if !t.enabled {
		return
	}

	event := BehaviorEvent{
		ID:      generateEventID(),
		Type:    BehaviorQuery,
		Content: query,
		Metadata: Metadata{
			SessionID: sessionID,
			Success:   true,
		},
		Timestamp: time.Now(),
	}

	t.addEvent(event)
}

// TrackResponse tracks an AI response
func (t *Tracker) TrackResponse(response string, duration time.Duration, tokenCount int, model string, sessionID string) {
	if !t.enabled {
		return
	}

	event := BehaviorEvent{
		ID:      generateEventID(),
		Type:    BehaviorResponse,
		Content: response,
		Metadata: Metadata{
			Duration:   duration,
			TokenCount: tokenCount,
			Model:      model,
			SessionID:  sessionID,
			Success:    true,
		},
		Timestamp: time.Now(),
	}

	t.addEvent(event)
}

// TrackCommand tracks a command execution
func (t *Tracker) TrackCommand(commandName string, sessionID string) {
	if !t.enabled {
		return
	}

	event := BehaviorEvent{
		ID:   generateEventID(),
		Type: BehaviorCommand,
		Metadata: Metadata{
			CommandName: commandName,
			SessionID:   sessionID,
			Success:     true,
		},
		Timestamp: time.Now(),
	}

	t.addEvent(event)
}

// TrackError tracks an error event
func (t *Tracker) TrackError(errorMsg string, sessionID string) {
	if !t.enabled {
		return
	}

	event := BehaviorEvent{
		ID:   generateEventID(),
		Type: BehaviorError,
		Metadata: Metadata{
			ErrorMessage: errorMsg,
			SessionID:    sessionID,
			Success:      false,
		},
		Timestamp: time.Now(),
	}

	t.addEvent(event)
}

// TrackSessionEnd tracks session end
func (t *Tracker) TrackSessionEnd() {
	if !t.enabled {
		return
	}

	duration := time.Since(t.sessionStart)
	event := BehaviorEvent{
		ID:   generateEventID(),
		Type: BehaviorSession,
		Metadata: Metadata{
			Duration:  duration,
			SessionID: t.sessionID,
			Success:   true,
		},
		Timestamp: time.Now(),
	}

	t.addEvent(event)

	// Persist events to file
	if err := t.persist(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to persist events: %v\n", err)
	}
}

// addEvent adds an event to the tracker and persists periodically
func (t *Tracker) addEvent(event BehaviorEvent) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.events = append(t.events, event)

	// Persist every 10 events or if it's a session end
	if len(t.events)%10 == 0 || event.Type == BehaviorSession {
		go func() {
			if err := t.persist(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to persist events: %v\n", err)
			}
		}()
	}
}

// persist saves events to disk
func (t *Tracker) persist() error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if len(t.events) == 0 {
		return nil
	}

	// Get today's file path
	filePath := t.getTodayFilePath()

	// Read existing events
	var allEvents []BehaviorEvent
	if data, err := os.ReadFile(filePath); err == nil {
		if err := json.Unmarshal(data, &allEvents); err != nil {
			return fmt.Errorf("failed to unmarshal existing events: %w", err)
		}
	}

	// Append new events
	allEvents = append(allEvents, t.events...)

	// Write back to file
	data, err := json.MarshalIndent(allEvents, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write events file: %w", err)
	}

	// Clear in-memory events after successful persist
	t.events = make([]BehaviorEvent, 0)

	return nil
}

// loadTodayEvents loads events from today's file
func (t *Tracker) loadTodayEvents() error {
	filePath := t.getTodayFilePath()

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's okay
		}
		return fmt.Errorf("failed to read events file: %w", err)
	}

	var events []BehaviorEvent
	if err := json.Unmarshal(data, &events); err != nil {
		return fmt.Errorf("failed to unmarshal events: %w", err)
	}

	return nil
}

// getTodayFilePath returns the file path for today's events
func (t *Tracker) getTodayFilePath() string {
	today := time.Now().Format("2006-01-02")
	return filepath.Join(t.dataPath, fmt.Sprintf("events-%s.json", today))
}

// GetSessionID returns the current session ID
func (t *Tracker) GetSessionID() string {
	return t.sessionID
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	return fmt.Sprintf("session-%d", time.Now().UnixNano())
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("event-%d", time.Now().UnixNano())
}
