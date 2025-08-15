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

	// Load configuration from config file specified by CONFIG_PATH environment variable
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("Error loading config from %s: %s", configPath, err.Error()))
		os.Exit(1)
	}

	// Validate configuration
	if cfg.Owner == "" || cfg.Repo == "" {
		utils.Logger.Error("Error: owner and repo must be specified in config.yaml")
		os.Exit(1)
	}

	if cfg.GithubToken == "" {
		utils.Logger.Error("Error: githubToken must be specified in config.yaml")
		utils.Logger.Error("Usage: add githubToken: your_personal_access_token to config.yaml")
		return
	}

	client := github.NewClient(cfg.GithubToken, cfg.Owner, cfg.Repo)
	commitHandler := monitor.DefaultCommitHandler(cfg)
	monitorService := monitor.NewService(client, cfg, commitHandler)

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
