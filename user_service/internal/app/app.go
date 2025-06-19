package app

import (
	"context"
	"log/slog"
	"strconv"

	grpcapp "user_service/internal/app/grpc"
	"user_service/internal/app/health"
	metricsapp "user_service/internal/app/metrics"
	"user_service/internal/config"
	application "user_service/internal/domain/app"
	"user_service/internal/domain/service"
	"user_service/internal/infra/in/kafka"
	pgrepo "user_service/internal/infra/out/postgres"
	"user_service/internal/storage/postgres"
	"user_service/internal/telemetry"
)

type App struct {
	cfg           *config.Config
	log           *slog.Logger
	db            *postgres.Storage
	grpcApp       *grpcapp.App
	health        *health.App
	metricsApp    *metricsapp.App
	telemetry     *telemetry.Telemetry
	kafkaConsumer *kafka.UserEventConsumer
}

func New(log *slog.Logger, cfg *config.Config) *App {
	ctx := context.Background()

	db := postgres.MustSetup(ctx, cfg.DB.GetURL(), log)

	// Инициализация telemetry
	telemetryInstance, err := telemetry.New(&cfg.Telemetry, log)
	if err != nil {
		panic("failed to initialize telemetry: " + err.Error())
	}

	// Инициализация репозиториев
	userRepo := pgrepo.NewUserRepository(db)
	staffRepo := pgrepo.NewStaffRepository(db)
	roleRepo := pgrepo.NewRoleRepository(db)
	userRoleRepo := pgrepo.NewUserRoleRepository(db.DB)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	staffService := service.NewStaffService(staffRepo)
	roleService := service.NewRoleService(roleRepo, userRoleRepo)

	// Инициализация application сервисов
	userAppService := application.NewUserAppService(userService)
	staffAppService := application.NewStaffAppService(staffService)
	roleAppService := application.NewRoleAppService(roleService)

	// Инициализация Kafka consumer
	kafkaConsumer, err := kafka.NewUserEventConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.GroupID,
		cfg.Kafka.Topic,
		userAppService,
	)
	if err != nil {
		panic("failed to create kafka consumer: " + err.Error())
	}

	grpc := grpcapp.New(
		log,
		userAppService,
		staffAppService,
		roleAppService,
		cfg.GRPCServer.Port,
		telemetryInstance,
	)

	health := health.New(
		[]func() error{
			func() error {
				return db.DB.Ping(context.Background())
			},
		},
		cfg.GRPCServer.Address,
		strconv.Itoa(cfg.GRPCServer.Port+100),
	)
	metricsApp := metricsapp.New(log, telemetryInstance, cfg.Telemetry.MetricsPort)

	return &App{
		cfg:           cfg,
		log:           log,
		db:            db,
		grpcApp:       grpc,
		health:        health,
		metricsApp:    metricsApp,
		telemetry:     telemetryInstance,
		kafkaConsumer: kafkaConsumer,
	}
}

func (a *App) MustRun() {
	// Инициализация gRPC сервера
	go func() {
		a.metricsApp.MustRun()
	}()
	go func() {
		if err := a.kafkaConsumer.Start(context.Background()); err != nil {
			a.log.Error("Kafka consumer stopped", "err", err)
		}
	}()
	a.grpcApp.MustRun()
}

func (a *App) Stop() {
	a.log.Info("stopping application")

	a.grpcApp.Stop()

	// Закрытие соединения с базой данных
	if a.db != nil {
		a.db.Stop(context.Background())
		a.log.Info("database connection closed")
	}
	a.metricsApp.Stop()
}
