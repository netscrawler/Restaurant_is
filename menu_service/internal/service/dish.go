package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/models/dto"
)

type Dish struct{}

// Create and save new dish.
func (d *Dish) Create(ctx context.Context, dish *dto.Dish) (*dto.Dish, error) {
	panic("not implemented") // TODO: Implement
}

// Get dish by uuid.
func (d *Dish) Get(ctx context.Context, dishID uuid.UUID) (*dto.Dish, error) {
	panic("not implemented") // TODO: Implement
}

// Update dish in storage.
func (d *Dish) Update(ctx context.Context, dish *dto.UpdateDishReq) (*dto.Dish, error) {
	panic("not implemented") // TODO: Implement
}

// List returns dish by filter, offset and limit.
func (d *Dish) List(ctx context.Context, filter *dto.ListDishReq) ([]dto.Dish, error) {
	panic("not implemented") // TODO: Implement
}
