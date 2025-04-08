package pgrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgClient struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgUser(db *postgres.Storage, log *zap.Logger) *pgClient {
	return &pgClient{
		log: log,
		db:  db,
	}
}

func (p *pgClient) CreateClient(ctx context.Context, client *models.Client) error {
	const op = "repository.pg.Client.Create"

	query, args, err := p.db.Builder.Insert("clients").Columns("phone").
		Values(client.Phone).ToSql()
	if err != nil {
		p.log.Error(op+"failed to build SQL query", zap.Error(err))

		return err
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		p.log.Error(op+"failed to execute SQL query", zap.Error(err))

		return err
	}

	return nil
}

func (p *pgClient) GetClientByPhone(ctx context.Context, phone string) (*models.Client, error) {
	const op = "repository.pg.Client.GetByPhone"

	query, args, err := p.db.Builder.Select(
		"id",
		"email",
		"password_hash",
		"is_active",
		"created_at",
		"updated_at",
	).
		From("clients").
		Where(squirrel.Eq{"phone": phone}).
		ToSql()
	if err != nil {
		wrapErr := fmt.Errorf("%w, %w", domain.ErrBuildQuery, err)
		p.log.Error(op+"failed to build SQL query", zap.Error(wrapErr))

		return nil, wrapErr
	}

	row := p.db.DB.QueryRow(ctx, query, args...)

	client := models.Client{}
	err = row.Scan(
		&client.ID,
		&client.Phone,
		&client.IsActive,
		&client.CreatedAt,
	)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {

			p.log.Info(op+"client not found", zap.String("phone", phone))

			return nil, domain.ErrNotFound
		}

		p.log.Error(op+"failed to execute SQL query", zap.Error(err))

		return nil, err
	}

	return &client, nil
}

func (p *pgClient) DeactivateClient(ctx context.Context, phone string) error {
	const op = "repository.pg.Client.DeactivateClient"

	query, args, err := p.db.Builder.Update("clients").
		Set("is_active", false).
		Where(squirrel.Eq{"phone": phone}).
		ToSql()
	if err != nil {
		p.log.Error(op+"failed to build SQL query", zap.Error(err))

		return err
	}
	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {

		p.log.Error(op+"failed to execute SQL query", zap.Error(err))

		return err
	}

	p.log.Info(op+"client deactivated", zap.String("email", phone))

	return nil
}
