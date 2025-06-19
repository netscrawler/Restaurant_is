package pgrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"user_service/internal/domain"
	"user_service/internal/domain/models"
	"user_service/internal/domain/repository"
	"user_service/internal/storage/postgres"
)

type roleRepository struct {
	db *postgres.Storage
}

// NewRoleRepository создает новый экземпляр RoleRepository.
func NewRoleRepository(db *postgres.Storage) *roleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *models.Role) error {
	query := squirrel.Insert("roles").
		Columns(
			"id",
			"name",
			"description",
			"created_at",
			"updated_at",
		).
		Values(
			role.ID,
			role.Name,
			role.Description,
			role.CreatedAt,
			role.UpdatedAt,
		).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create role query: %w", err)
	}

	_, err = r.db.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	query := squirrel.Select(
		"id",
		"name",
		"description",
		"created_at",
		"updated_at",
	).
		From("roles").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get role by ID query: %w", err)
	}

	var role models.Role

	err = r.db.DB.QueryRow(ctx, sql, args...).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoleNotFound
		}

		return nil, fmt.Errorf("failed to get role by ID: %w", err)
	}

	return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	query := squirrel.Select(
		"id",
		"name",
		"description",
		"created_at",
		"updated_at",
	).
		From("roles").
		Where(squirrel.Eq{"name": name}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get role by name query: %w", err)
	}

	var role models.Role

	err = r.db.DB.QueryRow(ctx, sql, args...).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoleNotFound
		}

		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}

	return &role, nil
}

func (r *roleRepository) Update(ctx context.Context, role *models.Role) error {
	query := squirrel.Update("roles").
		Set("name", role.Name).
		Set("description", role.Description).
		Set("updated_at", role.UpdatedAt).
		Where(squirrel.Eq{"id": role.ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update role query: %w", err)
	}

	result, err := r.db.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRoleNotFound
	}

	return nil
}

func (r *roleRepository) List(ctx context.Context) ([]*models.Role, error) {
	query := squirrel.Select(
		"id",
		"name",
		"description",
		"created_at",
		"updated_at",
	).
		From("roles").
		OrderBy("name ASC").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list roles query: %w", err)
	}

	rows, err := r.db.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []*models.Role

	for rows.Next() {
		var role models.Role

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}

		roles = append(roles, &role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating roles: %w", err)
	}

	return roles, nil
}

type userRoleRepository struct {
	db *pgxpool.Pool
}

// NewUserRoleRepository создает новый экземпляр UserRoleRepository.
func NewUserRoleRepository(db *pgxpool.Pool) repository.UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) AssignRole(ctx context.Context, userRole *models.UserRole) error {
	query := squirrel.Insert("user_roles").
		Columns("id", "user_id", "role_id", "assigned_at").
		Values(userRole.ID, userRole.UserID, userRole.RoleID, userRole.AssignedAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build assign role query: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

func (r *userRoleRepository) RevokeRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := squirrel.Delete("user_roles").
		Where(squirrel.Eq{"user_id": userID, "role_id": roleID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build revoke role query: %w", err)
	}

	result, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserRoleNotFound
	}

	return nil
}

func (r *userRoleRepository) GetUserRoles(
	ctx context.Context,
	userID uuid.UUID,
) ([]*models.Role, error) {
	query := squirrel.Select(
		"r.id",
		"r.name",
		"r.description",
		"r.created_at",
		"r.updated_at",
	).
		From("user_roles ur").
		Join("roles r ON ur.role_id = r.id").
		Where(squirrel.Eq{"ur.user_id": userID}).
		OrderBy("r.name ASC").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get user roles query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []*models.Role

	for rows.Next() {
		var role models.Role

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}

		roles = append(roles, &role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user roles: %w", err)
	}

	return roles, nil
}

func (r *userRoleRepository) HasRole(ctx context.Context, userID, roleID uuid.UUID) (bool, error) {
	query := squirrel.Select("COUNT(*)").
		From("user_roles").
		Where(squirrel.Eq{"user_id": userID, "role_id": roleID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build has role query: %w", err)
	}

	var count int

	err = r.db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user has role: %w", err)
	}

	return count > 0, nil
}
