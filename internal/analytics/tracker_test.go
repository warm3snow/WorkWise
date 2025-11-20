package analytics

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewTracker(t *testing.T) {
	// Test disabled tracker
	tracker, err := NewTracker(false, "")
	if err != nil {
		t.Errorf("NewTracker with disabled should not error: %v", err)
	}
	if tracker.enabled {
		t.Error("Disabled tracker should have enabled=false")
	}

	// Test enabled tracker
	tmpDir := t.TempDir()
	tracker, err = NewTracker(true, tmpDir)
	if err != nil {
		t.Errorf("NewTracker with valid path should not error: %v", err)
	}
	if !tracker.enabled {
		t.Error("Enabled tracker should have enabled=true")
	}
	if tracker.sessionID == "" {
		t.Error("Session ID should not be empty")
	}
}

func TestTrackQuery(t *testing.T) {
	tmpDir := t.TempDir()
	tracker, err := NewTracker(true, tmpDir)
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	sessionID := "test-session"
	query := "test query"
	
	tracker.TrackQuery(query, sessionID)
	
	if len(tracker.events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(tracker.events))
	}
	
	event := tracker.events[0]
	if event.Type != BehaviorQuery {
		t.Errorf("Expected type %s, got %s", BehaviorQuery, event.Type)
	}
	if event.Content != query {
		t.Errorf("Expected content %s, got %s", query, event.Content)
	}
	if event.Metadata.SessionID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, event.Metadata.SessionID)
	}
}

func TestTrackResponse(t *testing.T) {
	tmpDir := t.TempDir()
	tracker, err := NewTracker(true, tmpDir)
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	sessionID := "test-session"
	response := "test response"
	duration := 100 * time.Millisecond
	tokenCount := 50
	model := "test-model"
	
	tracker.TrackResponse(response, duration, tokenCount, model, sessionID)
	
	if len(tracker.events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(tracker.events))
	}
	
	event := tracker.events[0]
	if event.Type != BehaviorResponse {
		t.Errorf("Expected type %s, got %s", BehaviorResponse, event.Type)
	}
	if event.Metadata.Duration != duration {
		t.Errorf("Expected duration %v, got %v", duration, event.Metadata.Duration)
	}
	if event.Metadata.TokenCount != tokenCount {
		t.Errorf("Expected token count %d, got %d", tokenCount, event.Metadata.TokenCount)
	}
	if event.Metadata.Model != model {
		t.Errorf("Expected model %s, got %s", model, event.Metadata.Model)
	}
}

func TestPersist(t *testing.T) {
	tmpDir := t.TempDir()
	tracker, err := NewTracker(true, tmpDir)
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	// Add some events
	tracker.TrackQuery("query 1", "session-1")
	tracker.TrackQuery("query 2", "session-1")
	
	// Persist
	err = tracker.persist()
	if err != nil {
		t.Errorf("Persist failed: %v", err)
	}
	
	// Check file exists
	filePath := tracker.getTodayFilePath()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Events file was not created")
	}
	
	// Events should be cleared after persist
	if len(tracker.events) != 0 {
		t.Errorf("Events should be cleared after persist, got %d", len(tracker.events))
	}
}

func TestTrackerDisabled(t *testing.T) {
	tracker, _ := NewTracker(false, "")
	
	// These should not panic or error when disabled
	tracker.TrackQuery("test", "session")
	tracker.TrackResponse("test", 0, 0, "", "session")
	tracker.TrackCommand("test", "session")
	tracker.TrackError("test", "session")
	tracker.TrackSessionEnd()
	
	if len(tracker.events) != 0 {
		t.Error("Disabled tracker should not track events")
	}
}

func TestGetTodayFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	tracker, _ := NewTracker(true, tmpDir)
	
	filePath := tracker.getTodayFilePath()
	today := time.Now().Format("2006-01-02")
	expected := filepath.Join(tmpDir, "events-"+today+".json")
	
	if filePath != expected {
		t.Errorf("Expected file path %s, got %s", expected, filePath)
	}
}
