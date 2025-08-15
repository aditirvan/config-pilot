package main

import (
	"fmt"
	"os"

	"github.com/aditirvan/config-pilot/internal/config"
	"github.com/aditirvan/config-pilot/internal/github"
	"github.com/aditirvan/config-pilot/internal/monitor"
	"github.com/aditirvan/config-pilot/internal/utils"
)

func main() {

	// monitor.DecryptFiles(nil)

	// return

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		utils.Logger.Error("Please set GITHUB_TOKEN environment variable")
		utils.Logger.Error("Usage: export GITHUB_TOKEN=your_personal_access_token")
		return
	}

	// Load configuration from config.yaml
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error loading config: %s", err.Error()))
		os.Exit(1)
	}

	// Validate configuration
	if cfg.Owner == "" || cfg.Repo == "" {
		utils.Logger.Error("Error: owner and repo must be specified in config.yaml")
		os.Exit(1)
	}

	client := github.NewClient(token, cfg.Owner, cfg.Repo)
	commitHandler := monitor.DefaultCommitHandler(cfg)
	monitorService := monitor.NewService(client, cfg.MonitorPath, cfg.Repo, cfg.Owner, cfg.Repo, commitHandler)

	utils.Logger.Info("Starting gitops automation")

	if cfg.MonitorPath != "" {
		utils.Logger.Info(fmt.Sprintf("Monitoring path: %s (including subdirectories)", cfg.MonitorPath))
	} else {
		utils.Logger.Info("Monitoring entire repository")
	}
	utils.Logger.Info("Press Ctrl+C to stop monitoringy")

	if err := monitorService.Start(); err != nil {
		utils.Logger.Info(fmt.Sprintf("Error starting monitor: %v", err))
		os.Exit(1)
	}
}
