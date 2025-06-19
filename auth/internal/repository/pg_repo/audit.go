package pgrepo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/infra/out/postgres"
	"github.com/netscrawler/Restaurant_is/auth/internal/utils"
)

type pgAudit struct {
	log *slog.Logger
	db  *postgres.Storage
}

func NewPgAudit(db *postgres.Storage, log *slog.Logger) *pgAudit {
	return &pgAudit{
		log: log,
		db:  db,
	}
}

func (p *pgAudit) LogAuthEvent(ctx context.Context, event *models.AuthEvent) error {
	const op = "repository.pg.Audit.LogAuthEvent"

	log := utils.LoggerWithTrace(ctx, p.log)

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
		log.Error(op+" error build sql", slog.Any("error", err))

		return fmt.Errorf("%w (%w)", domain.ErrBuildQuery, err)
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		log.Error(op+" error insert log event", slog.Any("error", err))

		return fmt.Errorf("%w (%w)", domain.ErrExecQuery, err)
	}

	return nil
}
