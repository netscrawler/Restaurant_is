package domain

import "errors"

// Repository errors.
var (
	ErrNotFound   = errors.New("ErrNotFound")
	ErrBuildQuery = errors.New("ErrBuildQuery")
	ErrExecQuery  = errors.New("ErrExecQuery")
	ErrScanRow    = errors.New("ErrScanRow")
)

// Service errors.
var (
	ErrInternal         = errors.New("InternalError")
	ErrFailedCreateCode = errors.New("ErrGenerateCode")
	ErrGenerateToken    = errors.New("ErrGenerateToken")
	ErrUnknownUserType  = errors.New("ErrUnknownUserType")
	ErrUserAlreadyExist = errors.New("UserWithThisEmailAlreadyExist")
	ErrGeneratePassword = errors.New("ErrGeneratePassword")
	ErrUserDeactivated  = errors.New("ErrUserDeactivated")
)

// Token errors.
var (
	ErrUnexpectedSignMethod = errors.New("ErrUnexpectedSignMethod")
	ErrInternalCodeGen      = errors.New("ErrInternalCodeGen")
	ErrInternalCodeParse    = errors.New("ErrInternalCodeParse")
	ErrSecret               = errors.New("ErrSecret")
	ErrEmptySecret          = errors.New("ErrEmptySecret")
)

var (
	ErrCodeConfirmGen = errors.New("ErrCodeConfirmGen")
	ErrMissingKeys    = errors.New("ErrMissingRSAKeys")
)

// Ошибки валидации.
var (
	ErrEmptyField      = errors.New("field cannot be empty")
	ErrInvalidPhone    = errors.New("invalid phone number format")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidCode     = errors.New("invalid confirmation code")
	ErrInvalidURI      = errors.New("invalid redirect URI")
	ErrInvalidToken    = errors.New("invalid token")

	ErrPasswordLen       = errors.New("password must be at least 8 characters long")
	ErrPasswordDigit     = errors.New("password must contain at least one digit")
	ErrPasswordUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordLowerCase = errors.New("password must contain at least one lowercase letter")

	ErrNilRequest = errors.New("request cannot be nil")

	ErrInvalidUserUUID = errors.New("ErrInvalidUserUUID")
)
