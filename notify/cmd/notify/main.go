package main

import (
	"log/slog"
	"notify/internal/app"
	"notify/internal/config"
	"os"
	"os/signal"
	"syscall"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := setupLogger(cfg.Env)
	log.Debug("start with config", "config", cfg)
	application := app.New(log, cfg)
	go func() {
		application.MustRun()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.Stop()
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler
	switch env {
	case local:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case dev, prod:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	return slog.New(handler)
}
