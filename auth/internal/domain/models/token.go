package models

import "time"

type RefreshToken struct {
	Token     string
	UserID    string
	UserType  UserType // client/staff
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}
