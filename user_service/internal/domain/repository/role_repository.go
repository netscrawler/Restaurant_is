package repository

import (
	"context"

	"user_service/internal/domain/models"

	"github.com/google/uuid"
)

// RoleRepository определяет интерфейс для работы с ролями
type RoleRepository interface {
	// Create создает новую роль
	Create(ctx context.Context, role *models.Role) error

	// GetByID получает роль по ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error)

	// GetByName получает роль по имени
	GetByName(ctx context.Context, name string) (*models.Role, error)

	// Update обновляет роль
	Update(ctx context.Context, role *models.Role) error

	// List возвращает список ролей
	List(ctx context.Context) ([]*models.Role, error)
}

// UserRoleRepository определяет интерфейс для работы с ролями пользователей
type UserRoleRepository interface {
	// AssignRole назначает роль пользователю
	AssignRole(ctx context.Context, userRole *models.UserRole) error

	// RevokeRole отзывает роль у пользователя
	RevokeRole(ctx context.Context, userID, roleID uuid.UUID) error

	// GetUserRoles получает роли пользователя
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*models.Role, error)

	// HasRole проверяет, есть ли у пользователя определенная роль
	HasRole(ctx context.Context, userID, roleID uuid.UUID) (bool, error)
}
