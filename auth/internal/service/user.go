package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
)

type UserService struct {
	clientRepo repository.ClientRepository
	staffRepo  repository.StaffRepository
	log        *slog.Logger
	ntf        NotifySender
}

func NewUserService(
	client repository.ClientRepository,
	staff repository.StaffRepository,
	notify NotifySender,
	log *slog.Logger,
) *UserService {
	return &UserService{
		clientRepo: client,
		staffRepo:  staff,
		log:        log,
		ntf:        notify,
	}
}

func (u *UserService) RegisterStaff(
	ctx context.Context,
	email, position string,
) (*models.Staff, error) {
	exist, err := u.staffExist(ctx, email)
	if err != nil {
		return nil, domain.ErrInternal
	}

	if exist {
		return nil, domain.ErrUserAlreadyExist
	}

	staff, pswd, err := models.NewStaff(email, position)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = u.staffRepo.CreateStaff(ctx, staff)
	if err != nil {
		return nil, domain.ErrInternal
	}

	msg := fmt.Sprintf("Your registered in system,\nLogin:%s\nPassword:%s", email, pswd)

	u.ntf.Send(ctx, email, msg)

	return staff, nil
}

func (u *UserService) DeactivateStaff(ctx context.Context, email string) (*models.Staff, error) {
	staff, err := u.staffRepo.GetStaffByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrInternal
	}

	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrNotFound
	}

	err = u.staffRepo.DeactivateStaff(ctx, email)
	if err != nil {
		return nil, domain.ErrInternal
	}

	staff.IsActive = false

	return staff, nil
}

func (u *UserService) UpdateStaff(
	ctx context.Context,
	email string,
	newEmail, pswd, position *string,
) (*models.Staff, error) {
	staff, err := u.staffRepo.GetStaffByEmail(ctx, email)

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrInternal
	}

	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrNotFound
	}

	if newEmail != nil {
		staff.UpdateEmail(*newEmail)
	}

	if pswd != nil {
		staff.UpdatePassword(*pswd)
	}

	if position != nil {
		staff.UpdatePosition(*position)
	}

	err = u.staffRepo.UpdateStaff(ctx, staff)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return staff, nil
}

func (u *UserService) staffExist(ctx context.Context, email string) (bool, error) {
	_, err := u.staffRepo.GetStaffByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return false, fmt.Errorf("%w (%w)", domain.ErrInternal, err)
	}

	if errors.Is(err, domain.ErrNotFound) {
		return false, nil
	}

	return true, nil
}
