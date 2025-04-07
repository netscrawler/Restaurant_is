package pgrepo

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

/* CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
); */

func (p *pgClient) CreateClient(ctx context.Context, client *models.Client) error {
	const op = "repository.pg.Client.Create"

	query, args, err := p.db.Builder.Insert("clients").Columns("email", "password_hash").
		Values(client.Email, client.PasswordHash).ToSql()
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

func (p *pgClient) GetClientByEmail(ctx context.Context, email string) (*models.Client, error) {
	const op = "repository.pg.Client.Create"

	query, args, err := p.db.Builder.Select(
		"id",
		"email",
		"password_hash",
		"is_active",
		"created_at",
		"updated_at",
	).
		From("clients").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		p.log.Error(op+"failed to build SQL query", zap.Error(err))

		return nil, err
	}

	row := p.db.DB.QueryRow(ctx, query, args...)

	client := models.Client{}
	err = row.Scan(
		&client.ID,
		&client.Email,
		&client.PasswordHash,
		&client.IsActive,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Info(op+"client not found", zap.String("email", email))

			return nil, err
		}

		p.log.Error(op+"failed to execute SQL query", zap.Error(err))

		return nil, err
	}

	return &client, nil
}

func (p *pgClient) UpdateClientPassword(ctx context.Context, email string, newHash string) error {
	const op = "repository.pg.Client.UpdateClientPassword"

	query, args, err := p.db.Builder.Update("clients").
		Set("password_hash", newHash).
		Set("updated_at", "NOW()").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		p.log.Error(op+"failed to build SQL query", zap.Error(err))

		return err
	}

	if _, err := p.db.DB.Exec(ctx, query, args...); err != nil {
		p.log.Error(op+"failed to execute SQL query", zap.Error(err))

		return err
	}

	return nil
}

func (p *pgClient) DeactivateClient(ctx context.Context, email string) error {
	const op = "repository.pg.Client.DeactivateClient"

	query, args, err := p.db.Builder.Update("clients").
		Set("is_active", false).
		Where(squirrel.Eq{"email": email}).
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

	p.log.Info(op+"client deactivated", zap.String("email", email))

	return nil
}
