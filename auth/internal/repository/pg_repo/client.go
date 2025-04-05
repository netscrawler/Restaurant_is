package pgrepo

import (
	"context"

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
	panic("not implemented") // TODO: Implement
}

func (p *pgClient) GetClientByEmail(ctx context.Context, email string) (*models.Client, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgClient) UpdateClientPassword(ctx context.Context, email string, newHash string) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgClient) DeactivateClient(ctx context.Context, email string) error {
	panic("not implemented") // TODO: Implement
}
