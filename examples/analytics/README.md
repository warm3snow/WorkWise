# Behavior Analytics Example

This example demonstrates the behavior analytics feature of WorkWise.

## Overview

WorkWise automatically tracks your interactions and provides insights into your usage patterns, helping you manage your time more effectively.

## Quick Start

1. **Enable Analytics** (enabled by default)

Analytics is automatically enabled in the default configuration. To verify or customize:

```yaml
extensions:
  analytics_enabled: true
  analytics_path: ~/.workwise/analytics
```

2. **Use WorkWise normally**

```bash
# Ask questions
workwise ask "What is the weather today?"
workwise ask "How do I create a Go project?"

# Start interactive session
workwise chat
```

3. **View Your Statistics**

```bash
# Today's stats
workwise stats today

# This week's stats
workwise stats week

# This month's stats
workwise stats month
```

## Sample Output

After using WorkWise, you might see statistics like:

```
=== Statistics for 2024-01-01 to 2024-01-02 ===

Total Queries:      15
Total Responses:    15
Total Commands:     3
Total Sessions:     2
Total Errors:       0
Error Rate:         0.0%
Avg Response Time:  1.2s
Avg Session Time:   25m
Total Tokens:       3500
Most Active Hour:   14:00
Productivity Score: 85.0/100

Top Commands:
  1. help: 2 times
  2. clear: 1 times

Daily Activity:
  2024-01-01: 10 events
  2024-01-02: 8 events

=== Insights ===

1. You made 15 queries during this period
2. Your most active hour is 14:00 (during work hours)
3. Great productivity! Score: 85.0/100
4. Most used command: 'help' (2 times)
```

## Generate Reports

Generate comprehensive reports for different time periods:

```bash
# Today's report
workwise report

# Weekly report
workwise report --period=week

# Monthly report
workwise report --period=month
```

## What Gets Tracked

### User Queries
Every question or request you make is logged with:
- Query content
- Timestamp
- Session ID

### AI Responses
Response information includes:
- Response content
- Duration (how long it took)
- Token count
- Model used
- Session ID

### Commands
Command usage tracking:
- Command name (help, clear, etc.)
- Timestamp
- Session ID

### Sessions
Session information:
- Session start/end
- Total duration
- Session ID

### Errors
Error tracking for debugging:
- Error message
- Timestamp
- Session ID

## Data Storage

All analytics data is stored locally in JSON format:

```
~/.workwise/analytics/
├── events-2024-01-01.json
├── events-2024-01-02.json
└── events-2024-01-03.json
```

Each file contains the events for that specific day in JSON format.

## Privacy

- **100% Local**: All data is stored on your machine
- **No Cloud**: Nothing is sent to external servers
- **Your Data**: You have full control over your analytics data
- **Easy to Disable**: Just set `analytics_enabled: false` in config

## Use Cases

### Time Management
- Track when you're most productive
- Identify time-consuming tasks
- Plan your day based on activity patterns

### Self-Improvement
- Monitor your learning progress
- Track which topics you explore most
- Measure productivity over time

### System Optimization
- Identify frequently used commands
- Find and fix recurring errors
- Optimize your workflow

## Tips

1. **Review Weekly**: Check your weekly stats every Monday to plan your week
2. **Monthly Reviews**: Generate monthly reports to track long-term progress
3. **Error Monitoring**: Watch your error rate to identify issues early
4. **Productivity Tracking**: Use the productivity score to measure improvement

## Disable Analytics

If you don't want analytics:

1. Edit your config file: `~/.workwise/config.yaml`
2. Set `analytics_enabled: false`
3. Restart WorkWise

Or create a new config with analytics disabled:

```bash
workwise config init
# Then edit the config file to disable analytics
```
