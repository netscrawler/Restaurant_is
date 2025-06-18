package pgrepo

import (
	"context"
	"fmt"

	"user_service/internal/domain"
	"user_service/internal/domain/models"
	"user_service/internal/storage/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type staffRepository struct {
	db *postgres.Storage
}

// NewStaffRepository создает новый экземпляр StaffRepository
func NewStaffRepository(db *postgres.Storage) *staffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) Create(ctx context.Context, staff *models.Staff) error {
	query := squirrel.Insert("staff").
		Columns(
			"id",
			"work_email",
			"work_phone",
			"full_name",
			"position",
			"is_active",
			"hire_date",
		).
		Values(
			staff.ID,
			staff.WorkEmail,
			staff.WorkPhone,
			staff.FullName,
			staff.Position,
			staff.IsActive,
			staff.HireDate,
		).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create staff query: %w", err)
	}

	_, err = r.db.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to create staff: %w", err)
	}

	return nil
}

func (r *staffRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Staff, error) {
	query := squirrel.Select(
		"id",
		"work_email",
		"work_phone",
		"full_name",
		"position",
		"is_active",
		"hire_date",
	).
		From("staff").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get staff by ID query: %w", err)
	}

	var staff models.Staff
	err = r.db.DB.QueryRow(ctx, sql, args...).Scan(
		&staff.ID,
		&staff.WorkEmail,
		&staff.WorkPhone,
		&staff.FullName,
		&staff.Position,
		&staff.IsActive,
		&staff.HireDate,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrStaffNotFound
		}
		return nil, fmt.Errorf("failed to get staff by ID: %w", err)
	}

	return &staff, nil
}

func (r *staffRepository) GetByWorkEmail(
	ctx context.Context,
	workEmail string,
) (*models.Staff, error) {
	query := squirrel.Select(
		"id",
		"work_email",
		"work_phone",
		"full_name",
		"position",
		"is_active",
		"hire_date",
	).
		From("staff").
		Where(squirrel.Eq{"work_email": workEmail}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get staff by work email query: %w", err)
	}

	var staff models.Staff
	err = r.db.DB.QueryRow(ctx, sql, args...).Scan(
		&staff.ID,
		&staff.WorkEmail,
		&staff.WorkPhone,
		&staff.FullName,
		&staff.Position,
		&staff.IsActive,
		&staff.HireDate,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrStaffNotFound
		}
		return nil, fmt.Errorf("failed to get staff by work email: %w", err)
	}

	return &staff, nil
}

func (r *staffRepository) Update(ctx context.Context, staff *models.Staff) error {
	query := squirrel.Update("staff").
		Set("work_email", staff.WorkEmail).
		Set("work_phone", staff.WorkPhone).
		Set("full_name", staff.FullName).
		Set("position", staff.Position).
		Set("is_active", staff.IsActive).
		Where(squirrel.Eq{"id": staff.ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update staff query: %w", err)
	}

	result, err := r.db.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update staff: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrStaffNotFound
	}

	return nil
}

func (r *staffRepository) List(
	ctx context.Context,
	onlyActive bool,
	offset, limit int,
) ([]*models.Staff, error) {
	query := squirrel.Select(
		"id",
		"work_email",
		"work_phone",
		"full_name",
		"position",
		"is_active",
		"hire_date",
	).
		From("staff").
		OrderBy("hire_date DESC").
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar)

	if onlyActive {
		query = query.Where(squirrel.Eq{"is_active": true})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list staff query: %w", err)
	}

	rows, err := r.db.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list staff: %w", err)
	}
	defer rows.Close()

	var staffList []*models.Staff
	for rows.Next() {
		var staff models.Staff
		err := rows.Scan(
			&staff.ID,
			&staff.WorkEmail,
			&staff.WorkPhone,
			&staff.FullName,
			&staff.Position,
			&staff.IsActive,
			&staff.HireDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan staff: %w", err)
		}
		staffList = append(staffList, &staff)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating staff: %w", err)
	}

	return staffList, nil
}

func (r *staffRepository) Count(ctx context.Context, onlyActive bool) (int, error) {
	query := squirrel.Select("COUNT(*)").
		From("staff").
		PlaceholderFormat(squirrel.Dollar)

	if onlyActive {
		query = query.Where(squirrel.Eq{"is_active": true})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count staff query: %w", err)
	}

	var count int
	err = r.db.DB.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count staff: %w", err)
	}

	return count, nil
}
