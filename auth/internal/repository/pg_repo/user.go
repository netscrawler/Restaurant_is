package pgrepo

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgUser struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgUser(db *postgres.Storage, log *zap.Logger) *pgUser {
	return &pgUser{
		log: log,
		db:  db,
	}
}

// Основные операции
func (p *pgUser) Create(ctx context.Context, user *models.User) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgUser) GetByID(ctx context.Context, id int64) (*models.User, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgUser) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgUser) UpdatePassword(ctx context.Context, id int64, newHash string) error {
	panic("not implemented") // TODO: Implement
}

func (p *pgUser) Deactivate(ctx context.Context, id int64) error {
	panic("not implemented") // TODO: Implement
}

// Валидация
func (p *pgUser) IsEmailExist(ctx context.Context, email string) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (p *pgUser) IsPhoneExist(ctx context.Context, phone string) (bool, error) {
	panic("not implemented") // TODO: Implement
}
