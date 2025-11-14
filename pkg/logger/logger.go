package logger

import (
	"log/slog"
	"os"
)

// Setup initializes the structured logger
func Setup(env string) {
	var handler slog.Handler

	if env == "production" {
		// JSON handler for production
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Text handler for development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
