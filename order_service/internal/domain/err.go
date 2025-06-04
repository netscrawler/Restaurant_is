package domain

import "errors"

var (
	ErrEmptyDishList           = errors.New("empty dish list")
	ErrInvalidStatus           = errors.New("invalid status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
