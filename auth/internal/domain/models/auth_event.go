package models

import "time"

type AuthEvent struct {
	UserID    string
	UserType  UserType
	Action    string // login, logout, refresh
	IP        string
	UserAgent string
	Success   bool
	CreatedAt time.Time
}
