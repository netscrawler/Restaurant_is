package repository

import "time"

type Event struct {
	ID        string
	Type      string
	Payload   []byte
	Published bool
	OccuredAt time.Time
}

func NewEvent(id, eventType string, payload []byte) *Event {
	return &Event{
		ID:        id,
		Type:      eventType,
		Payload:   payload,
		Published: false,
		OccuredAt: time.Now(),
	}
}
