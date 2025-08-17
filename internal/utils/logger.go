package utils

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	// Menginisialisasi logger di dalam fungsi init
	handler := slog.NewTextHandler(os.Stdout, nil)
	Logger = slog.New(handler)
}
