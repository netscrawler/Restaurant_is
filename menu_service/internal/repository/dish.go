package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/models/dto"
)

type DishRepository interface {
	Create(ctx context.Context, dish *dto.Dish) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.Dish, error)
	GetByFilter(ctx context.Context, filter *dto.ListDishFilter) ([]dto.Dish, error)
	Update(ctx context.Context, dish *dto.Dish) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type dish struct {
	DishRepository
}

func NewDishRepository(repo DishRepository) *dish {
	return &dish{
		DishRepository: repo,
	}
}
