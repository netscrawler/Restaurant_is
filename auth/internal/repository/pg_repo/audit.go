package pgrepo

import (
	"context"

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
			event.IP,
			event.UserAgent,
			event.CreatedAt,
		).
		ToSql()
	if err != nil {
		// TODO: Add wrap.

		p.log.Error(op+" error build sql", zap.Error(err))

		return err
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		p.log.Error(op+" error insert log event", zap.Error(err))

		return err
	}

	return nil
}

func (p *pgAudit) GetAuthEvents(
	ctx context.Context,
	filter models.AuthFilter,
) ([]*models.AuthEvent, error) {
	const op = "repository.pg.Audit.GetAuthEvents"

	panic("not implemented") // TODO: Implement.
}
