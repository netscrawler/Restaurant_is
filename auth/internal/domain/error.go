package domain

import "errors"

// Repository errors.
var (
	ErrNotFound   = errors.New("ErrNotFound")
	ErrBuildQuery = errors.New("ErrBuildQuery")
	ErrExecQuery  = errors.New("ErrExecQuery")
)

// Service errors.
var (
	ErrInternal         = errors.New("InternalError")
	ErrFailedCreateCode = errors.New("ErrGenerateCode")
)

// Token errors.
var (
	ErrInvalidToken         = errors.New("ErrInvalidToken")
	ErrUnexpectedSignMethod = errors.New("ErrUnexpectedSignMethod")
	ErrInternalCodeGen      = errors.New("ErrInternalCodeGen")
	ErrInternalCodeParse    = errors.New("ErrInternalCodeParse")
)

var ErrCodeConfirmGen = errors.New("ErrCodeConfirmGen")
