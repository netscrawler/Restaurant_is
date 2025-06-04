package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/order_service/internal/domain"
)

type DishQuantity struct {
	DishID   uuid.UUID
	Quantity uint8
}

type Order struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Num           uint64
	Price         uint64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Status        string
	OrderType     string
	Address       string
	DishQuantites []DishQuantity
}

type Dish struct {
	ID    uuid.UUID
	Name  string
	Price uint64
}

func NewOrder(o *domain.Order) *Order {
	domainDishes := o.Dishes()
	dishes := make([]DishQuantity, len(domainDishes))
	for i, d := range domainDishes {
		dishes[i] = DishQuantity{
			DishID:   d.Dish,
			Quantity: d.Quantity,
		}
	}

	return &Order{
		ID:            o.ID(),
		UserID:        o.UserID(),
		Price:         o.Price(),
		CreatedAt:     o.CreatedAt(),
		UpdatedAt:     o.UpdatedAt(),
		Status:        string(o.Status()),
		OrderType:     string(o.OrderType()),
		Address:       o.Address(),
		DishQuantites: dishes,
	}
}
