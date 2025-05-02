// File: menu_service/internal/app/app.go
package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/netscrawler/Restaurant_is/menu_service/internal/app/grpc"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/config"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/repository"
	pgrepo "github.com/netscrawler/Restaurant_is/menu_service/internal/repository/pg_repo"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/service"
	minio "github.com/netscrawler/Restaurant_is/menu_service/internal/storage/mini_io"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/storage/postgres"
)

type App struct {
	grpcServer *grpcapp.App
	db         *postgres.Storage
	minio      *minio.Storage
}

func New(log *slog.Logger, cfg *config.Config) *App {
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
	dishService := service.NewDishService(dishRepo)
	imageService := service.NewImageService(minioStorage.Client, cfg.MinIO.Bucket, cfg.UrlExpiry)

	// Настройка gRPC сервера
	grpcApp := grpcapp.New(
		log,
		dishService,
		imageService,
		cfg.GRPCServer.Port,
	)

	return &App{
		grpcServer: grpcApp,
		db:         pgStorage,
		minio:      minioStorage,
	}
}

func (a *App) MustRun() {
	go a.grpcServer.MustRun()
}

func (a *App) Stop() {
	a.grpcServer.Stop()
	a.db.Stop(context.Background())
}
