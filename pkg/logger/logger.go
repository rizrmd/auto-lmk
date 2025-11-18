package logger

import (
	"log/slog"
	"os"
)

// Setup initializes the structured logger
func Setup(env string) {
	// Open log file
	logFile, err := os.OpenFile("/tmp/auto-lmk.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fallback to stdout if file creation fails
		logFile = os.Stdout
	}

	var handler slog.Handler

	if env == "production" {
		// JSON handler for production
		handler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Text handler for development
		handler = slog.NewTextHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
