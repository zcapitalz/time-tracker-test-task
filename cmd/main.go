package main

import (
	"fmt"
	"log/slog"
	"os"
	"time-tracker/internal/app"
	"time-tracker/internal/config"
)

func main() {
	cfg := config.MustNew()
	logger := mustNewLogger(cfg.Env)
	app.Run(cfg, logger)
}

func mustNewLogger(env config.Env) (logger *slog.Logger) {
	switch env {
	case config.EnvLocal:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvDev, config.EnvProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic(fmt.Errorf("unknown env: %v", env))
	}

	return logger
}
