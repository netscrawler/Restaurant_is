package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"user_service/internal/domain/models"
	"user_service/internal/domain/repository"
)

// RoleService представляет доменный сервис для работы с ролями.
type RoleService struct {
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
}

// NewRoleService создает новый экземпляр RoleService.
func NewRoleService(
	roleRepo repository.RoleRepository,
	userRoleRepo repository.UserRoleRepository,
) *RoleService {
	return &RoleService{
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
	}
}

// AssignRole назначает роль пользователю.
func (s *RoleService) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	// Проверяем, есть ли уже такая роль у пользователя
	hasRole, err := s.userRoleRepo.HasRole(ctx, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to check user role: %w", err)
	}

	if hasRole {
		return errors.New("user already has this role")
	}

	userRole := models.NewUserRole(userID, roleID)

	if err := s.userRoleRepo.AssignRole(ctx, userRole); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// RevokeRole отзывает роль у пользователя.
func (s *RoleService) RevokeRole(ctx context.Context, userID, roleID uuid.UUID) error {
	// Проверяем, есть ли такая роль у пользователя
	hasRole, err := s.userRoleRepo.HasRole(ctx, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to check user role: %w", err)
	}

	if !hasRole {
		return errors.New("user does not have this role")
	}

	if err := s.userRoleRepo.RevokeRole(ctx, userID, roleID); err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	return nil
}

// GetUserRoles получает роли пользователя.
func (s *RoleService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*models.Role, error) {
	roles, err := s.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return roles, nil
}
