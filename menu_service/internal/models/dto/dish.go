package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDishReq struct {
	Name           string
	Description    string
	Price          float64
	CategoryID     uuid.UUID
	CookingTimeMin int32
	ImageURL       *string
	IsAvailable    bool
	Calories       *int32
}

type UpdateDishReq struct {
	ID uuid.UUID
	CreateDishReq
}

type ListDishReq struct {
	CategoryID    *uuid.UUID
	OnlyAvailable bool
	Page          int32
	PageSize      int32
}

type Dish struct {
	ID             uuid.UUID
	Name           string
	Description    string
	Price          float64
	CategoryID     uuid.UUID
	CookingTimeMin int32
	ImageURL       string
	IsAvailable    bool
	Calories       int32
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}
