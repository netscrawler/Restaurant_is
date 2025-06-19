package application

import (
	"context"
	"errors"

	"user_service/internal/domain/models"
	"user_service/internal/domain/service"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserAppService представляет application сервис для работы с пользователями.
type UserAppService struct {
	userService *service.UserService
}

// NewUserAppService создает новый экземпляр UserAppService.
func NewUserAppService(userService *service.UserService) *UserAppService {
	return &UserAppService{
		userService: userService,
	}
}

// CreateUserRequest представляет запрос на создание пользователя.
type CreateUserRequest struct {
	Email    string
	Phone    string
	FullName string
	Password string
}

// CreateUserResponse представляет ответ на создание пользователя.
type CreateUserResponse struct {
	ID        string
	Email     string
	Phone     string
	FullName  string
	IsActive  bool
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
	Roles     []string
}

// CreateUser создает нового пользователя.
func (s *UserAppService) CreateUser(
	ctx context.Context,
	req *CreateUserRequest,
) (*CreateUserResponse, error) {
	user, err := s.userService.CreateUser(ctx, req.Email, req.Phone, req.FullName)
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		ID:        user.ID.String(), // Преобразуем UUID в int64
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		IsActive:  user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Roles:     []string{}, // Пока пустой список ролей
	}, nil
}

// HandleUserCreatedEvent обрабатывает событие user_created из Kafka
func (s *UserAppService) HandleUserCreatedEvent(
	ctx context.Context,
	id, email, phone string,
) error {
	// Проверяем, существует ли пользователь с таким id/email/phone
	// Если нет — создаем, если есть — можно обновить или пропустить (по бизнес-логике)
	_, err := s.userService.GetUserByEmail(ctx, email)
	if err == nil {
		return nil // Уже есть, не создаём
	}
	_, err = s.userService.GetUserByPhone(ctx, phone)
	if err == nil {
		return nil // Уже есть, не создаём
	}
	// Создаём пользователя
	_, err = s.userService.CreateUser(ctx, email, phone, "")
	return err
}

// GetUserRequest представляет запрос на получение пользователя.
type GetUserRequest struct {
	ID    *string
	Email *string
	Phone *string
}

// GetUserResponse представляет ответ на получение пользователя.
type GetUserResponse struct {
	ID        string
	Email     string
	Phone     string
	FullName  string
	IsActive  bool
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
	Roles     []string
}

// GetUser получает пользователя по ID, email или телефону.
func (s *UserAppService) GetUser(
	ctx context.Context,
	req *GetUserRequest,
) (*GetUserResponse, error) {
	var user *models.User

	if req.ID != nil {
		// Создаем UUID из int64
		userID, err := uuid.Parse(*req.ID)
		if err != nil {
			return nil, err
		}
		user, err = s.userService.GetUser(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	if req.Email != nil {
		user, _ = s.userService.GetUserByEmail(ctx, *req.Email)
	}
	if req.Phone != nil {
		user, _ = s.userService.GetUserByPhone(ctx, *req.Phone)
	}

	if user == nil {
		return nil, errors.New("no identifier provided")
	}

	return &GetUserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		IsActive:  user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Roles:     []string{}, // Пока пустой список ролей
	}, nil
}

// UpdateUserRequest представляет запрос на обновление пользователя.
type UpdateUserRequest struct {
	ID       string
	Email    *string
	Phone    *string
	FullName *string
	IsActive *bool
}

// UpdateUserResponse представляет ответ на обновление пользователя.
type UpdateUserResponse struct {
	ID        string
	Email     string
	Phone     string
	FullName  string
	IsActive  bool
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
	Roles     []string
}

// UpdateUser обновляет данные пользователя.
func (s *UserAppService) UpdateUser(
	ctx context.Context,
	req *UpdateUserRequest,
) (*UpdateUserResponse, error) {
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

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
		ID:        user.ID.String(),
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		IsActive:  user.IsActive,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Roles:     []string{}, // Пока пустой список ролей
	}, nil
}

// DeleteUserRequest представляет запрос на удаление пользователя.
type DeleteUserRequest struct {
	ID string
}

// DeleteUser удаляет пользователя.
func (s *UserAppService) DeleteUser(ctx context.Context, req *DeleteUserRequest) error {
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return err
	}

	return s.userService.DeleteUser(ctx, userID)
}

// ListUsersRequest представляет запрос на получение списка пользователей.
type ListUsersRequest struct {
	OnlyActive bool
	Page       int32
	PageSize   int32
}

// ListUsersResponse представляет ответ на получение списка пользователей.
type ListUsersResponse struct {
	Users      []*GetUserResponse
	TotalCount int32
}

// ListUsers возвращает список пользователей.
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
			ID:        user.ID.String(),
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
