package application

import (
	"context"
	"fmt"

	"user_service/internal/domain/models"
	"user_service/internal/domain/service"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserAppService представляет application сервис для работы с пользователями
type UserAppService struct {
	userService *service.UserService
}

// NewUserAppService создает новый экземпляр UserAppService
func NewUserAppService(userService *service.UserService) *UserAppService {
	return &UserAppService{
		userService: userService,
	}
}

// CreateUserRequest представляет запрос на создание пользователя
type CreateUserRequest struct {
	Email    string
	Phone    string
	FullName string
	Password string
}

// CreateUserResponse представляет ответ на создание пользователя
type CreateUserResponse struct {
	ID        int64
	Email     string
	Phone     string
	FullName  string
	IsActive  bool
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
	Roles     []string
}

// CreateUser создает нового пользователя
func (s *UserAppService) CreateUser(
	ctx context.Context,
	req *CreateUserRequest,
) (*CreateUserResponse, error) {
	user, err := s.userService.CreateUser(ctx, req.Email, req.Phone, req.FullName)
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		ID:        int64(user.ID.ID()), // Преобразуем UUID в int64
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		IsActive:  user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Roles:     []string{}, // Пока пустой список ролей
	}, nil
}

// GetUserRequest представляет запрос на получение пользователя
type GetUserRequest struct {
	ID    *int64
	Email *string
	Phone *string
}

// GetUserResponse представляет ответ на получение пользователя
type GetUserResponse struct {
	ID        int64
	Email     string
	Phone     string
	FullName  string
	IsActive  bool
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
	Roles     []string
}

// GetUser получает пользователя по ID, email или телефону
func (s *UserAppService) GetUser(
	ctx context.Context,
	req *GetUserRequest,
) (*GetUserResponse, error) {
	var user *models.User
	var err error

	if req.ID != nil {
		// Создаем UUID из int64
		userID := uuid.MustParse(fmt.Sprintf("%016x-0000-0000-0000-000000000000", *req.ID))
		user, err = s.userService.GetUser(ctx, userID)
	} else if req.Email != nil {
		user, err = s.userService.GetUserByEmail(ctx, *req.Email)
	} else if req.Phone != nil {
		user, err = s.userService.GetUserByPhone(ctx, *req.Phone)
	} else {
		return nil, fmt.Errorf("no identifier provided")
	}

	if err != nil {
		return nil, err
	}

	return &GetUserResponse{
		ID:        int64(user.ID.ID()),
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		IsActive:  user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Roles:     []string{}, // Пока пустой список ролей
	}, nil
}

// UpdateUserRequest представляет запрос на обновление пользователя
type UpdateUserRequest struct {
	ID       int64
	Email    *string
	Phone    *string
	FullName *string
	IsActive *bool
}

// UpdateUserResponse представляет ответ на обновление пользователя
type UpdateUserResponse struct {
	ID        int64
	Email     string
	Phone     string
	FullName  string
	IsActive  bool
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
	Roles     []string
}

// UpdateUser обновляет данные пользователя
func (s *UserAppService) UpdateUser(
	ctx context.Context,
	req *UpdateUserRequest,
) (*UpdateUserResponse, error) {
	userID := uuid.MustParse(fmt.Sprintf("%016x-0000-0000-0000-000000000000", req.ID))

	email := ""
	if req.Email != nil {
		email = *req.Email
	}

	phone := ""
	if req.Phone != nil {
		phone = *req.Phone
	}

	fullName := ""
	if req.FullName != nil {
		fullName = *req.FullName
	}

	user, err := s.userService.UpdateUser(ctx, userID, email, phone, fullName, req.IsActive)
	if err != nil {
		return nil, err
	}

	return &UpdateUserResponse{
		ID:        int64(user.ID.ID()),
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		IsActive:  user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Roles:     []string{}, // Пока пустой список ролей
	}, nil
}

// DeleteUserRequest представляет запрос на удаление пользователя
type DeleteUserRequest struct {
	ID int64
}

// DeleteUser удаляет пользователя
func (s *UserAppService) DeleteUser(ctx context.Context, req *DeleteUserRequest) error {
	userID := uuid.MustParse(fmt.Sprintf("%016x-0000-0000-0000-000000000000", req.ID))
	return s.userService.DeleteUser(ctx, userID)
}

// ListUsersRequest представляет запрос на получение списка пользователей
type ListUsersRequest struct {
	OnlyActive bool
	Page       int32
	PageSize   int32
}

// ListUsersResponse представляет ответ на получение списка пользователей
type ListUsersResponse struct {
	Users      []*GetUserResponse
	TotalCount int32
}

// ListUsers возвращает список пользователей
func (s *UserAppService) ListUsers(
	ctx context.Context,
	req *ListUsersRequest,
) (*ListUsersResponse, error) {
	offset := int(req.Page) * int(req.PageSize)
	limit := int(req.PageSize)

	users, total, err := s.userService.ListUsers(ctx, req.OnlyActive, offset, limit)
	if err != nil {
		return nil, err
	}

	userResponses := make([]*GetUserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &GetUserResponse{
			ID:        int64(user.ID.ID()),
			Email:     user.Email,
			Phone:     user.Phone,
			FullName:  user.FullName,
			IsActive:  user.IsActive,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			Roles:     []string{}, // Пока пустой список ролей
		}
	}

	return &ListUsersResponse{
		Users:      userResponses,
		TotalCount: int32(total),
	}, nil
}
