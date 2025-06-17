package dto

import "github.com/google/uuid"

type Dish struct {
	ID    uuid.UUID
	Name  string
	Price uint64
}

func NewDish(id []byte, name string, Price uint64) (*Dish, error) {
	uuid, err := uuid.ParseBytes(id)
	if err != nil {
		return nil, err
	}
	return &Dish{
		ID:    uuid,
		Name:  name,
		Price: Price,
	}, nil
}
