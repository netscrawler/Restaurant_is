package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/order_service/internal/domain"
	"github.com/netscrawler/Restaurant_is/order_service/internal/models/repository"
)

type Order interface {
	Save(ctx context.Context, order *repository.Order) (uint64, error)
	Get(ctx context.Context, orderID uuid.UUID) (*domain.Order, error)
}
