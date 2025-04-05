package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

// TokenRepository - управление токенами
type TokenRepository interface {
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllTokens(ctx context.Context, userID string, userType models.UserType) error
}

type Token struct {
	TokenRepository
}

func NewToken(repo TokenRepository) *Token {
	return &Token{
		TokenRepository: repo,
	}
}
