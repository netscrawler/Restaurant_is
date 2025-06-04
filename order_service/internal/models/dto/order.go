package dto

import (
	"github.com/google/uuid"
)

type OrderToCreate struct {
	UserID          uuid.UUID
	OrderType       string
	DeliveryAddress []byte
	Items           []OrderItem
}

type OrderCreated struct {
	ID     []byte
	NUM    uint64
	UserID []byte
	Total  uint64
	Status string
}

func NewOrderCreated(
	id uuid.UUID,
	num uint64,
	userID uuid.UUID,
	total uint64,
	status string,
) (*OrderCreated, error) {
	oID, err := id.MarshalBinary()
	if err != nil {
	}

	uID, err := userID.MarshalBinary()
	if err != nil {
	}

	return &OrderCreated{
		ID:     oID,
		NUM:    num,
		UserID: uID,
		Total:  total,
		Status: status,
	}, nil
}

func NewOrder(
	userID uuid.UUID,
	orderType string,
	deliveryAddress []byte,
	items []OrderItem,
) *OrderToCreate {
	return &OrderToCreate{
		UserID:          userID,
		OrderType:       orderType,
		DeliveryAddress: deliveryAddress,
		Items:           items,
	}
}

type OrderItem struct {
	Item    uuid.UUID
	Quanity uint8
}
