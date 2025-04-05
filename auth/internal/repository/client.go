package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

// ClientRepository - работа с клиентами
type ClientRepository interface {
	CreateClient(ctx context.Context, client *models.Client) error
	GetClientByEmail(ctx context.Context, email string) (*models.Client, error)
	UpdateClientPassword(ctx context.Context, email, newHash string) error
	DeactivateClient(ctx context.Context, email string) error
}

type Client struct {
	ClientRepository
}

func NewClient(repo ClientRepository) *Client {
	return &Client{
		ClientRepository: repo,
	}
}
