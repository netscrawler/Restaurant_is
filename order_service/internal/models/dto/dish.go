package dto

import "github.com/google/uuid"

type Dish struct {
	ID    uuid.UUID
	Name  string
	Price uint64
}
