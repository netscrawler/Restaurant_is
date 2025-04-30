package pg

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	dto "github.com/netscrawler/Restaurant_is/order_service/internal/models/repository"
	"github.com/netscrawler/Restaurant_is/order_service/internal/storage/postgres"
)

type PgEvent struct {
	storage *postgres.Storage
	builder sq.StatementBuilderType
}

func NewPgEvent(db *postgres.Storage) *PgEvent {
	return &PgEvent{
		storage: db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (e *PgEvent) Save(ctx context.Context, event *dto.Event) error {
	query := e.builder.Insert("events").
		Columns("id", "type", "payload", "published").
		Values(event.ID, event.Type, event.Payload, false)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = e.storage.DB.Exec(ctx, sql, args...)
	return err
}

func (e *PgEvent) GetUnpublishedEvents(ctx context.Context) ([]dto.Event, error) {
	query := e.builder.Select("id", "type", "payload").
		From("events").
		Where(sq.Eq{"published": false})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := e.storage.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []dto.Event
	for rows.Next() {
		var event dto.Event
		if err := rows.Scan(&event.ID, &event.Type, &event.Payload); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (e *PgEvent) MarkAsPublished(ctx context.Context, eventID string) error {
	query := e.builder.Update("events").
		Set("published", true).
		Where(sq.Eq{"id": eventID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = e.storage.DB.Exec(ctx, sql, args...)
	return err
}
