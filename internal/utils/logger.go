package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

var Logger *slog.Logger

// LoggingConfig holds configuration for logging
type LoggingConfig struct {
	LogFilePath string `yaml:"logFilePath"`
	LogLevel    string `yaml:"logLevel"`
	LogToFile   bool   `yaml:"logToFile"`
}

func init() {
	// Initialize logger with default configuration
	InitializeLogger(LoggingConfig{
		LogLevel:  "info",
		LogToFile: false,
	})
}

// InitializeLogger initializes the logger with the given configuration
func InitializeLogger(config LoggingConfig) {
	// Parse log level
	var level slog.Level
	switch config.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create options for handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Create handlers
	var handlers []slog.Handler

	// Always add stdout handler
	stdoutHandler := slog.NewTextHandler(os.Stdout, opts)
	handlers = append(handlers, stdoutHandler)

	// Add file handler if enabled and path is provided
	if config.LogToFile && config.LogFilePath != "" {
		// Ensure directory exists
		dir := filepath.Dir(config.LogFilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
		} else {
			file, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
			} else {
				fileHandler := slog.NewTextHandler(file, opts)
				handlers = append(handlers, fileHandler)
			}
		}
	}

	// Create multi-handler if we have multiple handlers
	if len(handlers) == 1 {
		Logger = slog.New(handlers[0])
	} else {
		Logger = slog.New(&multiHandler{handlers: handlers})
	}
}

// multiHandler implements a handler that writes to multiple handlers
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	var lastErr error
	for _, h := range m.handlers {
		if err := h.Handle(ctx, record); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}

// LoadLoggingConfig loads logging configuration from environment variables and config
func LoadLoggingConfig(cfg *LoggingConfig) LoggingConfig {
	// Start with provided config
	result := *cfg

	// Override with environment variables if provided
	if logFilePath := os.Getenv("LOG_FILE_PATH"); logFilePath != "" {
		result.LogFilePath = logFilePath
		result.LogToFile = true
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		result.LogLevel = logLevel
	}

	if logToFile := os.Getenv("LOG_TO_FILE"); logToFile != "" {
		result.LogToFile = logToFile == "true" || logToFile == "1"
	}

	return result
}
