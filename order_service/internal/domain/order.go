package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	Status    string
	OrderType string
)

func (s Status) Valid() bool {
	status := map[Status]struct{}{
		StatusCreated:   {},
		StatusProcess:   {},
		StatusOnKitchen: {},
		StatusDelivery:  {},
		StatusDelivered: {},
		StatusDeclined:  {},
	}

	_, ok := status[s]

	return ok
}

type DishQuantity struct {
	Dish     uuid.UUID
	Quantity uint8
}

type Order struct {
	id        uuid.UUID
	num       uint64
	user      uuid.UUID
	price     uint64
	createdAt time.Time
	updatedAt time.Time
	status    Status
	orderType OrderType
	address   string
	dishes    []DishQuantity
}

func NewOrder(
	user uuid.UUID,
	dishes DishList,
	orderType OrderType,
	address string,
) (*Order, error) {
	dishItems := make([]DishQuantity, 0, len(dishes))

	var totalPrice uint64

	for dish, quantity := range dishes {
		dishItems = append(dishItems, DishQuantity{
			Dish:     dish.ID,
			Quantity: quantity,
		})
		totalPrice += dish.Price * uint64(quantity)
	}

	if len(dishItems) == 0 {
		return nil, ErrEmptyDishList
	}

	return &Order{
		id:        uuid.New(),
		createdAt: time.Now(),
		updatedAt: time.Now(),
		status:    StatusCreated,
		price:     totalPrice,
		user:      user,
		dishes:    dishItems,
		orderType: orderType,
		address:   address,
		num:       0,
	}, nil
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) NUM() uint64 {
	return o.num
}

func (o *Order) SetNUM(n uint64) {
	if o.num != 0 {
		o.num = n
	}
}

func (o *Order) UserID() uuid.UUID {
	return o.user
}

func (o *Order) Price() uint64 {
	return o.price
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) UpdatedAt() time.Time {
	return o.updatedAt
}

func (o *Order) Status() string {
	return string(o.status)
}

func (o *Order) OrderType() string {
	return string(o.orderType)
}

func (o *Order) Address() string {
	return o.address
}

func (o *Order) Dishes() []DishQuantity {
	dishes := make([]DishQuantity, len(o.dishes))
	copy(dishes, o.dishes)

	return dishes
}

func (o *Order) statusChange(s Status) *OrderEvent {
	o.status = s
	o.updatedAt = time.Now()
	event, _ := NewOrderEvent(
		o.id,
		o.num,
		o.user,
		s,
		o.updatedAt,
	)

	return event
}

func (o *Order) Process() (*OrderEvent, error) {
	if o.status != StatusCreated {
		return nil, ErrInvalidStatusTransition
	}

	return o.statusChange(StatusProcess), nil
}

func (o *Order) SendToKitchen() (*OrderEvent, error) {
	if o.status != StatusProcess {
		return nil, ErrInvalidStatusTransition
	}

	return o.statusChange(StatusOnKitchen), nil
}

func (o *Order) StartDelivery() (*OrderEvent, error) {
	if o.status != StatusOnKitchen {
		return nil, ErrInvalidStatusTransition
	}

	return o.statusChange(StatusDelivery), nil
}

func (o *Order) Complete() (*OrderEvent, error) {
	if o.status != StatusDelivery {
		return nil, ErrInvalidStatusTransition
	}

	return o.statusChange(StatusDelivered), nil
}

func (o *Order) Decline() (*OrderEvent, error) {
	if o.status == StatusDelivered || o.status == StatusDeclined {
		return nil, ErrInvalidStatusTransition
	}

	return o.statusChange(StatusDeclined), nil
}
