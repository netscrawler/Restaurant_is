package models

import "time"

type User struct {
	ID           int64
	AccountType  string
	Email        string
	WorkEmail    string
	Phone        string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
