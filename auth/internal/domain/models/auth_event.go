package models

import (
	"time"

	"github.com/google/uuid"
)

// AuthEventAction определяет типы действий для аудита.
type AuthEventAction string

const (
	// Действия аутентификации.
	ActionLogin        AuthEventAction = "login"
	ActionLogout       AuthEventAction = "logout"
	ActionTokenRefresh AuthEventAction = "token_refresh"
	ActionTokenRevoke  AuthEventAction = "token_revoke"
)

// AuthEvent представляет событие аутентификации.
type AuthEvent struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"userId"`
	UserType  UserType        `json:"userType"`
	Action    AuthEventAction `json:"action"`
	IPAddress string          `json:"ipAddress"`
	UserAgent string          `json:"userAgent"`
	CreatedAt time.Time       `json:"createdAt"`
}
