package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

// ClientRepository - работа с клиентами
type ClientRepository interface {
	CreateClient(ctx context.Context, client *models.Client) error
	GetClientByPhone(ctx context.Context, phone string) (*models.Client, error)
	DeactivateClient(ctx context.Context, phone string) error
}

type Client struct {
	ClientRepository
}

func NewClient(repo ClientRepository) *Client {
	return &Client{
		ClientRepository: repo,
	}
}
