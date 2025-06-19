package pgrepo

import (
	"context"
	"log/slog"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/infra/out/postgres"
)

type pgToken struct {
	log *slog.Logger
	db  *postgres.Storage
}

func NewPgToken(db *postgres.Storage, log *slog.Logger) *pgToken {
	return &pgToken{
		log: log,
		db:  db,
	}
}

func (p *pgToken) StoreRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgToken) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgToken) RevokeRefreshToken(ctx context.Context, token string) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgToken) RevokeAllTokens(
	ctx context.Context,
	userID string,
	userType models.UserType,
) error {
	panic("not implemented") // TODO: Implement
}
