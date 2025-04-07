package pgrepo

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/storage/postgres"
	"go.uber.org/zap"
)

type pgStaff struct {
	log *zap.Logger
	db  *postgres.Storage
}

func NewPgStaff(db *postgres.Storage, log *zap.Logger) *pgStaff {
	return &pgStaff{
		log: log,
		db:  db,
	}
}

func (p *pgStaff) CreateStaff(ctx context.Context, staff *models.Staff) error {
	const op = "repository.pg.Staff.Create"

	query, args, err := p.db.Builder.
		Insert("staff").
		Columns("work_email", "password_hash").
		Values(staff.WorkEmail, staff.PasswordHash).
		ToSql()
	if err != nil {
		p.log.Error(op+"failed to build SQL query", zap.Error(err))

		return err
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		p.log.Error(op+"failed to execute SQL query", zap.Error(err))

		return err
	}

	return nil
}

func (p *pgStaff) GetStaffByEmail(ctx context.Context, workEmail string) (*models.Staff, error) {
	const op = "repository.pg.Staff.GetStaffByEmail"

	query, args, err := p.db.Builder.
		Select("id", "work_email", "password_hash", "is_active", "created_at", "updated_at").
		From("staff").
		Where(squirrel.Eq{"work_email": workEmail}).
		ToSql()
	if err != nil {
		p.log.Error(op+"failed to build SQL query", zap.Error(err))

		return nil, err
	}

	row := p.db.DB.QueryRow(ctx, query, args...)

	staff := models.Staff{}

	err = row.Scan(
		&staff.ID,
		&staff.WorkEmail,
		&staff.PasswordHash,
		&staff.IsActive,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Info(op+"staff not found", zap.String("work_email", workEmail))
		}

		p.log.Error(op+"failed to scan row", zap.Error(err))

		return nil, err
	}

	return &staff, nil
}

func (p *pgStaff) UpdateStaffPassword(ctx context.Context, workEmail string, newHash string) error {
	const op = "repository.pg.Staff.UpdateStaffPassword"

	query, args, err := p.db.Builder.Update("staff").
		Set("password_hash", newHash).
		Set("updated_at", "NOW()").
		Where(squirrel.Eq{"work_email": workEmail}).
		ToSql()
	if err != nil {
		p.log.Error(op+"failed to build SQL query", zap.Error(err))

		return err
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		p.log.Error(op+"failed to execute SQL query", zap.Error(err))

		return err
	}

	return nil
}

func (p *pgStaff) DeactivateStaff(ctx context.Context, workEmail string) error {
	const op = "repository.pg.Staff.DeactivateStaff"

	query, args, err := p.db.Builder.Update("staff").
		Set("is_active", false).
		Set("updated_at", "NOW()").
		Where(squirrel.Eq{"work_email": workEmail}).
		ToSql()
	if err != nil {
		p.log.Error(op+"failed to build SQL query", zap.Error(err))

		return err
	}

	_, err = p.db.DB.Exec(ctx, query, args...)
	if err != nil {
		p.log.Error(op+"failed to execute SQL query", zap.Error(err))

		return err
	}

	p.log.Info(op+"staff deactivated", zap.String("work_email", workEmail))

	return nil
}
