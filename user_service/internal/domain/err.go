package domain

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrStaffNotFound    = errors.New("staff not found")
	ErrRoleNotFound     = errors.New("role not found")
	ErrUserRoleNotFound = errors.New("user role not found")
)
