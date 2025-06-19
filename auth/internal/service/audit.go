package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	"github.com/netscrawler/Restaurant_is/auth/internal/utils"
)

type AuditService struct {
	auditRepo repository.AuditRepository
	log       *slog.Logger
}

func NewAuditService(auditRepo repository.AuditRepository, log *slog.Logger) *AuditService {
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

	log := utils.LoggerWithTrace(ctx, s.log)

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
		log.Info(op+"failed log event", slog.Any("error", err))
	}

	return nil
}
