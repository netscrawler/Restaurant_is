package app

import (
	"context"
	"log/slog"
	"strconv"

	grpcapp "user_service/internal/app/grpc"
	"user_service/internal/app/health"
	"user_service/internal/config"
	application "user_service/internal/domain/app"
	"user_service/internal/domain/service"
	pgrepo "user_service/internal/infra/out/postgres"
	"user_service/internal/storage/postgres"
)

type App struct {
	cfg     *config.Config
	log     *slog.Logger
	db      *postgres.Storage
	grpcApp *grpcapp.App
	health  *health.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	ctx := context.Background()

	db := postgres.MustSetup(ctx, cfg.DB.GetURL(), log)

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

	grpc := grpcapp.New(log, userAppService, staffAppService, roleAppService, cfg.GRPCServer.Port)

	health := health.New(
		[]func() error{
			func() error {
				return db.DB.Ping(context.Background())
			},
		},
		cfg.GRPCServer.Address,
		strconv.Itoa(cfg.GRPCServer.Port+1),
	)

	return &App{
		cfg:     cfg,
		log:     log,
		db:      db,
		grpcApp: grpc,
		health:  health,
	}
}

func (a *App) MustRun() {
	// Инициализация gRPC сервера
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
}
