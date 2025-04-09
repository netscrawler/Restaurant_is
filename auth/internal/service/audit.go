package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	"go.uber.org/zap"
)

type AuditService struct {
	auditRepo repository.AuditRepository
	log       *zap.Logger
}

func NewAuditService(auditRepo repository.AuditRepository, log *zap.Logger) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
		log:       log,
	}
}

func (s *AuditService) LogEvent(
	ctx context.Context,
	userID string,
	userType string,
	eventAction models.AuthEventAction,
	ipAddress, userAgent string,
) error {
	const op = "service.audit.LogEvent"

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("%w (%w)", domain.ErrInvalidUserUUID, err)
	}

	event := models.NewAuthEvent(
		userUUID,
		models.UserType(userType),
		eventAction,
		ipAddress,
		userAgent,
	)

	err = s.auditRepo.LogAuthEvent(ctx, event)
	if err != nil {
		s.log.Info(op+"failed log event", zap.String("error", err.Error()))
	}

	return nil
}
