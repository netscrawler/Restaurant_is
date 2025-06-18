package models

import (
	"time"

	"github.com/google/uuid"
)

// Staff представляет доменную модель сотрудника
type Staff struct {
	ID        uuid.UUID `json:"id"`
	WorkEmail string    `json:"work_email"`
	WorkPhone string    `json:"work_phone"`
	FullName  string    `json:"full_name"`
	Position  string    `json:"position"`
	IsActive  bool      `json:"is_active"`
	HireDate  time.Time `json:"hire_date"`
}

// NewStaff создает нового сотрудника
func NewStaff(workEmail, workPhone, fullName, position string) *Staff {
	now := time.Now()
	return &Staff{
		ID:        uuid.New(),
		WorkEmail: workEmail,
		WorkPhone: workPhone,
		FullName:  fullName,
		Position:  position,
		IsActive:  true,
		HireDate:  now,
	}
}

// Update обновляет данные сотрудника
func (s *Staff) Update(workPhone, position string) {
	if workPhone != "" {
		s.WorkPhone = workPhone
	}
	if position != "" {
		s.Position = position
	}
}

// Deactivate деактивирует сотрудника
func (s *Staff) Deactivate() {
	s.IsActive = false
}

// Activate активирует сотрудника
func (s *Staff) Activate() {
	s.IsActive = true
}
