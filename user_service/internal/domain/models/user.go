package models

import (
	"time"

	"github.com/google/uuid"
)

// User представляет доменную модель пользователя
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	FullName  string    `json:"full_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser создает нового пользователя
func NewUser(email, phone, fullName string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Phone:     phone,
		FullName:  fullName,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update обновляет данные пользователя
func (u *User) Update(email, phone, fullName string) {
	if email != "" {
		u.Email = email
	}
	if phone != "" {
		u.Phone = phone
	}
	if fullName != "" {
		u.FullName = fullName
	}
	u.UpdatedAt = time.Now()
}

// Deactivate деактивирует пользователя
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate активирует пользователя
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}
