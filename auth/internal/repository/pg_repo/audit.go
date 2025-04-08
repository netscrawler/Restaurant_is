package pgrepo

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
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

	// Начинаем строить запрос на выборку
	builder := p.db.Builder.
		Select(
			"id",
			"user_id",
			"user_type",
			"action",
			"ip_address",
			"user_agent",
			"created_at",
		).
		From("auth_logs")

	// Добавляем условия фильтрации
	if filter.UserID != nil {
		builder = builder.Where(squirrel.Eq{"user_id": filter.UserID})
	}

	if filter.UserType != nil {
		builder = builder.Where(squirrel.Eq{"user_type": *filter.UserType})
	}

	if filter.Action != nil {
		builder = builder.Where(squirrel.Eq{"action": *filter.Action})
	}

	if filter.IPAddress != nil {
		builder = builder.Where(squirrel.Eq{"ip_address": *filter.IPAddress})
	}

	if filter.DateFrom != nil {
		builder = builder.Where("created_at >= ?", *filter.DateFrom)
	}

	if filter.DateTo != nil {
		builder = builder.Where("created_at <= ?", *filter.DateTo)
	}

	// Добавляем сортировку по дате создания (от новых к старым)
	builder = builder.OrderBy("created_at DESC")

	// Добавляем ограничение по количеству и смещение
	if filter.Limit > 0 {
		builder = builder.Limit(uint64(filter.Limit))
	}

	if filter.Offset > 0 {
		builder = builder.Offset(uint64(filter.Offset))
	}

	// Формируем итоговый запрос
	query, args, err := builder.ToSql()
	if err != nil {
		p.log.Error(op+" error build sql", zap.Error(err))

		return nil, err
	}

	// Выполняем запрос
	rows, err := p.db.DB.Query(ctx, query, args...)
	if err != nil {
		p.log.Error(op+" error query auth events", zap.Error(err),
			zap.String("query", query), zap.Any("args", args))
		return nil, err
	}
	defer rows.Close()

	// Обрабатываем результаты
	var events []*models.AuthEvent
	for rows.Next() {
		var event models.AuthEvent
		err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.UserType,
			&event.Action,
			&event.IPAddress,
			&event.UserAgent,
			&event.CreatedAt,
		)
		if err != nil {
			p.log.Error(op+" error scan auth event", zap.Error(err))

			return nil, err
		}

		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		p.log.Error(op+" error iterating rows", zap.Error(err))
		return nil, err
	}

	return events, nil
}

// GetLoginAttempts получает количество попыток входа для IP-адреса за указанный период
func (p *pgAudit) GetLoginAttempts(
	ctx context.Context,
	ipAddress string,
	minutes int,
) (int, error) {
	const op = "repository.pg.Audit.GetLoginAttempts"

	// Строим запрос для подсчета попыток входа
	query, args, err := p.db.Builder.
		Select("COUNT(*)").
		From("auth_logs").
		Where("ip_address = ?", ipAddress).
		Where("action = ?", models.ActionLogin).
		Where("created_at >= ?", time.Now().Add(-time.Duration(minutes)*time.Minute)).
		ToSql()
	if err != nil {
		p.log.Error(op+" error build sql", zap.Error(err))

		return 0, err
	}

	// Выполняем запрос
	var count int
	err = p.db.DB.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		p.log.Error(op+" error query login attempts", zap.Error(err),
			zap.String("query", query), zap.Any("args", args))

		return 0, err
	}

	return count, nil
}
