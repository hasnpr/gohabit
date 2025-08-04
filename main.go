package main

import (
	"log/slog"
	"os"

	"github.com/hasnpr/gohabit/cmd"
)

func setupJSONLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func init() {
	setupJSONLogger()
}

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("can't execute app.", slog.Any("error", err))
	}
}
