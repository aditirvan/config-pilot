package monitor

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aditirvan/config-pilot/internal/config"
	"github.com/aditirvan/config-pilot/internal/github"
	"github.com/aditirvan/config-pilot/internal/utils"
)

// CommitHandler is a function type for handling new commits
type CommitHandler func(commit *github.Commit)

// Service continuously monitors for new commits
type Service struct {
	client       *github.Client
	lastKnownSHA string
	handler      CommitHandler
	interval     time.Duration
	monitorPath  string
	repoPath     string // Path to the local repository
	owner        string // Repository owner
	repo         string // Repository name
}

// NewService creates a new commit monitoring service
func NewService(client *github.Client, config *config.Config, handler CommitHandler) *Service {
	return &Service{
		client:      client,
		handler:     handler,
		interval:    time.Duration(config.Interval) * time.Second,
		monitorPath: config.MonitorPath,
		repoPath:    config.Repo,
		owner:       config.Owner,
		repo:        config.Repo,
	}
}

// Start begins the continuous monitoring
func (s *Service) Start() error {
	// Get initial commit to establish baseline
	initialCommit, err := s.client.GetLatestCommit(s.monitorPath)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("failed to get initial commit: %s", err.Error()))
		return err
	}

	s.lastKnownSHA = initialCommit.SHA
	utils.Logger.Info(fmt.Sprintf("Monitoring started. Initial commit: %s", initialCommit.SHA[:7]))
	if s.monitorPath != "" {
		utils.Logger.Info(fmt.Sprintf("Monitoring path: %s", s.monitorPath))
	}

	// Start monitoring loop
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		s.checkForUpdates()
	}

	return nil
}

// checkForUpdates checks for new commits and triggers handler if found
func (s *Service) checkForUpdates() {
	latestCommit, err := s.client.GetLatestCommit(s.monitorPath)
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("Error checking for updates: %v", err))
		return
	}

	if latestCommit.SHA != s.lastKnownSHA {
		s.lastKnownSHA = latestCommit.SHA
		s.handler(latestCommit)
	}
}

// DefaultCommitHandler returns a commit handler that uses configurable repository values
func DefaultCommitHandler(config *config.Config) CommitHandler {
	return func(commit *github.Commit) {
		utils.Logger.Info("new commit detected", slog.String("sha", commit.SHA[:7]), slog.String("author", commit.Commit.Author.Name), slog.String("commit_msg", commit.Commit.Message))
		// Pull the repository with configurable values
		if err := pullRepository(config); err != nil {
			utils.Logger.Error(err.Error())
			return
		}
	}
}

// pullRepository pulls the latest changes from the repository
func pullRepository(config *config.Config) error {

	if _, err := os.Stat("data/files"); os.IsNotExist(err) {
		// Jika folder tidak ada, buat folder
		err := os.MkdirAll("data/files", os.ModePerm)
		if err != nil {
			return err
		}
	}

	err := os.RemoveAll(fmt.Sprintf("data/%s", config.Repo))
	if err != nil {
		return err
	}

	err = os.RemoveAll("data/files")
	if err != nil {
		return err
	}

	err = os.RemoveAll("data/files/script.sh")
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	cmd := exec.Command("git", "clone", fmt.Sprintf("https://git:%s@github.com/%s/%s.git", config.GithubToken, config.Owner, config.Repo))
	cmd.Dir = "data" // Ensure we are in the correct directory
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %s, %w", output, err)
	}
	err = os.Rename(fmt.Sprintf("data/%s/%s", config.Repo, config.MonitorPath), "data/files")
	if err != nil {
		return err
	}
	err = os.RemoveAll(fmt.Sprintf("data/%s", config.Repo))
	if err != nil {
		return err
	}

	err = DecryptFiles(config)
	if err != nil {
		return err
	}

	err = ExecutionScript(config)
	if err != nil {
		return err
	}

	return nil
}

func ExecutionScript(config *config.Config) error {
	utils.Logger.Info("Starting exection script", slog.String("input", config.Script))
	defaultContent := `#!/bin/bash
		cd data/files`
	err := os.WriteFile("data/files/script.sh", []byte(fmt.Sprintf("%s\n\n%s", defaultContent, config.Script)), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("bash", "data/files/script.sh")
	// Menangkap output dan error
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// Menampilkan output dari skrip

	utils.Logger.Info("execution done", slog.String("output", string(output)))

	return nil
}

func DecryptFiles(config *config.Config) error {
	root := "data/files" // Ganti dengan path direktori yang ingin Anda telusuri
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() { // Hanya menampilkan file, bukan direktori

			cmd := exec.Command("sops", "-d", "-i", path)
			cmd.Env = append(os.Environ(), "SOPS_AGE_KEY="+config.AgeKey)

			_, err := cmd.CombinedOutput()
			if err == nil {
				utils.Logger.Info("File decrypted", slog.String("name", strings.ReplaceAll(path, "data/files/", "")))
			}

		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
