package analytics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Statistics represents aggregated behavior statistics
type Statistics struct {
	Period            string         `json:"period"`
	TotalQueries      int            `json:"total_queries"`
	TotalResponses    int            `json:"total_responses"`
	TotalCommands     int            `json:"total_commands"`
	TotalErrors       int            `json:"total_errors"`
	TotalSessions     int            `json:"total_sessions"`
	AvgResponseTime   time.Duration  `json:"avg_response_time"`
	AvgSessionTime    time.Duration  `json:"avg_session_time"`
	TotalTokens       int            `json:"total_tokens"`
	TopCommands       map[string]int `json:"top_commands"`
	HourlyActivity    map[int]int    `json:"hourly_activity"`
	DailyActivity     map[string]int `json:"daily_activity"`
	ErrorRate         float64        `json:"error_rate"`
	MostActiveHour    int            `json:"most_active_hour"`
	ProductivityScore float64        `json:"productivity_score"`
}

// Analyzer analyzes user behavior data
type Analyzer struct {
	dataPath string
}

// NewAnalyzer creates a new behavior analyzer
func NewAnalyzer(dataPath string) *Analyzer {
	return &Analyzer{
		dataPath: dataPath,
	}
}

// AnalyzePeriod analyzes behaviors for a specific period
func (a *Analyzer) AnalyzePeriod(startDate, endDate time.Time) (*Statistics, error) {
	// Load all events in the period
	events, err := a.loadEventsByPeriod(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to load events: %w", err)
	}

	if len(events) == 0 {
		return &Statistics{
			Period: fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		}, nil
	}

	stats := &Statistics{
		Period:         fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TopCommands:    make(map[string]int),
		HourlyActivity: make(map[int]int),
		DailyActivity:  make(map[string]int),
	}

	var totalResponseTime time.Duration
	var totalSessionTime time.Duration
	var responseCount int
	var sessionCount int

	// Analyze each event
	for _, event := range events {
		switch event.Type {
		case BehaviorQuery:
			stats.TotalQueries++
		case BehaviorResponse:
			stats.TotalResponses++
			if event.Metadata.Duration > 0 {
				totalResponseTime += event.Metadata.Duration
				responseCount++
			}
			stats.TotalTokens += event.Metadata.TokenCount
		case BehaviorCommand:
			stats.TotalCommands++
			if event.Metadata.CommandName != "" {
				stats.TopCommands[event.Metadata.CommandName]++
			}
		case BehaviorSession:
			stats.TotalSessions++
			if event.Metadata.Duration > 0 {
				totalSessionTime += event.Metadata.Duration
				sessionCount++
			}
		case BehaviorError:
			stats.TotalErrors++
		}

		// Track hourly activity
		hour := event.Timestamp.Hour()
		stats.HourlyActivity[hour]++

		// Track daily activity
		day := event.Timestamp.Format("2006-01-02")
		stats.DailyActivity[day]++
	}

	// Calculate averages
	if responseCount > 0 {
		stats.AvgResponseTime = totalResponseTime / time.Duration(responseCount)
	}
	if sessionCount > 0 {
		stats.AvgSessionTime = totalSessionTime / time.Duration(sessionCount)
	}

	// Calculate error rate
	totalEvents := stats.TotalQueries + stats.TotalResponses + stats.TotalCommands
	if totalEvents > 0 {
		stats.ErrorRate = float64(stats.TotalErrors) / float64(totalEvents)
	}

	// Find most active hour
	maxActivity := 0
	for hour, count := range stats.HourlyActivity {
		if count > maxActivity {
			maxActivity = count
			stats.MostActiveHour = hour
		}
	}

	// Calculate productivity score (simple heuristic)
	stats.ProductivityScore = a.calculateProductivityScore(stats)

	return stats, nil
}

// AnalyzeToday analyzes today's behaviors
func (a *Analyzer) AnalyzeToday() (*Statistics, error) {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	return a.AnalyzePeriod(startOfDay, endOfDay)
}

// AnalyzeWeek analyzes the past week's behaviors
func (a *Analyzer) AnalyzeWeek() (*Statistics, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)
	return a.AnalyzePeriod(startDate, endDate)
}

// AnalyzeMonth analyzes the past month's behaviors
func (a *Analyzer) AnalyzeMonth() (*Statistics, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)
	return a.AnalyzePeriod(startDate, endDate)
}

// loadEventsByPeriod loads all events within a date range
func (a *Analyzer) loadEventsByPeriod(startDate, endDate time.Time) ([]BehaviorEvent, error) {
	var allEvents []BehaviorEvent

	// Iterate through each day in the period
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		dateStr := currentDate.Format("2006-01-02")
		filePath := filepath.Join(a.dataPath, fmt.Sprintf("events-%s.json", dateStr))

		// Read events from file
		data, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				// File doesn't exist for this day, skip
				currentDate = currentDate.AddDate(0, 0, 1)
				continue
			}
			return nil, fmt.Errorf("failed to read events file %s: %w", filePath, err)
		}

		var dayEvents []BehaviorEvent
		if err := json.Unmarshal(data, &dayEvents); err != nil {
			return nil, fmt.Errorf("failed to unmarshal events from %s: %w", filePath, err)
		}

		// Filter events by time range
		for _, event := range dayEvents {
			if (event.Timestamp.After(startDate) || event.Timestamp.Equal(startDate)) &&
				event.Timestamp.Before(endDate) {
				allEvents = append(allEvents, event)
			}
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return allEvents, nil
}

// calculateProductivityScore calculates a productivity score based on various metrics
func (a *Analyzer) calculateProductivityScore(stats *Statistics) float64 {
	score := 0.0

	// Higher score for more queries (engagement)
	if stats.TotalQueries > 0 {
		score += float64(stats.TotalQueries) * 0.5
	}

	// Lower score for errors
	score -= float64(stats.TotalErrors) * 2.0

	// Bonus for consistent activity
	if len(stats.DailyActivity) > 1 {
		score += float64(len(stats.DailyActivity)) * 2.0
	}

	// Normalize to 0-100 range
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// GetInsights generates insights from statistics
func GetInsights(stats *Statistics) []string {
	insights := make([]string, 0)

	// Activity insights
	if stats.TotalQueries > 0 {
		insights = append(insights, fmt.Sprintf("You made %d queries during this period", stats.TotalQueries))
	}

	// Error insights
	if stats.ErrorRate > 0.1 {
		insights = append(insights, fmt.Sprintf("Warning: Error rate is %.1f%%. Consider reviewing common issues", stats.ErrorRate*100))
	}

	// Time management insights
	if stats.MostActiveHour >= 9 && stats.MostActiveHour <= 17 {
		insights = append(insights, fmt.Sprintf("Your most active hour is %d:00 (during work hours)", stats.MostActiveHour))
	} else {
		insights = append(insights, fmt.Sprintf("Your most active hour is %d:00 (outside typical work hours)", stats.MostActiveHour))
	}

	// Session insights
	if stats.AvgSessionTime > 30*time.Minute {
		insights = append(insights, fmt.Sprintf("Average session time: %v - Consider taking breaks", stats.AvgSessionTime.Round(time.Minute)))
	}

	// Productivity insights
	if stats.ProductivityScore > 70 {
		insights = append(insights, fmt.Sprintf("Great productivity! Score: %.1f/100", stats.ProductivityScore))
	} else if stats.ProductivityScore > 40 {
		insights = append(insights, fmt.Sprintf("Good productivity. Score: %.1f/100", stats.ProductivityScore))
	} else {
		insights = append(insights, fmt.Sprintf("Productivity could be improved. Score: %.1f/100", stats.ProductivityScore))
	}

	// Top commands insight
	if len(stats.TopCommands) > 0 {
		topCmd := ""
		maxCount := 0
		for cmd, count := range stats.TopCommands {
			if count > maxCount {
				maxCount = count
				topCmd = cmd
			}
		}
		insights = append(insights, fmt.Sprintf("Most used command: '%s' (%d times)", topCmd, maxCount))
	}

	return insights
}

// FormatStatistics formats statistics for display
func FormatStatistics(stats *Statistics) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== Statistics for %s ===\n\n", stats.Period))
	sb.WriteString(fmt.Sprintf("Total Queries:      %d\n", stats.TotalQueries))
	sb.WriteString(fmt.Sprintf("Total Responses:    %d\n", stats.TotalResponses))
	sb.WriteString(fmt.Sprintf("Total Commands:     %d\n", stats.TotalCommands))
	sb.WriteString(fmt.Sprintf("Total Sessions:     %d\n", stats.TotalSessions))
	sb.WriteString(fmt.Sprintf("Total Errors:       %d\n", stats.TotalErrors))
	sb.WriteString(fmt.Sprintf("Error Rate:         %.1f%%\n", stats.ErrorRate*100))

	if stats.AvgResponseTime > 0 {
		sb.WriteString(fmt.Sprintf("Avg Response Time:  %v\n", stats.AvgResponseTime.Round(time.Millisecond)))
	}
	if stats.AvgSessionTime > 0 {
		sb.WriteString(fmt.Sprintf("Avg Session Time:   %v\n", stats.AvgSessionTime.Round(time.Minute)))
	}

	sb.WriteString(fmt.Sprintf("Total Tokens:       %d\n", stats.TotalTokens))
	sb.WriteString(fmt.Sprintf("Most Active Hour:   %d:00\n", stats.MostActiveHour))
	sb.WriteString(fmt.Sprintf("Productivity Score: %.1f/100\n", stats.ProductivityScore))

	// Top commands
	if len(stats.TopCommands) > 0 {
		sb.WriteString("\nTop Commands:\n")

		// Sort commands by count
		type cmdCount struct {
			cmd   string
			count int
		}
		cmds := make([]cmdCount, 0, len(stats.TopCommands))
		for cmd, count := range stats.TopCommands {
			cmds = append(cmds, cmdCount{cmd, count})
		}
		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i].count > cmds[j].count
		})

		for i, cc := range cmds {
			if i >= 5 {
				break
			}
			sb.WriteString(fmt.Sprintf("  %d. %s: %d times\n", i+1, cc.cmd, cc.count))
		}
	}

	// Daily activity
	if len(stats.DailyActivity) > 0 {
		sb.WriteString("\nDaily Activity:\n")

		// Sort days
		days := make([]string, 0, len(stats.DailyActivity))
		for day := range stats.DailyActivity {
			days = append(days, day)
		}
		sort.Strings(days)

		for _, day := range days {
			count := stats.DailyActivity[day]
			sb.WriteString(fmt.Sprintf("  %s: %d events\n", day, count))
		}
	}

	// Insights
	insights := GetInsights(stats)
	if len(insights) > 0 {
		sb.WriteString("\n=== Insights ===\n\n")
		for i, insight := range insights {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, insight))
		}
	}

	return sb.String()
}
