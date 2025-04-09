package pgrepo

import (
	"context"
	"fmt"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgAudit struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgAudit(db *postgres.Storage, log *zap.Logger) *pgAudit {
	return &pgAudit{
		log: log,
		db:  db,
	}
}

func (p *pgAudit) LogAuthEvent(ctx context.Context, event *models.AuthEvent) error {
	const op = "repository.pg.Audit.LogAuthEvent"

	query, args, err := p.db.Builder.
		Insert("auth_logs").
		Columns(
			"user_id",
			"user_type",
			"action",
			"ip_address",
			"user_agent",
			"created_at",
		).
		Values(
			event.UserID,
			event.UserType,
			event.Action,
			event.IPAddress,
			event.UserAgent,
			event.CreatedAt,
		).
		ToSql()
	if err != nil {
		p.log.Error(op+" error build sql", zap.Error(err))

		return fmt.Errorf("%w (%w)", domain.ErrBuildQuery, err)
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		p.log.Error(op+" error insert log event", zap.Error(err))

		return fmt.Errorf("%w (%w)", domain.ErrExecQuery, err)
	}

	return nil
}
