package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/netscrawler/Restaurant_is/order_service/internal/app"
	"github.com/netscrawler/Restaurant_is/order_service/internal/config"
	"github.com/netscrawler/Restaurant_is/order_service/internal/telemetry"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "production"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log = log.With(slog.String("service", "order"))

	log.Debug("started with config", slog.Any("config", cfg))

	slog.SetDefault(log)

	telemetryInstance, err := telemetry.New(&cfg.Telemetry, log)
	if err != nil {
		log.Error("failed to init telemetry", slog.Any("err", err))
		os.Exit(1)
	}

	shutdownCtx, shutdownCancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer shutdownCancel()
	defer func() {
		if err := telemetryInstance.Shutdown(shutdownCtx); err != nil {
			log.Error("failed to shutdown telemetry", slog.Any("err", err))
		}
	}()

	application := app.New(log, cfg, telemetryInstance)

	go func() {
		application.MustRun()
	}()

	<-shutdownCtx.Done()

	application.Stop()
	log.Info("Gracefully stopped")
}

//nolint:exhaustruct
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
