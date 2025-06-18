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

type ImageURLProvider interface {
	// GetDownloadURL генерирует pre-signed URL для скачивания изображения.
	GetDownloadURL(ctx context.Context, objectKey string) (string, error)
}

type Dish struct {
	repo          DishRepository
	imageProvider ImageURLProvider
}

func NewDishService(repo DishRepository, imageProvider ImageURLProvider) *Dish {
	return &Dish{
		repo:          repo,
		imageProvider: imageProvider,
	}
}

// Create and save new dish.
func (d *Dish) Create(ctx context.Context, dish *dto.Dish) (*dto.Dish, error) {
	// Если есть objectKey в ImageURL, генерируем URL для скачивания
	var downloadURL string
	if dish.ImageURL != "" {
		var err error
		downloadURL, err = d.imageProvider.GetDownloadURL(ctx, dish.ImageURL)
		if err != nil {
			return nil, err
		}
	}

	err := d.repo.Create(ctx, dish)
	if err != nil {
		return nil, err
	}
	dish.ImageURL = downloadURL

	return dish, nil
}

// Get dish by uuid.
func (d *Dish) Get(ctx context.Context, dishID uuid.UUID) (*dto.Dish, error) {
	dish, err := d.repo.GetByID(ctx, dishID)
	if err != nil {
		return nil, err
	}

	if dish.ImageURL != "" {
		downloadURL, err := d.imageProvider.GetDownloadURL(ctx, dish.ImageURL)
		if err != nil {
			return nil, err
		}
		dish.ImageURL = downloadURL
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
	if existingDish.ImageURL != "" {
		downloadURL, err := d.imageProvider.GetDownloadURL(ctx, existingDish.ImageURL)
		if err != nil {
			return nil, err
		}
		existingDish.ImageURL = downloadURL
	}

	return existingDish, nil
}

func (d *Dish) List(ctx context.Context, filter *dto.ListDishFilter) ([]dto.Dish, error) {
	dishes, err := d.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	var editedDishes []dto.Dish

	for _, dish := range dishes {
		if dish.ImageURL != "" {
			downloadURL, err := d.imageProvider.GetDownloadURL(ctx, dish.ImageURL)
			if err != nil {
				return nil, err
			}
			dish.ImageURL = downloadURL
		}
		editedDishes = append(editedDishes, dish)
	}

	return editedDishes, nil
}
