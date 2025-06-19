package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/netscrawler/Restaurant_is/auth/internal/app/grpc"
	metricsapp "github.com/netscrawler/Restaurant_is/auth/internal/app/metrics"
	notifyclient "github.com/netscrawler/Restaurant_is/auth/internal/app/notifyclient"
	"github.com/netscrawler/Restaurant_is/auth/internal/config"
	notify "github.com/netscrawler/Restaurant_is/auth/internal/infra/out/grpc"
	kafkainfra "github.com/netscrawler/Restaurant_is/auth/internal/infra/out/kafka"
	"github.com/netscrawler/Restaurant_is/auth/internal/infra/out/postgres"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	inmemcache "github.com/netscrawler/Restaurant_is/auth/internal/repository/in_mem_cache"
	kafkarepo "github.com/netscrawler/Restaurant_is/auth/internal/repository/kafka_repo"
	pgrepo "github.com/netscrawler/Restaurant_is/auth/internal/repository/pg_repo"
	"github.com/netscrawler/Restaurant_is/auth/internal/service"
	"github.com/netscrawler/Restaurant_is/auth/internal/telemetry"
	"github.com/netscrawler/Restaurant_is/auth/internal/utils"
)

type App struct {
	log          *slog.Logger
	gRPCServ     *grpcapp.App
	db           *postgres.Storage
	notyfyclient *notifyclient.Client
	cfg          *config.Config
	telemetry    *metricsapp.App
}

func New(log *slog.Logger, cfg config.Config, tel *telemetry.Telemetry) *App {
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

	// --- Kafka init ---
	kafkaClient, err := kafkainfra.NewKafka(cfg.Kafka)
	if err != nil {
		panic(err)
	}
	userEventProducer := kafkarepo.NewUserEventProducer(kafkaClient, cfg.Kafka.Topic)

	authService := service.NewAuthService(
		log,
		clientRepo,
		stafRepo,
		// tokenRepo,
		// oauthRepo,
		notifySender,
		codeProvider,
		jwt,
		userEventProducer,
	)
	audit := service.NewAuditService(auditRepo, log)
	token := service.NewTokenService(tokenRepo, jwt, log)
	user := service.NewUserService(clientRepo, stafRepo, notifySender, log, userEventProducer)
	gRPCServ := grpcapp.New(log, authService, audit, token, user, cfg.GRPCServer.Port)

	metr := metricsapp.New(log, tel, cfg.Telemetry.MetricsPort)

	return &App{
		log:          log,
		gRPCServ:     gRPCServ,
		db:           db,
		cfg:          &cfg,
		notyfyclient: notifyClient,
		telemetry:    metr,
	}
}

func (a *App) Run() error {
	go a.telemetry.Run()
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
	a.telemetry.Stop()
}
