package app

import (
	"context"
	"log/slog"
	"strconv"

	grpcapp "github.com/netscrawler/Restaurant_is/order_service/internal/app/grpc"
	"github.com/netscrawler/Restaurant_is/order_service/internal/app/health"
	"github.com/netscrawler/Restaurant_is/order_service/internal/config"
	"github.com/netscrawler/Restaurant_is/order_service/internal/infra/in/postgres"
	grpcinfra "github.com/netscrawler/Restaurant_is/order_service/internal/infra/out/grpc"
	"github.com/netscrawler/Restaurant_is/order_service/internal/infra/out/kafka"
	pg "github.com/netscrawler/Restaurant_is/order_service/internal/repository/pg_repo"
	"github.com/netscrawler/Restaurant_is/order_service/internal/service"
)

type App struct {
	app    *grpcapp.App
	healz  *health.App
	cancel context.CancelFunc
	kafka  *service.Event
}

func New(log *slog.Logger, cfg *config.Config) *App {
	ctx, cancel := context.WithCancel(context.Background())
	storage := postgres.MustSetup(ctx, cfg.DB.GetURL(), log)

	orderRepo := pg.NewPgOrder(storage)
	eventRepo := pg.NewPgEvent(storage)

	kafka, err := kafka.NewKafkaPublisher(&cfg.Kafka, log)
	if err != nil {
		panic(err)
	}

	menuClient, err := grpcinfra.New(ctx, cfg.MenuClient)
	if err != nil {
		panic(err)
	}

	eventSerice := service.NewEventService(eventRepo, kafka, cfg, log)

	orderService := service.NewOrder(orderRepo, eventRepo, menuClient)

	grpcapp := grpcapp.New(log, orderService, cfg.GRPCServer.Port)

	heal := health.New(
		[]func() error{
			func() error { return storage.DB.Ping(context.Background()) },
		},
		cfg.GRPCServer.Address,
		strconv.Itoa(cfg.GRPCServer.Port+1),
	)

	return &App{
		app:    grpcapp,
		healz:  heal,
		cancel: cancel,
		kafka:  eventSerice,
	}
}

func (a *App) MustRun() {
	a.healz.Start()
	a.app.MustRun()
	// panic("err")
}

func (a *App) Stop() {
	panic("err")
}
