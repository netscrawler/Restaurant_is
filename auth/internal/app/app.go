package app

import (
	grpcapp "github.com/netscrawler/Restaurant_is/auth/internal/app/grpc"
	"github.com/netscrawler/Restaurant_is/auth/internal/config"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type App struct {
	log      *zap.Logger
	gRPCServ *grpcapp.App
	db       *postgres.Storage
	cfg      *config.Config
}

func New(log *zap.Logger, cfg config.Config) *App {
	const op = "app.New"

	// TODO: implement new app logic.
	return &App{
		log:      log,
		gRPCServ: &grpcapp.App{},
		db:       &postgres.Storage{},
		cfg:      &config.Config{},
	}
}

func (a *App) Run() error {
	err := a.gRPCServ.Run()
	if err != nil {
		a.db.Stop()

		return err
	}

	return nil
}

func (a *App) Stop() {
	a.db.Stop()
	a.gRPCServ.Stop()
}
