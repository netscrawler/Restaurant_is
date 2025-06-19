package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"user_service/internal/domain/models"
	"user_service/internal/domain/repository"
)

// StaffService представляет доменный сервис для работы с сотрудниками.
type StaffService struct {
	staffRepo repository.StaffRepository
}

// NewStaffService создает новый экземпляр StaffService.
func NewStaffService(staffRepo repository.StaffRepository) *StaffService {
	return &StaffService{
		staffRepo: staffRepo,
	}
}

// CreateStaff создает нового сотрудника.
func (s *StaffService) CreateStaff(
	ctx context.Context,
	workEmail, workPhone, fullName, position string,
) (*models.Staff, error) {
	// Проверяем, не существует ли уже сотрудник с таким email
	existingStaff, err := s.staffRepo.GetByWorkEmail(ctx, workEmail)
	if err == nil && existingStaff != nil {
		return nil, fmt.Errorf("staff with email %s already exists", workEmail)
	}

	staff := models.NewStaff(workEmail, workPhone, fullName, position)

	if err := s.staffRepo.Create(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	return staff, nil
}

// UpdateStaff обновляет данные сотрудника.
func (s *StaffService) UpdateStaff(
	ctx context.Context,
	id uuid.UUID,
	workPhone, position string,
	isActive *bool,
) (*models.Staff, error) {
	staff, err := s.staffRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	if staff == nil {
		return nil, errors.New("staff not found")
	}

	staff.Update(workPhone, position)

	if isActive != nil {
		if *isActive {
			staff.Activate()
		} else {
			staff.Deactivate()
		}
	}

	if err := s.staffRepo.Update(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	return staff, nil
}

// ListStaff возвращает список сотрудников.
func (s *StaffService) ListStaff(
	ctx context.Context,
	onlyActive bool,
	offset, limit int,
) ([]*models.Staff, int, error) {
	staff, err := s.staffRepo.List(ctx, onlyActive, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list staff: %w", err)
	}

	total, err := s.staffRepo.Count(ctx, onlyActive)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count staff: %w", err)
	}

	return staff, total, nil
}
