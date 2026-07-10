// Package logging provides a simple logging utility for the RunLog API application. It uses the slog package to create a logger that outputs log messages in JSON format to standard output. The logger is configured to log messages at the Info level and above, making it suitable for capturing important events and errors in the application. This package is intended to be used throughout the application to ensure consistent logging practices.
package logging

import (
	"log/slog"
	"os"
)

// New creates a new instance of slog.Logger with a JSON handler that writes log messages to standard output. The logger is configured to log messages at the Info level and above. This function is used to initialize the logging system for the RunLog API application.
func New() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}
