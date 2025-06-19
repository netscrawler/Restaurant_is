package app

import (
	"log/slog"
	botapp "notify/internal/app/bot"
	grpcapp "notify/internal/app/grpc"
	"notify/internal/config"
	service "notify/internal/service/stub"
)

type App struct {
	botApp  *botapp.Bot
	grpcApp *grpcapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	log = log.With("channel", "telegram chat")

	bot, err := botapp.New(log, cfg.Bot.TelegramToken, cfg.Bot.BotPoll)
	if err != nil {
		panic(err)
	}

	notifyer := service.NewTelegramSender(log, bot.Bot, cfg.StubRecipient)
	grpc := grpcapp.New(log, notifyer, cfg.GRRPC.Port)
	return &App{
		botApp:  bot,
		grpcApp: grpc,
	}
}

func (a *App) MustRun() {
	a.botApp.Start()
	a.grpcApp.MustRun()
}

func (a *App) Stop() {
	a.grpcApp.Stop()
	a.botApp.Stop()
}
