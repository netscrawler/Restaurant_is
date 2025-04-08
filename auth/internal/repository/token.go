package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

// TokenRepository определяет интерфейс для работы с хранилищем токенов
type TokenRepository interface {
	// StoreRefreshToken сохраняет refresh токен в хранилище
	StoreRefreshToken(ctx context.Context, token *models.RefreshToken) error

	// GetRefreshToken получает refresh токен из хранилища
	GetRefreshToken(ctx context.Context, tokenString string) (*models.RefreshToken, error)

	// DeleteRefreshToken отзывает refresh токен (устанавливает флаг revoked)
	DeleteRefreshToken(ctx context.Context, tokenString string) error

	// DeleteAllUserTokens отзывает все токены пользователя
	DeleteAllUserTokens(ctx context.Context, userID string) error

	// LogTokenAction записывает действие с токеном в аудит
	LogTokenAction(ctx context.Context, userID, userType, action, ipAddress, userAgent string) error

	// CleanupExpiredTokens удаляет все истекшие токены, возвращает количество удаленных токенов
	CleanupExpiredTokens(ctx context.Context) (int64, error)
}

type Token struct {
	TokenRepository
}

func NewToken(repo TokenRepository) *Token {
	return &Token{
		TokenRepository: repo,
	}
}
