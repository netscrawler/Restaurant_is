package pgrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"user_service/internal/domain"
	"user_service/internal/domain/models"
	"user_service/internal/storage/postgres"
)

type userRepository struct {
	db *postgres.Storage
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *postgres.Storage) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := squirrel.Insert("users").
		Columns(
			"id",
			"email",
			"phone",
			"full_name",
			"is_active",
			"created_at",
			"updated_at",
		).
		Values(
			user.ID,
			user.Email,
			user.Phone,
			user.FullName,
			user.IsActive,
			user.CreatedAt,
			user.UpdatedAt,
		).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create user query: %w", err)
	}

	_, err = r.db.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := squirrel.Select(
		"id",
		"email",
		"phone",
		"full_name",
		"is_active",
		"created_at",
		"updated_at",
	).
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get user by ID query: %w", err)
	}

	var user models.User

	err = r.db.DB.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.FullName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := squirrel.Select(
		"id",
		"email",
		"phone",
		"full_name",
		"is_active",
		"created_at",
		"updated_at",
	).
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get user by email query: %w", err)
	}

	var user models.User

	err = r.db.DB.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.FullName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := squirrel.Select(
		"id",
		"email",
		"phone",
		"full_name",
		"is_active",
		"created_at",
		"updated_at",
	).
		From("users").
		Where(squirrel.Eq{"phone": phone}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get user by phone query: %w", err)
	}

	var user models.User

	err = r.db.DB.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.FullName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := squirrel.Update("users").
		Set("email", user.Email).
		Set("phone", user.Phone).
		Set("full_name", user.FullName).
		Set("is_active", user.IsActive).
		Set("updated_at", user.UpdatedAt).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update user query: %w", err)
	}

	result, err := r.db.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := squirrel.Delete("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete user query: %w", err)
	}

	result, err := r.db.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) List(
	ctx context.Context,
	onlyActive bool,
	offset, limit int,
) ([]*models.User, error) {
	query := squirrel.Select(
		"id",
		"email",
		"phone",
		"full_name",
		"is_active",
		"created_at",
		"updated_at",
	).
		From("users").
		OrderBy("created_at DESC").
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar)

	if onlyActive {
		query = query.Where(squirrel.Eq{"is_active": true})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list users query: %w", err)
	}

	rows, err := r.db.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Phone,
			&user.FullName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *userRepository) Count(ctx context.Context, onlyActive bool) (int, error) {
	query := squirrel.Select("COUNT(*)").
		From("users").
		PlaceholderFormat(squirrel.Dollar)

	if onlyActive {
		query = query.Where(squirrel.Eq{"is_active": true})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count users query: %w", err)
	}

	var count int

	err = r.db.DB.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
