package utils

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	authv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
)

// Регулярные выражения для валидации
var (
	// Валидация телефона (международный формат, например +79991234567)
	phoneRegex = regexp.MustCompile(`^\+[1-9]\d{10,14}$`)

	// Базовая валидация email
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// Код подтверждения (обычно 4-6 цифр)
	codeRegex = regexp.MustCompile(`^\d{4}$`)
)

// RequestValidator предоставляет API для валидации всех запросов
type RequestValidator struct{}

// NewRequestValidator создает новый экземпляр валидатора запросов
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

// ValidateRequest валидирует любой запрос из proto-сервиса
func (v *RequestValidator) ValidateRequest(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *authv1.LoginClientInitRequest:
		return ValidateLoginClientRequest(r)
	case *authv1.LoginClientConfirmRequest:
		return ValidateLoginClientConfirmRequest(r)
	case *authv1.LoginStaffRequest:
		return ValidateLoginStaffRequest(r)
	case *authv1.OAuthYandexRequest:
		return ValidateOAuthYandexRequest(r)
	case *authv1.ValidateRequest:
		return ValidateValidateRequest(r)
	case *authv1.RefreshRequest:
		return ValidateRefreshRequest(r)
	default:
		return nil // Неизвестный тип запроса
	}
}

// ValidateLoginClientRequest валидирует запрос на логин клиента
func ValidateLoginClientRequest(req *authv1.LoginClientInitRequest) error {
	if req == nil {
		return domain.ErrNilRequest
	}

	// Валидация телефона
	if req.GetPhone() == "" {
		return fmt.Errorf("phone: %w", domain.ErrEmptyField)
	}

	if !phoneRegex.MatchString(req.GetPhone()) {
		return fmt.Errorf("phone: %w", domain.ErrInvalidPhone)
	}

	return nil
}

// ValidateLoginClientConfirmRequest валидирует запрос на подтверждение логина клиента
func ValidateLoginClientConfirmRequest(req *authv1.LoginClientConfirmRequest) error {
	if req == nil {
		return domain.ErrNilRequest
	}

	// Валидация телефона
	if req.GetPhone() == "" {
		return fmt.Errorf("phone: %w", domain.ErrEmptyField)
	}

	if !phoneRegex.MatchString(req.GetPhone()) {
		return fmt.Errorf("phone: %w", domain.ErrInvalidPhone)
	}

	// Валидация кода
	if req.GetCode() == "" {
		return fmt.Errorf("code: %w", domain.ErrEmptyField)
	}

	if !codeRegex.MatchString(req.GetCode()) {
		return fmt.Errorf("code: %w", domain.ErrInvalidCode)
	}

	return nil
}

// ValidateLoginStaffRequest валидирует запрос на логин персонала
func ValidateLoginStaffRequest(req *authv1.LoginStaffRequest) error {
	if req == nil {
		return domain.ErrNilRequest
	}

	// Валидация рабочего email
	if req.GetWorkEmail() == "" {
		return fmt.Errorf("work_email: %w", domain.ErrEmptyField)
	}

	if !emailRegex.MatchString(req.GetWorkEmail()) {
		return fmt.Errorf("work_email: %w", domain.ErrInvalidEmail)
	}

	// Валидация пароля
	if req.GetPassword() == "" {
		return fmt.Errorf("password: %w", domain.ErrEmptyField)
	}

	// Проверка минимальной длины пароля и других требований
	if err := validatePassword(req.GetPassword()); err != nil {
		return fmt.Errorf("password: %w", err)
	}

	return nil
}

// ValidateOAuthYandexRequest валидирует запрос на OAuth через Яндекс
func ValidateOAuthYandexRequest(req *authv1.OAuthYandexRequest) error {
	if req == nil {
		return domain.ErrNilRequest
	}

	// Валидация кода авторизации
	if req.GetCode() == "" {
		return fmt.Errorf("code: %w", domain.ErrEmptyField)
	}

	// Валидация redirect_uri
	if req.GetRedirectUri() == "" {
		return fmt.Errorf("redirect_uri: %w", domain.ErrEmptyField)
	}

	// Проверка валидности URL
	_, err := url.ParseRequestURI(req.GetRedirectUri())
	if err != nil {
		return fmt.Errorf("redirect_uri: %w", domain.ErrInvalidURI)
	}

	return nil
}

// ValidateValidateRequest валидирует запрос на валидацию токена
func ValidateValidateRequest(req *authv1.ValidateRequest) error {
	if req == nil {
		return domain.ErrNilRequest
	}

	// Валидация токена
	if req.GetToken() == "" {
		return fmt.Errorf("token: %w", domain.ErrEmptyField)
	}

	// Проверяем минимальную длину и формат JWT (три секции, разделенные точками)
	if !isValidJWTFormat(req.GetToken()) {
		return fmt.Errorf("token: %w", domain.ErrInvalidToken)
	}

	return nil
}

// ValidateRefreshRequest валидирует запрос на обновление токена
func ValidateRefreshRequest(req *authv1.RefreshRequest) error {
	if req == nil {
		return domain.ErrNilRequest
	}

	// Валидация refresh токена
	if req.GetRefreshToken() == "" {
		return fmt.Errorf("refresh_token: %w", domain.ErrEmptyField)
	}

	// Проверяем минимальную длину и формат JWT (три секции, разделенные точками)
	if !isValidJWTFormat(req.GetRefreshToken()) {
		return fmt.Errorf("refresh_token: %w", domain.ErrInvalidToken)
	}

	return nil
}

// Вспомогательные функции

// validatePassword проверяет пароль на соответствие требованиям безопасности
func validatePassword(password string) error {
	// Минимальная длина

	if len(password) < 8 {
		return domain.ErrPasswordLen
	}

	// Проверка наличия букв верхнего регистра
	hasDigit := false
	hasUpper := false
	hasLower := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		}

		if char >= 'a' && char <= 'z' {
			hasLower = true
		}

		if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	if !hasUpper {
		return domain.ErrPasswordUppercase
	}

	if !hasDigit {
		return domain.ErrPasswordDigit
	}

	if !hasLower {
		return domain.ErrPasswordLowerCase
	}

	return nil
}

// isValidJWTFormat проверяет базовый формат JWT
func isValidJWTFormat(token string) bool {
	parts := strings.Split(token, ".")

	return len(parts) == 3 && len(token) >= 30 // Минимальная разумная длина JWT
}
