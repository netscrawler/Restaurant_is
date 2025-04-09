package pgrepo

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgToken struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgToken(db *postgres.Storage, log *zap.Logger) *pgToken {
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
