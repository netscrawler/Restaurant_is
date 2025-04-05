package pgrepo

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgOauth struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgOauth(db *postgres.Storage, log *zap.Logger) *pgOauth {
	return &pgOauth{
		log: log,
		db:  db,
	}
}

func (p *pgOauth) Create(ctx context.Context, token *models.RefreshToken) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) Revoke(ctx context.Context, token string) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) RevokeAllForUser(ctx context.Context, userID int64) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) DeleteExpired(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) LinkAccount(
	ctx context.Context,
	userID int64,
	provider *models.OAuthProvider,
) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) GetByProviderID(
	ctx context.Context,
	providerName string,
	providerID string,
) (*models.OAuthProvider, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) GetUserProviders(
	ctx context.Context,
	userID int64,
) ([]*models.OAuthProvider, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) UnlinkAccount(ctx context.Context, userID int64, providerName string) error {
	panic("not implemented") // TODO: Implement
}
