package domain

import (
	"time"
)

type EventType string

type Event interface {
	OccurredAt() time.Time
	Payload() []byte
	EventType() EventType
}
