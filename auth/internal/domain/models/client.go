package models

import "time"

type Client struct {
	ID           string // UUID
	Email        string
	PasswordHash string
	Phone        string
	FullName     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
