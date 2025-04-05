package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

type OAuthProviderRepository interface {
	LinkAccount(ctx context.Context, userID int64, provider *models.OAuthProvider) error
	GetByProviderID(
		ctx context.Context,
		providerName, providerID string,
	) (*models.OAuthProvider, error)
	GetUserProviders(ctx context.Context, userID int64) ([]*models.OAuthProvider, error)
	UnlinkAccount(ctx context.Context, userID int64, providerName string) error
}
