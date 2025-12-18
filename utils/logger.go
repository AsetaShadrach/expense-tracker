package utils

import (
	"log/slog"
	"os"
)

var GeneralLogger *slog.Logger

func InitiateLogger() {
	var LogLevel slog.Level
	if os.Getenv("ENV") == "production" {
		LogLevel = slog.LevelInfo
	} else {
		LogLevel = slog.Level(-8) //LevelTrace
	}

	GeneralLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: LogLevel,
	}))

	GeneralLogger.Info("General Logger has been loaded --")
}
