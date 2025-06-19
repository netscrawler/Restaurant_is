package grpcapp

import (
	"fmt"
	"net"

	notifygrpc "notify/internal/grpc/notify"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *zap.Logger, notifyService notifygrpc.NotifySender, port int) *App {
	grpc_zap.ReplaceGrpcLoggerV2(log)

	//nolint:exhaustive
	logOpts := []grpc_zap.Option{
		grpc_zap.WithLevels(func(code codes.Code) zapcore.Level {
			switch code {
			case codes.OK:
				return zapcore.DebugLevel
			case codes.InvalidArgument, codes.NotFound:
				return zapcore.WarnLevel
			default:
				return zapcore.ErrorLevel
			}
		}),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) error {
			log.Error("Recovered from panic", zap.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(log, logOpts...),
			recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc.ChainStreamInterceptor(
			grpc_zap.StreamServerInterceptor(log, logOpts...),
			recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)
	notifygrpc.Register(gRPCServer, notifyService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
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

	a.log.Info("grpc server started", zap.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(zap.String("op", op)).
		Info("stopping gRPC server", zap.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
