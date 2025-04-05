package models

import "time"

type AuthFilter struct {
	UserID   *string
	UserType *UserType
	Actions  []string
	DateFrom *time.Time
	DateTo   *time.Time
}
