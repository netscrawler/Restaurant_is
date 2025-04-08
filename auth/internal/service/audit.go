package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
)

// AuditService предоставляет функции для аудита событий аутентификации
type AuditService struct {
	auditRepo repository.AuditRepository
}

// NewAuditService создает новый экземпляр сервиса аудита
func NewAuditService(auditRepo repository.AuditRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
	}
}

// LogLoginEvent записывает событие успешного входа
func (s *AuditService) LogLoginEvent(
	ctx context.Context,
	userID string,
	userType string,
	ipAddress, userAgent string,
) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	event := &models.AuthEvent{
		ID:        uuid.New(),
		UserID:    userUUID,
		UserType:  models.UserType(userType),
		Action:    models.ActionLogin,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	return s.auditRepo.LogAuthEvent(ctx, event)
}

// LogLogoutEvent записывает событие выхода
func (s *AuditService) LogLogoutEvent(
	ctx context.Context,
	userID string,
	userType string,
	ipAddress, userAgent string,
) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	event := &models.AuthEvent{
		ID:        uuid.New(),
		UserID:    userUUID,
		UserType:  models.UserType(userType),
		Action:    models.ActionLogout,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	return s.auditRepo.LogAuthEvent(ctx, event)
}

// LogTokenRefreshEvent записывает событие обновления токена
func (s *AuditService) LogTokenRefreshEvent(
	ctx context.Context,
	userID string,
	userType string,
	ipAddress, userAgent string,
) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	event := &models.AuthEvent{
		ID:        uuid.New(),
		UserID:    userUUID,
		UserType:  models.UserType(userType),
		Action:    models.ActionTokenRefresh,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	return s.auditRepo.LogAuthEvent(ctx, event)
}

// LogTokenRevokeEvent записывает событие отзыва токена
func (s *AuditService) LogTokenRevokeEvent(
	ctx context.Context,
	userID string,
	userType string,
	ipAddress, userAgent string,
) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	event := &models.AuthEvent{
		ID:        uuid.New(),
		UserID:    userUUID,
		UserType:  models.UserType(userType),
		Action:    models.ActionTokenRevoke,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	return s.auditRepo.LogAuthEvent(ctx, event)
}

// LogAuthEvent записывает произвольное событие аутентификации
func (s *AuditService) LogAuthEvent(
	ctx context.Context,
	userID string,
	userType string,
	action string,
	ipAddress, userAgent string,
) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	event := &models.AuthEvent{
		ID:        uuid.New(),
		UserID:    userUUID,
		UserType:  models.UserType(userType),
		Action:    models.AuthEventAction(action),
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	return s.auditRepo.LogAuthEvent(ctx, event)
}

// GetUserAuthEvents получает события аутентификации для конкретного пользователя
func (s *AuditService) GetUserAuthEvents(
	ctx context.Context,
	userID string,
	userType string,
	limit, offset int,
) ([]*models.AuthEvent, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	userTypeParsed := models.UserType(userType)

	filter := models.AuthFilter{
		UserID:   &userUUID,
		UserType: &userTypeParsed,
		Limit:    limit,
		Offset:   offset,
	}

	return s.auditRepo.GetAuthEvents(ctx, filter)
}

// GetAuthEvents получает события аутентификации по фильтру
func (s *AuditService) GetAuthEvents(
	ctx context.Context,
	filter models.AuthFilter,
) ([]*models.AuthEvent, error) {
	return s.auditRepo.GetAuthEvents(ctx, filter)
}

// GetRecentLoginEvents получает недавние события входа для указанного IP-адреса
// Полезно для обнаружения множественных попыток входа
func (s *AuditService) GetRecentLoginEvents(
	ctx context.Context,
	ipAddress string,
	minutes int,
) ([]*models.AuthEvent, error) {
	action := models.ActionLogin
	from := time.Now().Add(-time.Duration(minutes) * time.Minute)

	filter := models.AuthFilter{
		Action:    &action,
		IPAddress: &ipAddress,
		DateFrom:  &from,
		Limit:     100, // Ограничиваем количество возвращаемых записей
	}

	return s.auditRepo.GetAuthEvents(ctx, filter)
}

// GetUserSessionHistory получает историю сессий пользователя
func (s *AuditService) GetUserSessionHistory(
	ctx context.Context,
	userID string,
	userType string,
	days int,
) ([]*models.AuthEvent, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	userTypeParsed := models.UserType(userType)
	from := time.Now().AddDate(0, 0, -days)

	filter := models.AuthFilter{
		UserID:   &userUUID,
		UserType: &userTypeParsed,
		DateFrom: &from,
		Limit:    1000, // Ограничиваем количество возвращаемых записей
	}

	return s.auditRepo.GetAuthEvents(ctx, filter)
}
