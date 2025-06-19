package models

import (
	"time"

	"github.com/google/uuid"
)

// Role представляет доменную модель роли.
type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewRole создает новую роль.
func NewRole(name, description string) *Role {
	now := time.Now()

	return &Role{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update обновляет данные роли.
func (r *Role) Update(name, description string) {
	r.Name = name
	r.Description = description
	r.UpdatedAt = time.Now()
}

// UserRole представляет связь пользователя с ролью.
type UserRole struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	RoleID     uuid.UUID `json:"role_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

// NewUserRole создает новую связь пользователя с ролью.
func NewUserRole(userID, roleID uuid.UUID) *UserRole {
	return &UserRole{
		ID:         uuid.New(),
		UserID:     userID,
		RoleID:     roleID,
		AssignedAt: time.Now(),
	}
}
