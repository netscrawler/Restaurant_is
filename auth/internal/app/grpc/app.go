package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	authgrpc "github.com/netscrawler/Restaurant_is/auth/internal/infra/in/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService authgrpc.Auth,
	audit authgrpc.Audit,
	token authgrpc.Token,
	user authgrpc.User,
	port int,
) *App {
	gRPCServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			UnaryLoggingInterceptor(log),
			recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(),
		),
	)
	authgrpc.Register(gRPCServer, authService, audit, token, user)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
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

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

//nolint:varnamelen
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

	a.log.Info("stopping gRPC server", slog.String("op", op), slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
