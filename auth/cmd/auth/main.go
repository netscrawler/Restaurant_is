package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/netscrawler/Restaurant_is/auth/internal/config"
	"go.uber.org/zap"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log, err := setupLogger(cfg.Env)
	if err != nil {
		panic(err)
	}

	log.Debug("start with config", zap.Any("config", cfg))
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	// application.GRPCServer.Stop()
}

func setupLogger(env string) (*zap.Logger, error) {
	var log *zap.Logger

	var err error
	// TODO: Разграничить уровни логгирования
	switch env {
	case local:
		log, err = zap.NewDevelopment()
		// log, err = zap.NewProduction()
	case dev:
		log, err = zap.NewProduction()
	case prod:
		log, err = zap.NewProduction()
	}

	return log, err
}
