package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/netscrawler/Restaurant_is/auth/internal/app"
	"github.com/netscrawler/Restaurant_is/auth/internal/config"
	"github.com/netscrawler/Restaurant_is/auth/internal/telemetry"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Debug("start with config", slog.Any("config", cfg))

	// --- TELEMETRY INIT ---
	telemetryInstance, err := telemetry.New(&cfg.Telemetry, log)
	if err != nil {
		log.Error("failed to init telemetry", slog.Any("error", err))
		os.Exit(1)
	}
	defer telemetryInstance.Shutdown(context.Background())
	// --- END TELEMETRY INIT ---

	application := app.New(log, *cfg)
	go func() {
		if err := application.Run(); err != nil {
			log.Error("failed to run application", slog.Any("error", err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case dev, prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
