package pg

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/order_service/internal/infra/in/postgres"
	dto "github.com/netscrawler/Restaurant_is/order_service/internal/models/repository"
)

type PgOrder struct {
	storage *postgres.Storage
	builder sq.StatementBuilderType
}

func NewPgOrder(db *postgres.Storage) *PgOrder {
	return &PgOrder{
		storage: db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (o *PgOrder) Save(ctx context.Context, order *dto.Order) (uint64, error) {
	query := o.builder.Insert("orders").
		Columns(
			"id",
			"user_id",
			"price",
			"created_at",
			"updated_at",
			"status",
			"order_type",
			"address",
			"dish_quantites",
		).
		Values(
			order.ID,
			order.UserID,
			order.Price,
			order.CreatedAt,
			order.UpdatedAt,
			order.Status,
			order.OrderType,
			order.Address,
			order.DishQuantites,
		).
		Suffix("RETURNING num")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := o.storage.DB.Query(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var num uint64
	if rows.Next() {
		if err := rows.Scan(&num); err != nil {
			return 0, err
		}
	}

	return num, rows.Err()
}

func (o *PgOrder) Get(ctx context.Context, id uuid.UUID) (*dto.Order, error) {
	uid := id.String()
	query := o.builder.Select(
		"id",
		"user_id",
		"num",
		"price",
		"created_at",
		"updated_at",
		"status",
		"order_type",
		"address",
		"dish_quantites",
	).
		From("orders").
		Where(sq.Eq{"id": uid})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := o.storage.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order dto.Order
	if rows.Next() {
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Num,
			&order.Price,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.Status,
			&order.OrderType,
			&order.Address,
			&order.DishQuantites,
		); err != nil {
			return nil, err
		}
	}

	return &order, rows.Err()
}

func (o *PgOrder) Update(ctx context.Context, order *dto.Order) error {
	return nil
}
