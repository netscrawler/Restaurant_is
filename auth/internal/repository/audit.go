package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

// AuditRepository - логирование событий
type AuditRepository interface {
	LogAuthEvent(ctx context.Context, event *models.AuthEvent) error
	GetAuthEvents(ctx context.Context, filter models.AuthFilter) ([]*models.AuthEvent, error)
}

type Audit struct {
	AuditRepository
}

func NewAudit(repo AuditRepository) *Audit {
	return &Audit{
		AuditRepository: repo,
	}
}
