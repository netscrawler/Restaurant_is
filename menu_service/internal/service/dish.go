package service

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

type Dish struct {
	repo DishRepository
}

func NewDishService(repo DishRepository) *Dish {
	return &Dish{
		repo: repo,
	}
}

// Create and save new dish.
func (d *Dish) Create(ctx context.Context, dish *dto.Dish) (*dto.Dish, error) {
	err := d.repo.Create(ctx, dish)
	if err != nil {
		return nil, err
	}

	return dish, nil
}

// Get dish by uuid.
func (d *Dish) Get(ctx context.Context, dishID uuid.UUID) (*dto.Dish, error) {
	dish, err := d.repo.GetByID(ctx, dishID)
	if err != nil {
		return nil, err
	}
	return dish, nil
}

func (d *Dish) Update(ctx context.Context, req *dto.UpdateDishReq) (*dto.Dish, error) {
	existingDish, err := d.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existingDish.Name = *req.Name
	}
	if req.Description != nil {
		existingDish.Description = *req.Description
	}
	if req.Price != nil {
		existingDish.Price = *req.Price
	}
	if req.CategoryID != nil {
		existingDish.CategoryID = *req.CategoryID
	}
	if req.CookingTimeMin != nil {
		existingDish.CookingTimeMin = *req.CookingTimeMin
	}
	if req.ImageURL != nil {
		existingDish.ImageURL = *req.ImageURL
	}
	if req.IsAvailable != nil {
		existingDish.IsAvailable = *req.IsAvailable
	}
	if req.Calories != nil {
		existingDish.Calories = *req.Calories
	}

	err = d.repo.Update(ctx, existingDish)
	if err != nil {
		return nil, err
	}

	return existingDish, nil
}

func (d *Dish) List(ctx context.Context, filter *dto.ListDishFilter) ([]dto.Dish, error) {
	dishes, err := d.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return dishes, nil
}
