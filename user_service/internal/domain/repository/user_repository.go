package repository

import (
	"context"

	"github.com/google/uuid"
	"user_service/internal/domain/models"
)

// UserRepository определяет интерфейс для работы с пользователями.
type UserRepository interface {
	// Create создает нового пользователя
	Create(ctx context.Context, user *models.User) error

	// GetByID получает пользователя по ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)

	// GetByEmail получает пользователя по email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetByPhone получает пользователя по телефону
	GetByPhone(ctx context.Context, phone string) (*models.User, error)

	// Update обновляет пользователя
	Update(ctx context.Context, user *models.User) error

	// Delete удаляет пользователя
	Delete(ctx context.Context, id uuid.UUID) error

	// List возвращает список пользователей с пагинацией
	List(ctx context.Context, onlyActive bool, offset, limit int) ([]*models.User, error)

	// Count возвращает общее количество пользователей
	Count(ctx context.Context, onlyActive bool) (int, error)
}
