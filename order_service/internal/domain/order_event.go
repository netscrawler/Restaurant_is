package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OrderEvent struct {
	occurredAt time.Time
	eventType  EventType
	payload    []byte
}

type orderEvent struct {
	OrderID  uuid.UUID `json:"order_id,omitempty"`
	OrderNum uint64    `json:"order_num,omitempty"`
	ClientID uuid.UUID `json:"client_id,omitempty"`
	Status   string    `json:"status,omitempty"`
}

func NewOrderEvent(
	orderID uuid.UUID,
	orderNum uint64,
	clientID uuid.UUID,
	status Status,
	occur time.Time,
) (*OrderEvent, error) {
	var event EventType

	switch status {
	case StatusCreated:
		event = eventCreated
	case StatusDelivered:
		event = eventFinalize
	case StatusDeclined:
		event = eventFinalize
	default:
		event = eventStatusChange
	}

	oEvent := orderEvent{
		OrderID:  orderID,
		Status:   string(status),
		OrderNum: orderNum,
		ClientID: clientID,
	}

	payload, err := json.Marshal(oEvent)
	if err != nil {
		return nil, err
	}

	return &OrderEvent{
		occurredAt: occur,
		payload:    payload,
		eventType:  event,
	}, nil
}

func (e OrderEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e *OrderEvent) Payload() []byte {
	return e.payload
}

func (e *OrderEvent) EventType() EventType {
	return e.eventType
}
