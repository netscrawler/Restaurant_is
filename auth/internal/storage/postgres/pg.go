package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	log     *zap.Logger
	Db      *pgxpool.Pool
	Builder *squirrel.StatementBuilderType
}

// TODO: Переделаать логирование
func MustSetup(ctx context.Context, dsn string, log *zap.Logger) *Storage {
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
	log.Info(fmt.Sprintf("%s Successfyly connect to database", op))
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Storage{
		log:     log,
		Db:      pgConn,
		Builder: &builder,
	}
}
