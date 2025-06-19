// File: menu_service/internal/app/app.go
package app

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	grpcapp "github.com/netscrawler/Restaurant_is/menu_service/internal/app/grpc"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/app/health"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/app/metrics"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/config"
	minio "github.com/netscrawler/Restaurant_is/menu_service/internal/infra/out/mini_io"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/infra/out/postgres"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/repository"
	pgrepo "github.com/netscrawler/Restaurant_is/menu_service/internal/repository/pg_repo"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/service"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/telemetry"
)

type App struct {
	log        *slog.Logger
	grpcApp    *grpcapp.App
	db         *postgres.Storage
	minio      *minio.Storage
	healz      *health.App
	metricsApp *metrics.App
	telemetry  *telemetry.Telemetry
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// Инициализация хранилищ
	pgStorage := postgres.MustSetup(context.Background(), cfg.DB.GetURL(), log)
	minioStorage := minio.MustSetup(
		context.Background(),
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.UseSSL,
		[]string{cfg.MinIO.Bucket},
		log,
	)

	dishRepo := repository.NewDishRepository(pgrepo.NewDishPgRepo(pgStorage))
	imageService := service.NewImageService(
		minioStorage,
		cfg.MinIO.Bucket,
		cfg.MinIO.UrlExpiry,
	)
	dishService := service.NewDishService(dishRepo, imageService)

	// Инициализация телеметрии
	telemetryApp, err := telemetry.New(&cfg.Telemetry, log)
	if err != nil {
		panic(err)
	}

	// Настройка gRPC сервера
	grpcApp := grpcapp.New(
		log,
		dishService,
		imageService,
		cfg.GRPCServer.Port,
		telemetryApp,
	)

	healthApp := health.New([]func() error{
		func() error {
			return pgStorage.DB.Ping(context.Background())
		},
		func() error {
			cancel, err := minioStorage.HealthCheck(5 * time.Second)
			if err != nil {
				return nil
			}
			defer cancel()
			if !minioStorage.IsOnline() {
				return errors.New("minio not available")
			}

			return nil
		},
	}, cfg.GRPCServer.Address, strconv.Itoa(cfg.GRPCServer.Port+100))

	metricsApp := metrics.New(log, telemetryApp, cfg.Telemetry.MetricsPort)

	return &App{
		log:        log,
		grpcApp:    grpcApp,
		db:         pgStorage,
		minio:      minioStorage,
		healz:      healthApp,
		metricsApp: metricsApp,
		telemetry:  telemetryApp,
	}
}

func (a *App) MustRun() {
	// Запускаем метрики в отдельной горутине
	go func() {
		a.metricsApp.MustRun()
	}()

	// Запускаем gRPC сервер в основной горутине
	a.grpcApp.MustRun()
}

func (a *App) Stop() {
	a.log.Info("stopping application")

	// Останавливаем gRPC сервер
	a.grpcApp.Stop()

	a.db.Stop(context.Background())

	// Останавливаем метрики
	a.metricsApp.Stop()

	// Останавливаем телеметрию
	if a.telemetry != nil {
		if err := a.telemetry.Shutdown(context.Background()); err != nil {
			a.log.Error("failed to shutdown telemetry", slog.Any("error", err))
		}
	}

	a.log.Info("application stopped")
}
