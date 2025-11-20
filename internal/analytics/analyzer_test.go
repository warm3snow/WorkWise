package analytics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAnalyzePeriod(t *testing.T) {
	tmpDir := t.TempDir()
	analyzer := NewAnalyzer(tmpDir)

	// Create some test events
	events := []BehaviorEvent{
		{
			ID:        "event-1",
			Type:      BehaviorQuery,
			Content:   "test query",
			Timestamp: time.Now(),
			Metadata:  Metadata{SessionID: "session-1", Success: true},
		},
		{
			ID:        "event-2",
			Type:      BehaviorResponse,
			Content:   "test response",
			Timestamp: time.Now(),
			Metadata: Metadata{
				SessionID:  "session-1",
				Success:    true,
				Duration:   100 * time.Millisecond,
				TokenCount: 50,
				Model:      "test-model",
			},
		},
		{
			ID:        "event-3",
			Type:      BehaviorCommand,
			Timestamp: time.Now(),
			Metadata: Metadata{
				SessionID:   "session-1",
				Success:     true,
				CommandName: "help",
			},
		},
	}

	// Write events to file
	today := time.Now().Format("2006-01-02")
	filePath := filepath.Join(tmpDir, "events-"+today+".json")
	data, _ := json.MarshalIndent(events, "", "  ")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		t.Fatalf("Failed to write test events: %v", err)
	}

	// Analyze
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now().Add(24 * time.Hour)
	stats, err := analyzer.AnalyzePeriod(startDate, endDate)
	if err != nil {
		t.Fatalf("AnalyzePeriod failed: %v", err)
	}

	// Verify statistics
	if stats.TotalQueries != 1 {
		t.Errorf("Expected 1 query, got %d", stats.TotalQueries)
	}
	if stats.TotalResponses != 1 {
		t.Errorf("Expected 1 response, got %d", stats.TotalResponses)
	}
	if stats.TotalCommands != 1 {
		t.Errorf("Expected 1 command, got %d", stats.TotalCommands)
	}
	if stats.TotalTokens != 50 {
		t.Errorf("Expected 50 tokens, got %d", stats.TotalTokens)
	}
	if stats.AvgResponseTime != 100*time.Millisecond {
		t.Errorf("Expected avg response time 100ms, got %v", stats.AvgResponseTime)
	}
	if stats.TopCommands["help"] != 1 {
		t.Errorf("Expected 'help' command count 1, got %d", stats.TopCommands["help"])
	}
}

func TestAnalyzeToday(t *testing.T) {
	tmpDir := t.TempDir()
	analyzer := NewAnalyzer(tmpDir)

	// Create empty statistics - should not error
	stats, err := analyzer.AnalyzeToday()
	if err != nil {
		t.Errorf("AnalyzeToday failed: %v", err)
	}

	if stats == nil {
		t.Error("Expected non-nil statistics")
	}
}

func TestAnalyzeEmptyPeriod(t *testing.T) {
	tmpDir := t.TempDir()
	analyzer := NewAnalyzer(tmpDir)

	startDate := time.Now().Add(-7 * 24 * time.Hour)
	endDate := time.Now()

	stats, err := analyzer.AnalyzePeriod(startDate, endDate)
	if err != nil {
		t.Errorf("AnalyzePeriod should not fail on empty data: %v", err)
	}

	if stats.TotalQueries != 0 {
		t.Errorf("Expected 0 queries, got %d", stats.TotalQueries)
	}
}

func TestCalculateProductivityScore(t *testing.T) {
	analyzer := NewAnalyzer("")

	// Test high productivity
	stats := &Statistics{
		TotalQueries: 100,
		TotalErrors:  0,
		DailyActivity: map[string]int{
			"2024-01-01": 10,
			"2024-01-02": 10,
			"2024-01-03": 10,
		},
	}
	score := analyzer.calculateProductivityScore(stats)
	if score < 0 || score > 100 {
		t.Errorf("Score should be between 0-100, got %f", score)
	}

	// Test with errors
	stats.TotalErrors = 50
	scoreWithErrors := analyzer.calculateProductivityScore(stats)
	if scoreWithErrors >= score {
		t.Error("Score with errors should be lower")
	}

	// Test score normalization
	stats.TotalQueries = 1000
	normalizedScore := analyzer.calculateProductivityScore(stats)
	if normalizedScore > 100 {
		t.Errorf("Score should not exceed 100, got %f", normalizedScore)
	}
}

func TestGetInsights(t *testing.T) {
	stats := &Statistics{
		Period:            "2024-01-01 to 2024-01-07",
		TotalQueries:      50,
		TotalErrors:       5,
		ErrorRate:         0.1,
		MostActiveHour:    14,
		AvgSessionTime:    45 * time.Minute,
		ProductivityScore: 75.0,
		TopCommands: map[string]int{
			"help": 10,
			"ask":  5,
		},
	}

	insights := GetInsights(stats)

	if len(insights) == 0 {
		t.Error("Expected insights, got none")
	}

	// Check that insights contain meaningful information
	found := false
	for _, insight := range insights {
		if len(insight) > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected non-empty insights")
	}
}

func TestFormatStatistics(t *testing.T) {
	stats := &Statistics{
		Period:            "2024-01-01 to 2024-01-07",
		TotalQueries:      50,
		TotalResponses:    45,
		TotalCommands:     20,
		TotalSessions:     5,
		TotalErrors:       2,
		ErrorRate:         0.04,
		AvgResponseTime:   150 * time.Millisecond,
		AvgSessionTime:    30 * time.Minute,
		TotalTokens:       5000,
		MostActiveHour:    14,
		ProductivityScore: 75.0,
		TopCommands: map[string]int{
			"help": 10,
			"ask":  5,
		},
		DailyActivity: map[string]int{
			"2024-01-01": 10,
			"2024-01-02": 15,
		},
	}

	formatted := FormatStatistics(stats)

	if len(formatted) == 0 {
		t.Error("Expected formatted output, got empty string")
	}

	// Check that formatted output contains key information
	if !contains(formatted, "Total Queries") {
		t.Error("Formatted output should contain 'Total Queries'")
	}
	if !contains(formatted, "50") {
		t.Error("Formatted output should contain the query count")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || contains(s[1:], substr)))
}
