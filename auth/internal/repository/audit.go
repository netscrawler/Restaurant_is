package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

type AuditRepository interface {
	LogAuthEvent(ctx context.Context, event *models.AuthEvent) error
}

type Audit struct {
	AuditRepository
}

func NewAudit(repo AuditRepository) *Audit {
	return &Audit{
		AuditRepository: repo,
	}
}
