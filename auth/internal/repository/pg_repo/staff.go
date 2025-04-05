package pgrepo

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgStaff struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgStaff(db *postgres.Storage, log *zap.Logger) *pgStaff {
	return &pgStaff{
		log: log,
		db:  db,
	}
}

func (p *pgStaff) CreateStaff(ctx context.Context, staff *models.Staff) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgStaff) GetStaffByEmail(ctx context.Context, workEmail string) (*models.Staff, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgStaff) UpdateStaffPassword(ctx context.Context, workEmail string, newHash string) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgStaff) DeactivateStaff(ctx context.Context, workEmail string) error {
	panic("not implemented") // TODO: Implement
}
