package metrics

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/netscrawler/Restaurant_is/order_service/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	log       *slog.Logger
	server    *http.Server
	port      int
	telemetry *telemetry.Telemetry
}

// New создает новое приложение для метрик.
func New(log *slog.Logger, telemetry *telemetry.Telemetry, port int) *App {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &App{
		log:       log,
		server:    server,
		port:      port,
		telemetry: telemetry,
	}
}

// MustRun запускает HTTP сервер метрик и паникует при ошибке.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run запускает HTTP сервер метрик.
func (a *App) Run() error {
	const op = "metrics.Run"

	a.log.Info("metrics server started", slog.String("addr", a.server.Addr))

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop останавливает HTTP сервер метрик.
func (a *App) Stop() {
	const op = "metrics.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping metrics server", slog.Int("port", a.port))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("failed to shutdown metrics server", slog.Any("error", err))
	}
}
