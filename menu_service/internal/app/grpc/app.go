package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	menugrpc "github.com/netscrawler/Restaurant_is/menu_service/internal/infra/in/grpc/menu"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	telemetry  *telemetry.Telemetry
}

// New creates new gRPC server app.
func New(
	log *slog.Logger,
	dishService menugrpc.DishProvider,
	imageService menugrpc.ImageUrlCreator,
	port int,
	telemetry *telemetry.Telemetry,
) *App {
	// loggingOpts := []logging.Option{
	// 	logging.WithLogOnEvents(
	// 		// logging.StartCall, logging.FinishCall,
	// 		logging.PayloadReceived, logging.PayloadSent,
	// 	),
	// 	// Add any other option (check functions starting with logging.With).
	// }
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(grpccodes.Internal, "internal error")
		}),
	}

	// Создаем interceptor для метрик
	metricsInterceptor := createMetricsInterceptor(telemetry.CustomMetrics)

	// Создаем tracing interceptor
	tracingInterceptor := createTracingInterceptor(telemetry.Tracer)

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		// logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
		UnaryLoggingInterceptor(log),
		tracingInterceptor,
		metricsInterceptor,
	))

	menugrpc.Register(gRPCServer, dishService, imageService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
		telemetry:  telemetry,
	}
}

// createMetricsInterceptor создает interceptor для сбора метрик.
func createMetricsInterceptor(metrics *telemetry.CustomMetrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		start := time.Now()

		// Вызываем основной обработчик
		resp, err = handler(ctx, req)

		// Записываем базовые метрики
		duration := time.Since(start).Seconds()
		recordBasicMetrics(metrics, info.FullMethod, duration, err)

		return resp, err
	}
}

// recordBasicMetrics записывает базовые метрики без извлечения данных из запросов.
func recordBasicMetrics(
	metrics *telemetry.CustomMetrics,
	method string,
	duration float64,
	err error,
) {
	ctx := context.Background()

	switch method {
	case "/menu.MenuService/CreateDish":
		// Записываем только время выполнения без category_id
		metrics.RecordDishCreateDuration(ctx, duration, 0) // 0 как default category_id

		if err == nil {
			metrics.RecordDishCreated(ctx, 0)
		}

	case "/menu.MenuService/UpdateDish":
		metrics.RecordDishUpdateDuration(ctx, duration, 0)

		if err == nil {
			metrics.RecordDishUpdated(ctx, 0)
		}

	case "/menu.MenuService/GetDish":
		metrics.RecordDishGetDuration(ctx, duration)

	case "/menu.MenuService/ListDishes":
		metrics.RecordDishListDuration(ctx, duration, nil)
		metrics.RecordDishListed(ctx, nil, false)
	}
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(lvl), msg, fields...)
		},
	)
}

func UnaryLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		traceID := trace.SpanContextFromContext(ctx).TraceID().String()
		log.Info("REQ",
			slog.String("method", info.FullMethod),
			slog.Any("request", req),
			slog.String("trace_id", traceID),
		)

		resp, err = handler(ctx, req)
		if err != nil {
			log.Error("REQ FAIL",
				slog.String("method", info.FullMethod),
				slog.Any("error", err),
				slog.String("trace_id", traceID),
			)

			return resp, err
		}

		log.Info("RESP",
			slog.String("method", info.FullMethod),
			slog.Any("response", resp),
			slog.String("trace_id", traceID),
		)

		return resp, err
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	// Останавливаем телеметрию
	if a.telemetry != nil {
		if err := a.telemetry.Shutdown(context.Background()); err != nil {
			a.log.Error("failed to shutdown telemetry", slog.Any("error", err))
		}
	}

	a.gRPCServer.GracefulStop()
}

// Создаем tracing interceptor.
func createTracingInterceptor(tracer trace.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		span := trace.SpanFromContext(ctx)
		if !span.SpanContext().IsValid() {
			var spanCtx context.Context

			spanCtx, span = tracer.Start(ctx, info.FullMethod)
			defer span.End()

			ctx = spanCtx
		}

		span.SetAttributes(
			attribute.String("grpc.method", info.FullMethod),
			attribute.String("grpc.service", "menu.MenuService"),
		)

		resp, err = handler(ctx, req)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		return resp, err
	}
}
