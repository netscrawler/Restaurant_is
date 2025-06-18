package usergrpc

import (
	"context"
	"fmt"

	service "user_service/internal/domain/app"

	userv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/user"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type serverAPI struct {
	userv1.UnimplementedUserServiceServer
	userAppService  *service.UserAppService
	staffAppService *service.StaffAppService
	roleAppService  *service.RoleAppService
}

// NewServerAPI создает новый экземпляр serverAPI
func newServerAPI(
	userAppService *service.UserAppService,
	staffAppService *service.StaffAppService,
	roleAppService *service.RoleAppService,
) *serverAPI {
	return &serverAPI{
		userAppService:  userAppService,
		staffAppService: staffAppService,
		roleAppService:  roleAppService,
	}
}

// Пользователи (клиенты)
func (s *serverAPI) CreateUser(
	ctx context.Context,
	r *userv1.CreateUserRequest,
) (*userv1.UserResponse, error) {
	req := &service.CreateUserRequest{
		Email:    r.GetEmail(),
		Phone:    r.GetPhone(),
		FullName: r.GetFullName(),
		Password: r.GetPassword(),
	}

	resp, err := s.userAppService.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &userv1.UserResponse{
		User: &userv1.User{
			Id:        resp.ID,
			Email:     resp.Email,
			Phone:     resp.Phone,
			FullName:  resp.FullName,
			IsActive:  resp.IsActive,
			CreatedAt: resp.CreatedAt,
			UpdatedAt: resp.UpdatedAt,
			Roles:     resp.Roles,
		},
	}, nil
}

func (s *serverAPI) GetUser(
	ctx context.Context,
	r *userv1.GetUserRequest,
) (*userv1.UserResponse, error) {
	var req *service.GetUserRequest

	switch identifier := r.GetIdentifier().(type) {
	case *userv1.GetUserRequest_Id:
		req = &service.GetUserRequest{
			ID: &identifier.Id,
		}
	case *userv1.GetUserRequest_Email:
		req = &service.GetUserRequest{
			Email: &identifier.Email,
		}
	case *userv1.GetUserRequest_Phone:
		req = &service.GetUserRequest{
			Phone: &identifier.Phone,
		}
	default:
		return nil, fmt.Errorf("no identifier provided")
	}

	resp, err := s.userAppService.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &userv1.UserResponse{
		User: &userv1.User{
			Id:        resp.ID,
			Email:     resp.Email,
			Phone:     resp.Phone,
			FullName:  resp.FullName,
			IsActive:  resp.IsActive,
			CreatedAt: resp.CreatedAt,
			UpdatedAt: resp.UpdatedAt,
			Roles:     resp.Roles,
		},
	}, nil
}

func (s *serverAPI) UpdateUser(
	ctx context.Context,
	r *userv1.UpdateUserRequest,
) (*userv1.UserResponse, error) {
	req := &service.UpdateUserRequest{
		ID: r.GetId(),
	}

	if r.Email != nil {
		email := r.GetEmail()
		req.Email = &email
	}
	if r.Phone != nil {
		phone := r.GetPhone()
		req.Phone = &phone
	}
	if r.FullName != nil {
		fullName := r.GetFullName()
		req.FullName = &fullName
	}
	if r.IsActive != nil {
		isActive := r.GetIsActive()
		req.IsActive = &isActive
	}

	resp, err := s.userAppService.UpdateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &userv1.UserResponse{
		User: &userv1.User{
			Id:        resp.ID,
			Email:     resp.Email,
			Phone:     resp.Phone,
			FullName:  resp.FullName,
			IsActive:  resp.IsActive,
			CreatedAt: resp.CreatedAt,
			UpdatedAt: resp.UpdatedAt,
			Roles:     resp.Roles,
		},
	}, nil
}

func (s *serverAPI) DeleteUser(
	ctx context.Context,
	r *userv1.DeleteUserRequest,
) (*emptypb.Empty, error) {
	req := &service.DeleteUserRequest{
		ID: r.GetId(),
	}

	err := s.userAppService.DeleteUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) ListUsers(
	ctx context.Context,
	r *userv1.ListUsersRequest,
) (*userv1.ListUsersResponse, error) {
	req := &service.ListUsersRequest{
		OnlyActive: r.GetOnlyActive(),
		Page:       r.GetPage(),
		PageSize:   r.GetPageSize(),
	}

	resp, err := s.userAppService.ListUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*userv1.User, len(resp.Users))
	for i, user := range resp.Users {
		users[i] = &userv1.User{
			Id:        user.ID,
			Email:     user.Email,
			Phone:     user.Phone,
			FullName:  user.FullName,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Roles:     user.Roles,
		}
	}

	return &userv1.ListUsersResponse{
		Users:      users,
		TotalCount: resp.TotalCount,
	}, nil
}

// Сотрудники
func (s *serverAPI) CreateStaff(
	ctx context.Context,
	r *userv1.CreateStaffRequest,
) (*userv1.StaffResponse, error) {
	req := &service.CreateStaffRequest{
		WorkEmail: r.GetWorkEmail(),
		WorkPhone: r.GetWorkPhone(),
		FullName:  r.GetFullName(),
		Position:  r.GetPosition(),
		Password:  r.GetPassword(),
	}

	resp, err := s.staffAppService.CreateStaff(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	return &userv1.StaffResponse{
		Staff: &userv1.Staff{
			Id:        resp.ID,
			WorkEmail: resp.WorkEmail,
			WorkPhone: resp.WorkPhone,
			FullName:  resp.FullName,
			Position:  resp.Position,
			IsActive:  resp.IsActive,
			HireDate:  resp.HireDate,
			Roles:     resp.Roles,
		},
	}, nil
}

func (s *serverAPI) UpdateStaff(
	ctx context.Context,
	r *userv1.UpdateStaffRequest,
) (*userv1.StaffResponse, error) {
	req := &service.UpdateStaffRequest{
		ID: r.GetId(),
	}

	if r.WorkPhone != nil {
		workPhone := r.GetWorkPhone()
		req.WorkPhone = &workPhone
	}
	if r.Position != nil {
		position := r.GetPosition()
		req.Position = &position
	}
	if r.IsActive != nil {
		isActive := r.GetIsActive()
		req.IsActive = &isActive
	}

	resp, err := s.staffAppService.UpdateStaff(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	return &userv1.StaffResponse{
		Staff: &userv1.Staff{
			Id:        resp.ID,
			WorkEmail: resp.WorkEmail,
			WorkPhone: resp.WorkPhone,
			FullName:  resp.FullName,
			Position:  resp.Position,
			IsActive:  resp.IsActive,
			HireDate:  resp.HireDate,
			Roles:     resp.Roles,
		},
	}, nil
}

func (s *serverAPI) ListStaff(
	ctx context.Context,
	r *userv1.ListStaffRequest,
) (*userv1.ListStaffResponse, error) {
	req := &service.ListStaffRequest{
		OnlyActive: r.GetOnlyActive(),
		Page:       r.GetPage(),
		PageSize:   r.GetPageSize(),
	}

	resp, err := s.staffAppService.ListStaff(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list staff: %w", err)
	}

	staff := make([]*userv1.Staff, len(resp.Staff))
	for i, s := range resp.Staff {
		staff[i] = &userv1.Staff{
			Id:        s.ID,
			WorkEmail: s.WorkEmail,
			WorkPhone: s.WorkPhone,
			FullName:  s.FullName,
			Position:  s.Position,
			IsActive:  s.IsActive,
			HireDate:  s.HireDate,
			Roles:     s.Roles,
		}
	}

	return &userv1.ListStaffResponse{
		Staff:      staff,
		TotalCount: resp.TotalCount,
	}, nil
}

// Роли
func (s *serverAPI) AssignRole(
	ctx context.Context,
	r *userv1.AssignRoleRequest,
) (*emptypb.Empty, error) {
	// Преобразуем int64 в string для UUID
	userID := fmt.Sprintf("%016x-0000-0000-0000-000000000000", r.GetUserId())

	req := &service.AssignRoleRequest{
		UserID: userID,
		RoleID: r.GetRole(),
	}

	err := s.roleAppService.AssignRole(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) RevokeRole(
	ctx context.Context,
	r *userv1.RevokeRoleRequest,
) (*emptypb.Empty, error) {
	// Преобразуем int64 в string для UUID
	userID := fmt.Sprintf("%016x-0000-0000-0000-000000000000", r.GetUserId())

	req := &service.RevokeRoleRequest{
		UserID: userID,
		RoleID: r.GetRole(),
	}

	err := s.roleAppService.RevokeRole(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke role: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func Register(
	gRPCServer *grpc.Server,
	userAppService *service.UserAppService,
	staffAppService *service.StaffAppService,
	roleAppService *service.RoleAppService,
) {
	serverAPI := newServerAPI(userAppService, staffAppService, roleAppService)
	userv1.RegisterUserServiceServer(gRPCServer, serverAPI)
}
