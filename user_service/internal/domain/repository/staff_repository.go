package repository

import (
	"context"

	"user_service/internal/domain/models"

	"github.com/google/uuid"
)

// StaffRepository определяет интерфейс для работы с сотрудниками
type StaffRepository interface {
	// Create создает нового сотрудника
	Create(ctx context.Context, staff *models.Staff) error

	// GetByID получает сотрудника по ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Staff, error)

	// GetByWorkEmail получает сотрудника по рабочему email
	GetByWorkEmail(ctx context.Context, workEmail string) (*models.Staff, error)

	// Update обновляет сотрудника
	Update(ctx context.Context, staff *models.Staff) error

	// List возвращает список сотрудников с пагинацией
	List(ctx context.Context, onlyActive bool, offset, limit int) ([]*models.Staff, error)

	// Count возвращает общее количество сотрудников
	Count(ctx context.Context, onlyActive bool) (int, error)
}
