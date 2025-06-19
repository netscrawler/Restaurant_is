package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/domain"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/infra/out/postgres"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/models/dto"
)

type dishPgRepo struct {
	storage *postgres.Storage
}

const (
	dishTable       = "dishes"
	defaultPageSize = 10
)

func (d *dishPgRepo) Create(ctx context.Context, dish *dto.Dish) error {
	sql, args, err := d.storage.Builder.Insert(dishTable).
		Columns(
			"id",
			"name",
			"description",
			"price",
			"category_id",
			"cooking_time_min",
			"image_url",
			"is_available",
			"calories",
			"created_at",
			"updated_at",
		).
		Values(
			dish.ID,
			dish.Name,
			dish.Description,
			dish.Price,
			dish.CategoryID,
			dish.CookingTimeMin,
			dish.ImageURL,
			dish.IsAvailable,
			dish.Calories,
			dish.CreatedAt,
			dish.UpdatedAt,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	_, err = d.storage.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return nil
}

func (d *dishPgRepo) GetByID(ctx context.Context, id uuid.UUID) (*dto.Dish, error) {
	sql, args, err := d.storage.Builder.Select(
		"id",
		"name",
		"description",
		"price",
		"category_id",
		"cooking_time_min",
		"image_url",
		"is_available",
		"calories",
		"created_at",
		"updated_at",
	).
		From(dishTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	var dish dto.Dish

	err = d.storage.DB.QueryRow(ctx, sql, args...).Scan(
		&dish.ID,
		&dish.Name,
		&dish.Description,
		&dish.Price,
		&dish.CategoryID,
		&dish.CookingTimeMin,
		&dish.ImageURL,
		&dish.IsAvailable,
		&dish.Calories,
		&dish.CreatedAt,
		&dish.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: dish not found", domain.ErrInvalid)
		}

		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return &dish, nil
}

func (d *dishPgRepo) GetByFilter(
	ctx context.Context,
	filter *dto.ListDishFilter,
) ([]dto.Dish, error) {
	query := d.storage.Builder.Select(
		"id",
		"name",
		"description",
		"price",
		"category_id",
		"cooking_time_min",
		"image_url",
		"is_available",
		"calories",
		"created_at",
		"updated_at",
	).
		From(dishTable).
		Where(squirrel.Eq{"is_available": filter.OnlyAvailable})

	if filter.CategoryID != nil {
		query = query.Where(squirrel.Eq{"category_id": *filter.CategoryID})
	}

	if filter.OnlyAvailable {
		query = query.Where(squirrel.Eq{"is_available": true})
	}

	offset := (filter.Page - 1) * filter.PageSize
	query = query.Limit(uint64(filter.PageSize))
	query = query.Offset(uint64(offset))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	rows, err := d.storage.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}
	defer rows.Close()

	var dishes []dto.Dish

	for rows.Next() {
		var dish dto.Dish

		err = rows.Scan(
			&dish.ID,
			&dish.Name,
			&dish.Description,
			&dish.Price,
			&dish.CategoryID,
			&dish.CookingTimeMin,
			&dish.ImageURL,
			&dish.IsAvailable,
			&dish.Calories,
			&dish.CreatedAt,
			&dish.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrInternal, err)
		}

		dishes = append(dishes, dish)
	}

	return dishes, nil
}

func (d *dishPgRepo) Update(ctx context.Context, dish *dto.Dish) error {
	updateMap := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if dish.Name != "" {
		updateMap["name"] = dish.Name
	}

	if dish.Description != "" {
		updateMap["description"] = dish.Description
	}

	if dish.Price > 0 {
		updateMap["price"] = dish.Price
	}

	if dish.CategoryID != 0 {
		updateMap["category_id"] = dish.CategoryID
	}

	if dish.CookingTimeMin != 0 {
		updateMap["cooking_time_min"] = dish.CookingTimeMin
	}

	if dish.ImageURL != "" {
		updateMap["image_url"] = dish.ImageURL
	}

	if dish.Calories != 0 {
		updateMap["calories"] = dish.Calories
	}

	sql, args, err := d.storage.Builder.Update(dishTable).
		SetMap(updateMap).
		Where(squirrel.Eq{"id": dish.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	_, err = d.storage.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return nil
}

func (d *dishPgRepo) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := d.storage.Builder.Delete(dishTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	_, err = d.storage.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return nil
}

func NewDishPgRepo(storage *postgres.Storage) *dishPgRepo {
	return &dishPgRepo{
		storage: storage,
	}
}
