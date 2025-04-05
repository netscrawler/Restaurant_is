package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

type UserRepository interface {
	// Основные операции
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdatePassword(ctx context.Context, id int64, newHash string) error
	Deactivate(ctx context.Context, id int64) error

	// Валидация
	IsEmailExist(ctx context.Context, email string) (bool, error)
	IsPhoneExist(ctx context.Context, phone string) (bool, error)
}
