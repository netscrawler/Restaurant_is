package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/infra/out/postgres"
	"github.com/netscrawler/Restaurant_is/auth/internal/utils"
)

type pgClient struct {
	log *slog.Logger
	db  *postgres.Storage
}

func NewPgClient(db *postgres.Storage, log *slog.Logger) *pgClient {
	return &pgClient{
		log: log,
		db:  db,
	}
}

func (p *pgClient) CreateClient(ctx context.Context, client *models.Client) error {
	const op = "repository.pg.Client.Create"

	log := utils.LoggerWithTrace(ctx, p.log)

	query, args, err := p.db.Builder.Insert("clients").Columns("id", "phone").
		Values(client.ID, client.Phone).ToSql()
	if err != nil {
		log.Error(op+"failed to build SQL query", slog.Any("error", err))

		return fmt.Errorf("%w %w", domain.ErrBuildQuery, err)
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		log.Error(op+"failed to execute SQL query", slog.Any("error", err))

		return fmt.Errorf("%w %w", domain.ErrExecQuery, err)
	}

	return nil
}

func (p *pgClient) GetClientByPhone(ctx context.Context, phone string) (*models.Client, error) {
	const op = "repository.pg.Client.GetByPhone"

	log := utils.LoggerWithTrace(ctx, p.log)

	query, args, err := p.db.Builder.Select(
		"id",
		"phone",
		"is_active",
		"created_at",
	).
		From("clients").
		Where(squirrel.Eq{"phone": phone}).
		ToSql()
	if err != nil {
		wrapErr := fmt.Errorf("%w, %w", domain.ErrBuildQuery, err)
		log.Error(op+"failed to build SQL query", slog.Any("error", wrapErr))

		return nil, wrapErr
	}

	row := p.db.DB.QueryRow(ctx, query, args...)

	client := new(models.Client)

	err = row.Scan(
		&client.ID,
		&client.Phone,
		&client.IsActive,
		&client.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info(op+"client not found", slog.String("phone", phone))

			return nil, domain.ErrNotFound
		}

		log.Error(op+"failed to execute SQL query", slog.Any("error", err))

		return nil, fmt.Errorf("%w %w", domain.ErrExecQuery, err)
	}

	return client, nil
}

func (p *pgClient) DeactivateClient(ctx context.Context, phone string) error {
	const op = "repository.pg.Client.DeactivateClient"

	log := utils.LoggerWithTrace(ctx, p.log)

	query, args, err := p.db.Builder.Update("clients").
		Set("is_active", false).
		Where(squirrel.Eq{"phone": phone}).
		ToSql()
	if err != nil {
		log.Error(op+"failed to build SQL query", slog.Any("error", err))

		return fmt.Errorf("%w %w", domain.ErrBuildQuery, err)
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		log.Error(op+"failed to execute SQL query", slog.Any("error", err))

		return fmt.Errorf("%w %w", domain.ErrExecQuery, err)
	}

	log.Info(op+"client deactivated", slog.String("email", phone))

	return nil
}
