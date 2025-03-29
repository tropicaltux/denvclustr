// Package logger initializes a default logger for project-x.
// It automatically initializes a default logger during package import
// that writes logs to a file in the appropriate data directory for the user's OS.
package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	defaultLogFileName = "project-x.log"
	envLogDir          = "PROJECT_X_LOG_DIR"
	appName            = "project-x"
)

// dummyLogger creates a logger that discards all log messages.
// This is used as a fallback when the file logger cannot be initialized.
func dummyLogger() *slog.Logger {
	discardWriter := slog.NewTextHandler(os.NewFile(0, os.DevNull), nil)
	return slog.New(discardWriter)
}

// getDefaultLogDir determines the appropriate directory for storing log files.
// It checks for a custom directory in the PROJECT_X_LOG_DIR environment variable first.
// Otherwise, it uses platform-specific directories according to XDG specifications:
//   - Linux: ~/.local/state/project-x/logs
//   - macOS: ~/Library/Application Support/project-x/logs
//   - Windows: %APPDATA%\project-x\logs
func getDefaultLogDir() (string, error) {
	if custom := os.Getenv(envLogDir); custom != "" {
		return custom, nil
	}

	return filepath.Join(xdg.StateHome, appName, "logs"), nil
}

// tryInitLogger attempts to initialize a file-based logger.
// It creates the log directory if it doesn't exist, opens/creates the log file,
// and configures a text-based logger with INFO level threshold.
// Returns the configured logger or an error if initialization fails.
func tryInitLogger() (*slog.Logger, error) {
	logDir, err := getDefaultLogDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine log dir: %w", err)
	}

	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("cannot create log dir %q: %w", logDir, err)
	}

	logPath := filepath.Join(logDir, defaultLogFileName)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("cannot open log file: %w", err)
	}

	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler), nil
}

// init is automatically called when the package is imported.
// It sets up the default logger for the application.
// If logger initialization fails, it will print a warning to stderr
// and fall back to a dummy logger that discards all messages.
func init() {
	logger, err := tryInitLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n🚨 [WARNING] logger FAILED to initialize: %v\n"+
			"🔇 Falling back to dummy logger. Logs will NOT be saved to file.\n\n", err)

		slog.SetDefault(dummyLogger())
		return
	}

	slog.SetDefault(logger)
}
