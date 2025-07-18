package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	service "user_service/internal/domain/app"
	usergrpc "user_service/internal/infra/in/grpc"
	"user_service/internal/telemetry"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	userAppService *service.UserAppService,
	staffAppService *service.StaffAppService,
	roleAppService *service.RoleAppService,
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

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),

		grpc.StatsHandler(otelgrpc.NewServerHandler()), grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			UnaryLoggingInterceptor(log),

			// logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
		),
		// grpc.ChainStreamInterceptor(
		// 	otelgrpc.StreamServerInterceptor(),
		// ),
	)

	usergrpc.Register(gRPCServer, userAppService, staffAppService, roleAppService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
		telemetry:  telemetry,
	}
}

func UnaryLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		traceID := oteltrace.SpanContextFromContext(ctx).TraceID().String()
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

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(lvl), msg, fields...)
		},
	)
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

	a.gRPCServer.GracefulStop()
}
