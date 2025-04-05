package pgrepo

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgAudit struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgAudit(db *postgres.Storage, log *zap.Logger) *pgAudit {
	return &pgAudit{
		log: log,
		db:  db,
	}
}

func (p *pgAudit) LogAuthEvent(ctx context.Context, event *models.AuthEvent) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgAudit) GetAuthEvents(
	ctx context.Context,
	filter models.AuthFilter,
) ([]*models.AuthEvent, error) {
	panic("not implemented") // TODO: Implement
}
