package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

type OAuthRepository interface {
	LinkAccount(ctx context.Context, clientID string, provider *models.OAuthProvider) error
	GetByProvider(ctx context.Context, provider, providerID string) (*models.OAuthProvider, error)
	UnlinkAccount(ctx context.Context, clientID, provider string) error
}

type OAuth struct {
	OAuthRepository
}

func NewOAuth(repo OAuthRepository) *OAuth {
	return &OAuth{
		OAuthRepository: repo,
	}
}
