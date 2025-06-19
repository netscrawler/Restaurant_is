package app

import (
	"context"
	"log/slog"
	"strconv"

	grpcapp "github.com/netscrawler/Restaurant_is/order_service/internal/app/grpc"
	"github.com/netscrawler/Restaurant_is/order_service/internal/app/health"
	metricsapp "github.com/netscrawler/Restaurant_is/order_service/internal/app/metrics"
	"github.com/netscrawler/Restaurant_is/order_service/internal/config"
	"github.com/netscrawler/Restaurant_is/order_service/internal/infra/in/postgres"
	grpcinfra "github.com/netscrawler/Restaurant_is/order_service/internal/infra/out/grpc"
	"github.com/netscrawler/Restaurant_is/order_service/internal/infra/out/kafka"
	pg "github.com/netscrawler/Restaurant_is/order_service/internal/repository/pg_repo"
	"github.com/netscrawler/Restaurant_is/order_service/internal/service"
	"github.com/netscrawler/Restaurant_is/order_service/internal/telemetry"
)

type App struct {
	app       *grpcapp.App
	healz     *health.App
	cancel    context.CancelFunc
	kafka     *service.Event
	telemetry *telemetry.Telemetry
	metrics   *metricsapp.App
	storage   *postgres.Storage
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

	telemetryApp, err := telemetry.New(&cfg.Telemetry, log)
	if err != nil {
		panic(err)
	}

	metricsApp := metricsapp.New(log, telemetryApp, cfg.Telemetry.MetricsPort)

	return &App{
		app:       grpcapp,
		healz:     heal,
		cancel:    cancel,
		kafka:     eventSerice,
		telemetry: telemetryApp,
		metrics:   metricsApp,
		storage:   storage,
	}
}

func (a *App) MustRun() {
	go func() {
		a.metrics.MustRun()
	}()
	a.healz.Start()
	a.app.MustRun()
}

func (a *App) Stop() {
	a.app.Stop()
	a.storage.Stop(context.Background())
	a.metrics.Stop()
}
