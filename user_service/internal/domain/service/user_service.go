package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"user_service/internal/domain/models"
	"user_service/internal/domain/repository"
)

// UserService представляет доменный сервис для работы с пользователями.
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser создает нового пользователя.
func (s *UserService) CreateUser(
	ctx context.Context,
	email, phone, fullName string,
) (*models.User, error) {
	// Проверяем, не существует ли уже пользователь с таким email
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Проверяем, не существует ли уже пользователь с таким телефоном
	existingUserByPhone, err := s.userRepo.GetByPhone(ctx, phone)
	if err == nil && existingUserByPhone != nil {
		return nil, fmt.Errorf("user with phone %s already exists", phone)
	}

	user := models.NewUser(email, phone, fullName)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser получает пользователя по ID.
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByEmail получает пользователя по email.
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByPhone получает пользователя по телефону.
func (s *UserService) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUser обновляет данные пользователя.
func (s *UserService) UpdateUser(
	ctx context.Context,
	id uuid.UUID,
	email, phone, fullName string,
	isActive *bool,
) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Update(email, phone, fullName)

	if isActive != nil {
		if *isActive {
			user.Activate()
		} else {
			user.Deactivate()
		}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser удаляет пользователя.
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return errors.New("user not found")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers возвращает список пользователей.
func (s *UserService) ListUsers(
	ctx context.Context,
	onlyActive bool,
	offset, limit int,
) ([]*models.User, int, error) {
	users, err := s.userRepo.List(ctx, onlyActive, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := s.userRepo.Count(ctx, onlyActive)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}
