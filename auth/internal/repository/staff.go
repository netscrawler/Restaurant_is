package repository

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
)

type StaffRepository interface {
	CreateStaff(ctx context.Context, staff *models.Staff) error
	GetStaffByEmail(ctx context.Context, workEmail string) (*models.Staff, error)
	UpdateStaffPassword(ctx context.Context, workEmail, newHash string) error
	DeactivateStaff(ctx context.Context, workEmail string) error
	UpdateStaff(ctx context.Context, staff *models.Staff) error
}

type Staff struct {
	StaffRepository
}

func NewStaff(repo StaffRepository) *Staff {
	return &Staff{
		StaffRepository: repo,
	}
}
