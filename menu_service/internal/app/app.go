package app

import (
	"log/slog"

	"github.com/netscrawler/Restaurant_is/menu_service/internal/config"
)

type App struct{}

func New(log *slog.Logger, cfg *config.Config) *App {
	return &App{}
}

func (a *App) MustRun() {
	panic("err")
}

func (a *App) Stop() {
	panic("err")
}
