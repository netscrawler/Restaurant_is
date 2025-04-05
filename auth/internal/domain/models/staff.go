package models

import "time"

type Staff struct {
	ID           string // UUID
	WorkEmail    string
	PasswordHash string
	Position     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
