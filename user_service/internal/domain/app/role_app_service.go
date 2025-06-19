package application

import (
	"context"

	"github.com/google/uuid"
	"user_service/internal/domain/service"
)

// RoleAppService представляет application сервис для работы с ролями.
type RoleAppService struct {
	roleService *service.RoleService
}

// NewRoleAppService создает новый экземпляр RoleAppService.
func NewRoleAppService(roleService *service.RoleService) *RoleAppService {
	return &RoleAppService{
		roleService: roleService,
	}
}

// AssignRoleRequest представляет запрос на назначение роли.
type AssignRoleRequest struct {
	UserID string
	RoleID string
}

// AssignRole назначает роль пользователю.
func (s *RoleAppService) AssignRole(ctx context.Context, req *AssignRoleRequest) error {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return err
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return err
	}

	return s.roleService.AssignRole(ctx, userID, roleID)
}

// RevokeRoleRequest представляет запрос на отзыв роли.
type RevokeRoleRequest struct {
	UserID string
	RoleID string
}

// RevokeRole отзывает роль у пользователя.
func (s *RoleAppService) RevokeRole(ctx context.Context, req *RevokeRoleRequest) error {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return err
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return err
	}

	return s.roleService.RevokeRole(ctx, userID, roleID)
}

// GetUserRolesRequest представляет запрос на получение ролей пользователя.
type GetUserRolesRequest struct {
	UserID string
}

// GetUserRolesResponse представляет ответ на получение ролей пользователя.
type GetUserRolesResponse struct {
	Roles []*RoleResponse
}

// RoleResponse представляет роль в ответе.
type RoleResponse struct {
	ID          string
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}

// GetUserRoles получает роли пользователя.
func (s *RoleAppService) GetUserRoles(
	ctx context.Context,
	req *GetUserRolesRequest,
) (*GetUserRolesResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	roles, err := s.roleService.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	roleResponses := make([]*RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = &RoleResponse{
			ID:          role.ID.String(),
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &GetUserRolesResponse{
		Roles: roleResponses,
	}, nil
}
