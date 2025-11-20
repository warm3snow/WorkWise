package main

import (
	"fmt"
	"os"

	"github.com/warm3snow/WorkWise/internal/cli"
	"github.com/warm3snow/WorkWise/internal/config"
)

var (
	// Version information, can be overridden during build
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Create and run CLI application
	app := cli.NewApp(cfg, Version, BuildTime, GitCommit)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
