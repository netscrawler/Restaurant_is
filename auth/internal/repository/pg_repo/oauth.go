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

func (p *pgOauth) LinkAccount(
	ctx context.Context,
	clientID string,
	provider *models.OAuthProvider,
) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) GetByProvider(
	ctx context.Context,
	provider string,
	providerID string,
) (*models.OAuthProvider, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgOauth) UnlinkAccount(ctx context.Context, clientID string, provider string) error {
	panic("not implemented") // TODO: Implement
}
