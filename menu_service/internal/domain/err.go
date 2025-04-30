package domain

import (
	"errors"
	"fmt"
)

var ErrInternal = errors.New("internal error")

var ErrInvalid = errors.New("invalid argument")

var ErrInvalidUUID = fmt.Errorf("%w: invalid uuid", ErrInvalid)
