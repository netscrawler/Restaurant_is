package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/netscrawler/Restaurant_is/auth/internal/app/grpc"
	notifyclient "github.com/netscrawler/Restaurant_is/auth/internal/app/notifyclient"
	"github.com/netscrawler/Restaurant_is/auth/internal/config"
	notify "github.com/netscrawler/Restaurant_is/auth/internal/infra/out/grpc"
	"github.com/netscrawler/Restaurant_is/auth/internal/infra/out/postgres"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	inmemcache "github.com/netscrawler/Restaurant_is/auth/internal/repository/in_mem_cache"
	pgrepo "github.com/netscrawler/Restaurant_is/auth/internal/repository/pg_repo"
	"github.com/netscrawler/Restaurant_is/auth/internal/service"
	"github.com/netscrawler/Restaurant_is/auth/internal/utils"
)

type App struct {
	log          *slog.Logger
	gRPCServ     *grpcapp.App
	db           *postgres.Storage
	notyfyclient *notifyclient.Client
	cfg          *config.Config
}

func New(log *slog.Logger, cfg config.Config) *App {
	db := postgres.MustSetup(context.Background(), cfg.DB.GetURL(), log)
	clientRepo := repository.NewClient(pgrepo.NewPgClient(db, log))
	auditRepo := repository.NewAudit(pgrepo.NewPgAudit(db, log))
	// oauthRepo := repository.NewOAuth(pgrepo.NewPgOauth(db, log))
	stafRepo := repository.NewStaff(pgrepo.NewPgStaff(db, log))
	tokenRepo := repository.NewToken(pgrepo.NewPgToken(db, log))

	notifyClient, err := notifyclient.New(context.Background(), cfg.NotifyClient)
	if err != nil {
		panic(err)
	}

	notifySender := notify.New(log, notifyClient)

	jwt, _ := utils.NewJWTManager(cfg.JWT)

	codeProvider := inmemcache.New(cfg.CodeLife)

	authService := service.NewAuthService(
		log,
		clientRepo,
		stafRepo,
		// tokenRepo,
		// oauthRepo,
		notifySender,
		codeProvider,
		jwt,
	)
	audit := service.NewAuditService(auditRepo, log)
	token := service.NewTokenService(tokenRepo, jwt, log)
	user := service.NewUserService(clientRepo, stafRepo, notifySender, log)
	gRPCServ := grpcapp.New(log, authService, audit, token, user, cfg.GRPCServer.Port)

	return &App{
		log:          log,
		gRPCServ:     gRPCServ,
		db:           db,
		cfg:          &cfg,
		notyfyclient: notifyClient,
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
