package models

import (
	"time"

	"github.com/google/uuid"
)

// AuthFilter представляет фильтр для поиска событий аутентификации.
type AuthFilter struct {
	UserID    *uuid.UUID       `json:"userId,omitempty"`
	UserType  *UserType        `json:"userType,omitempty"`
	Action    *AuthEventAction `json:"action,omitempty"`
	IPAddress *string          `json:"ipAddress,omitempty"`
	DateFrom  *time.Time       `json:"dateFrom,omitempty"`
	DateTo    *time.Time       `json:"dateTo,omitempty"`
	Limit     int              `json:"limit,omitempty"`
	Offset    int              `json:"offset,omitempty"`
}
