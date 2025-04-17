package main

import (
	"os"
	"os/signal"
	"syscall"

	"notify/internal/app"
	"notify/internal/config"

	"go.uber.org/zap"
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
	log, err := setupLogger(cfg.Env)
	if err != nil {
		panic(err)
	}
	log.Debug("start with config", zap.Any("config", cfg))

	application := app.New(log, cfg)
	go func() {
		application.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
}

func setupLogger(env string) (*zap.Logger, error) {
	var log *zap.Logger

	var err error

	switch env {
	case local:
		log, err = zap.NewDevelopment()
	case dev:
		log, err = zap.NewProduction()
	case prod:
		log, err = zap.NewProduction()
	default:
		log, err = zap.NewProduction()
	}

	return log, err
}
