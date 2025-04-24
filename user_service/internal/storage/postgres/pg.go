package postgres

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	log     *slog.Logger
	DB      *pgxpool.Pool
	Builder *squirrel.StatementBuilderType
}

func MustSetup(ctx context.Context, dsn string, log *slog.Logger) *Storage {
	logger := log.With("storage", "postgres")
	stLogger := logger.With("func", "setup")

	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		stLogger.ErrorContext(ctx, "error parse config", slog.String("error", err.Error()))
		panic(err)
	}

	pgConn, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		stLogger.ErrorContext(ctx, "error create client", slog.String("error", err.Error()))
		panic(err)
	}

	err = pgConn.Ping(ctx)
	if err != nil {
		stLogger.ErrorContext(ctx, "connection error", slog.String("error", err.Error()))
		panic(err)
	}

	stLogger.InfoContext(ctx, "Successfyly connect to database")

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Storage{
		log:     logger,
		DB:      pgConn,
		Builder: &builder,
	}
}

func (s *Storage) Stop(ctx context.Context) {
	logger := s.log.With("func", "stop")

	s.DB.Close()

	logger.InfoContext(ctx, "Connection to database closed")
}
