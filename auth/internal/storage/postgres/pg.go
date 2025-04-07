package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	log     *zap.Logger
	DB      *pgxpool.Pool
	Builder *squirrel.StatementBuilderType
}

// TODO: Переделаать логирование.
func MustSetup(ctx context.Context, dsn string, log *zap.Logger) *Storage {
	//nolint: varnamelen
	const op = "storage.postgres.Setup"

	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Debug(op, zap.String("error", err.Error()))
		panic(err)
	}

	pgConn, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		log.Debug(op, zap.String("error", err.Error()))
		panic(err)
	}

	err = pgConn.Ping(ctx)
	if err != nil {
		log.Debug(op, zap.String("error", err.Error()))
		panic(err)
	}

	log.Info(op + "Successfyly connect to database")

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Storage{
		log:     log,
		DB:      pgConn,
		Builder: &builder,
	}
}

func (s *Storage) Stop() {
	const op = "storage.pg.Stop"

	s.DB.Close()

	s.log.Info(op + "Connection to database closed")
}
