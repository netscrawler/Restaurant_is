package models

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID        uuid.UUID // UUID
	Phone     string
	IsActive  bool
	CreatedAt time.Time
}

func NewClient(phone string) *Client {
	defer func() { recover() }()

	uid := uuid.New()

	return &Client{
		ID:        uid,
		Phone:     phone,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
}
