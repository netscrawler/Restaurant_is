package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken представляет хранимый refresh токен.
type RefreshToken struct {
	UserID       uuid.UUID // ID пользователя
	UserType     UserType  // Тип пользователя: "client" или "staff"
	RefreshToken string    // Сам токен
	ExpiresAt    time.Time // Время истечения токена
	Revoked      bool      // Флаг отзыва токена
	CreatedAt    time.Time // Время создания
}

func NewRefreshToken(
	userID uuid.UUID,
	userType UserType,
	token string,
	expiresAt time.Time,
) *RefreshToken {
	return &RefreshToken{
		UserID:       userID,
		UserType:     userType,
		RefreshToken: token,
		ExpiresAt:    expiresAt,
		Revoked:      false,
		CreatedAt:    time.Now(),
	}
}
