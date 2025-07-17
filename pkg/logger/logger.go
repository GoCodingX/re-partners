package logger

import (
	"log/slog"
	"os"
)

func Init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// set as the default logger globally
	slog.SetDefault(logger)

	logger.Info("logger initialized")
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Fatal(message string, err error) {
	slog.Error(message, slog.String("error", err.Error()))
	os.Exit(1)
}
