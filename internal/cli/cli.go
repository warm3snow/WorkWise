package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/warm3snow/WorkWise/internal/agent"
	"github.com/warm3snow/WorkWise/internal/analytics"
	"github.com/warm3snow/WorkWise/internal/config"
)

// App represents the CLI application
type App struct {
	config    *config.Config
	version   string
	buildTime string
	gitCommit string
	agent     *agent.Agent
	tracker   *analytics.Tracker
	analyzer  *analytics.Analyzer
}

// NewApp creates a new CLI application
func NewApp(cfg *config.Config, version, buildTime, gitCommit string) *App {
	// Initialize tracker
	tracker, err := analytics.NewTracker(cfg.Extensions.AnalyticsEnabled, cfg.Extensions.AnalyticsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to initialize analytics tracker: %v\n", err)
		tracker, _ = analytics.NewTracker(false, "") // Disabled tracker as fallback
	}

	// Initialize analyzer
	analyzer := analytics.NewAnalyzer(cfg.Extensions.AnalyticsPath)

	return &App{
		config:    cfg,
		version:   version,
		buildTime: buildTime,
		gitCommit: gitCommit,
		tracker:   tracker,
		analyzer:  analyzer,
	}
}

// Run starts the CLI application
func (a *App) Run(args []string) error {
	app := &cli.App{
		Name:    "workwise",
		Usage:   "Intelligent Desktop Assistant - Your AI-powered work companion",
		Version: fmt.Sprintf("%s (built: %s, commit: %s)", a.version, a.buildTime, a.gitCommit),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "Start in interactive mode",
				Value:   true,
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to configuration file",
				EnvVars: []string{"WORKWISE_CONFIG"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "chat",
				Aliases: []string{"c"},
				Usage:   "Start an interactive chat session",
				Action:  a.chatCommand,
			},
			{
				Name:    "ask",
				Aliases: []string{"a"},
				Usage:   "Ask a single question",
				Action:  a.askCommand,
			},
			{
				Name:  "config",
				Usage: "Manage configuration",
				Subcommands: []*cli.Command{
					{
						Name:   "show",
						Usage:  "Show current configuration",
						Action: a.configShowCommand,
					},
					{
						Name:   "init",
						Usage:  "Initialize configuration file",
						Action: a.configInitCommand,
					},
				},
			},
			{
				Name:   "version",
				Usage:  "Show version information",
				Action: a.versionCommand,
			},
			{
				Name:    "stats",
				Aliases: []string{"s"},
				Usage:   "Show usage statistics",
				Subcommands: []*cli.Command{
					{
						Name:   "today",
						Usage:  "Show today's statistics",
						Action: a.statsToday,
					},
					{
						Name:   "week",
						Usage:  "Show this week's statistics",
						Action: a.statsWeek,
					},
					{
						Name:   "month",
						Usage:  "Show this month's statistics",
						Action: a.statsMonth,
					},
				},
				Action: a.statsToday, // Default to today
			},
			{
				Name:    "report",
				Aliases: []string{"r"},
				Usage:   "Generate behavior analysis report",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "period",
						Aliases: []string{"p"},
						Usage:   "Report period (today, week, month)",
						Value:   "today",
					},
				},
				Action: a.reportCommand,
			},
		},
		Action: a.defaultAction,
	}

	return app.Run(args)
}

// defaultAction is the default action when no command is specified
func (a *App) defaultAction(c *cli.Context) error {
	if c.Bool("interactive") || c.NArg() == 0 {
		return a.chatCommand(c)
	}
	// If arguments are provided, treat them as a question
	return a.askCommand(c)
}

// chatCommand starts an interactive chat session
func (a *App) chatCommand(c *cli.Context) error {
	fmt.Println("WorkWise - Intelligent Desktop Assistant")
	fmt.Printf("Version: %s\n", a.version)
	fmt.Println("Type 'help' for commands, 'exit' or 'quit' to leave")
	fmt.Println("---")

	// Initialize agent
	if err := a.initAgent(); err != nil {
		return fmt.Errorf("failed to initialize agent: %w", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	for {
		fmt.Print(a.config.CLI.Prompt)
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle special commands
		switch strings.ToLower(input) {
		case "exit", "quit", "q":
			a.tracker.TrackSessionEnd()
			fmt.Println("Goodbye!")
			return nil
		case "help", "?":
			a.tracker.TrackCommand("help", a.tracker.GetSessionID())
			a.printHelp()
			continue
		case "clear", "cls":
			a.tracker.TrackCommand("clear", a.tracker.GetSessionID())
			fmt.Print("\033[H\033[2J")
			continue
		}

		// Track query
		a.tracker.TrackQuery(input, a.tracker.GetSessionID())

		// Process user input with agent
		startTime := time.Now()
		response, err := a.agent.Process(ctx, input)
		duration := time.Since(startTime)
		
		if err != nil {
			a.tracker.TrackError(err.Error(), a.tracker.GetSessionID())
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		// Track response
		a.tracker.TrackResponse(response, duration, 0, a.config.AI.Model, a.tracker.GetSessionID())

		fmt.Println(response)
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

// askCommand processes a single question
func (a *App) askCommand(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("please provide a question")
	}

	question := strings.Join(c.Args().Slice(), " ")

	// Initialize agent
	if err := a.initAgent(); err != nil {
		return fmt.Errorf("failed to initialize agent: %w", err)
	}

	// Track query
	sessionID := a.tracker.GetSessionID()
	a.tracker.TrackQuery(question, sessionID)

	ctx := context.Background()
	startTime := time.Now()
	response, err := a.agent.Process(ctx, question)
	duration := time.Since(startTime)
	
	if err != nil {
		a.tracker.TrackError(err.Error(), sessionID)
		return fmt.Errorf("error processing question: %w", err)
	}

	// Track response
	a.tracker.TrackResponse(response, duration, 0, a.config.AI.Model, sessionID)
	a.tracker.TrackSessionEnd()

	fmt.Println(response)
	return nil
}

// configShowCommand shows the current configuration
func (a *App) configShowCommand(c *cli.Context) error {
	fmt.Println("Current Configuration:")
	fmt.Println("---")
	fmt.Printf("Provider: %s\n", a.config.AI.Provider)
	fmt.Printf("Model: %s\n", a.config.AI.Model)
	if a.config.AI.BaseURL != "" {
		fmt.Printf("Base URL: %s\n", a.config.AI.BaseURL)
	}
	fmt.Printf("Max Iterations: %d\n", a.config.AI.Agent.MaxIterations)
	fmt.Printf("Temperature: %.2f\n", a.config.AI.Agent.Temperature)
	fmt.Printf("History Enabled: %v\n", a.config.AI.Agent.HistoryEnabled)
	fmt.Printf("MCP Enabled: %v\n", a.config.Extensions.MCPEnabled)
	fmt.Printf("Skills Enabled: %v\n", a.config.Extensions.SkillsEnabled)
	fmt.Printf("Desktop Enabled: %v\n", a.config.Extensions.DesktopEnabled)
	fmt.Printf("Analytics Enabled: %v\n", a.config.Extensions.AnalyticsEnabled)
	return nil
}

// configInitCommand initializes a configuration file
func (a *App) configInitCommand(c *cli.Context) error {
	if err := a.config.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	fmt.Println("Configuration file created successfully")
	return nil
}

// versionCommand shows version information
func (a *App) versionCommand(c *cli.Context) error {
	fmt.Printf("WorkWise version %s\n", a.version)
	fmt.Printf("Built: %s\n", a.buildTime)
	fmt.Printf("Commit: %s\n", a.gitCommit)
	return nil
}

// initAgent initializes the AI agent
func (a *App) initAgent() error {
	if a.agent != nil {
		return nil
	}

	var err error
	a.agent, err = agent.NewAgent(a.config)
	return err
}

// printHelp prints available commands
func (a *App) printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help, ?     - Show this help message")
	fmt.Println("  clear, cls  - Clear the screen")
	fmt.Println("  exit, quit  - Exit the application")
	fmt.Println("\nJust type your question or request and press Enter!")
}

// statsToday shows today's statistics
func (a *App) statsToday(c *cli.Context) error {
	if !a.config.Extensions.AnalyticsEnabled {
		fmt.Println("Analytics is not enabled. Enable it in your config file.")
		return nil
	}

	stats, err := a.analyzer.AnalyzeToday()
	if err != nil {
		return fmt.Errorf("failed to analyze statistics: %w", err)
	}

	fmt.Println(analytics.FormatStatistics(stats))
	return nil
}

// statsWeek shows this week's statistics
func (a *App) statsWeek(c *cli.Context) error {
	if !a.config.Extensions.AnalyticsEnabled {
		fmt.Println("Analytics is not enabled. Enable it in your config file.")
		return nil
	}

	stats, err := a.analyzer.AnalyzeWeek()
	if err != nil {
		return fmt.Errorf("failed to analyze statistics: %w", err)
	}

	fmt.Println(analytics.FormatStatistics(stats))
	return nil
}

// statsMonth shows this month's statistics
func (a *App) statsMonth(c *cli.Context) error {
	if !a.config.Extensions.AnalyticsEnabled {
		fmt.Println("Analytics is not enabled. Enable it in your config file.")
		return nil
	}

	stats, err := a.analyzer.AnalyzeMonth()
	if err != nil {
		return fmt.Errorf("failed to analyze statistics: %w", err)
	}

	fmt.Println(analytics.FormatStatistics(stats))
	return nil
}

// reportCommand generates a behavior analysis report
func (a *App) reportCommand(c *cli.Context) error {
	if !a.config.Extensions.AnalyticsEnabled {
		fmt.Println("Analytics is not enabled. Enable it in your config file.")
		return nil
	}

	period := c.String("period")
	var stats *analytics.Statistics
	var err error

	switch period {
	case "today":
		stats, err = a.analyzer.AnalyzeToday()
	case "week":
		stats, err = a.analyzer.AnalyzeWeek()
	case "month":
		stats, err = a.analyzer.AnalyzeMonth()
	default:
		return fmt.Errorf("invalid period: %s (must be today, week, or month)", period)
	}

	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	fmt.Println(analytics.FormatStatistics(stats))
	return nil
}
