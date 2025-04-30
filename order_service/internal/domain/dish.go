package domain

import (
	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/order_service/internal/models/dto"
)

type Dish struct {
	ID    uuid.UUID
	Name  string
	Price uint64
}

func NewDish(d dto.Dish) *Dish {
	return &Dish{
		ID:    d.ID,
		Name:  d.Name,
		Price: d.Price,
	}
}

type DishList map[Dish]uint8

func (d DishList) Validate() error {
	if len(d) == 0 {
		return ErrEmptyDishList
	}

	return nil
}

func (d DishList) CalculatePrice() uint64 {
	var c uint64 = 0
	for k, v := range d {
		c += k.Price * uint64(v)
	}

	return c
}
