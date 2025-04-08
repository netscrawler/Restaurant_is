package models

import "time"

// RefreshToken представляет хранимый refresh токен.
type RefreshToken struct {
	UserID       string    // ID пользователя
	UserType     string    // Тип пользователя: "client" или "staff"
	RefreshToken string    // Сам токен
	ExpiresAt    time.Time // Время истечения токена
	Revoked      bool      // Флаг отзыва токена
	CreatedAt    time.Time // Время создания
}
