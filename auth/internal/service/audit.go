package service

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
)

// Вспомогательные сервисы.
type AuditService struct {
	auditRepo repository.AuditRepository
}

func NewAuditService(auditRepo repository.AuditRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
	}
}

func (s *AuditService) LogAuthEvent(ctx context.Context, event *models.AuthEvent) error {
	// Реализация
	panic("implement me")
}
