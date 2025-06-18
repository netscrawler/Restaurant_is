package application

import (
	"context"
	"fmt"

	"user_service/internal/domain/service"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// StaffAppService представляет application сервис для работы с сотрудниками
type StaffAppService struct {
	staffService *service.StaffService
}

// NewStaffAppService создает новый экземпляр StaffAppService
func NewStaffAppService(staffService *service.StaffService) *StaffAppService {
	return &StaffAppService{
		staffService: staffService,
	}
}

// CreateStaffRequest представляет запрос на создание сотрудника
type CreateStaffRequest struct {
	WorkEmail string
	WorkPhone string
	FullName  string
	Position  string
	Password  string
}

// CreateStaffResponse представляет ответ на создание сотрудника
type CreateStaffResponse struct {
	ID        int64
	WorkEmail string
	WorkPhone string
	FullName  string
	Position  string
	IsActive  bool
	HireDate  *timestamppb.Timestamp
	Roles     []string
}

// CreateStaff создает нового сотрудника
func (s *StaffAppService) CreateStaff(
	ctx context.Context,
	req *CreateStaffRequest,
) (*CreateStaffResponse, error) {
	staff, err := s.staffService.CreateStaff(
		ctx,
		req.WorkEmail,
		req.WorkPhone,
		req.FullName,
		req.Position,
	)
	if err != nil {
		return nil, err
	}

	return &CreateStaffResponse{
		ID:        int64(staff.ID.ID()),
		WorkEmail: staff.WorkEmail,
		WorkPhone: staff.WorkPhone,
		FullName:  staff.FullName,
		Position:  staff.Position,
		IsActive:  staff.IsActive,
		HireDate:  timestamppb.New(staff.HireDate),
		Roles:     []string{}, // Пока пустой список ролей
	}, nil
}

// UpdateStaffRequest представляет запрос на обновление сотрудника
type UpdateStaffRequest struct {
	ID        int64
	WorkPhone *string
	Position  *string
	IsActive  *bool
}

// UpdateStaffResponse представляет ответ на обновление сотрудника
type UpdateStaffResponse struct {
	ID        int64
	WorkEmail string
	WorkPhone string
	FullName  string
	Position  string
	IsActive  bool
	HireDate  *timestamppb.Timestamp
	Roles     []string
}

// UpdateStaff обновляет данные сотрудника
func (s *StaffAppService) UpdateStaff(
	ctx context.Context,
	req *UpdateStaffRequest,
) (*UpdateStaffResponse, error) {
	staffID := uuid.MustParse(fmt.Sprintf("%016x-0000-0000-0000-000000000000", req.ID))

	workPhone := ""
	if req.WorkPhone != nil {
		workPhone = *req.WorkPhone
	}

	position := ""
	if req.Position != nil {
		position = *req.Position
	}

	staff, err := s.staffService.UpdateStaff(ctx, staffID, workPhone, position, req.IsActive)
	if err != nil {
		return nil, err
	}

	return &UpdateStaffResponse{
		ID:        int64(staff.ID.ID()),
		WorkEmail: staff.WorkEmail,
		WorkPhone: staff.WorkPhone,
		FullName:  staff.FullName,
		Position:  staff.Position,
		IsActive:  staff.IsActive,
		HireDate:  timestamppb.New(staff.HireDate),
		Roles:     []string{}, // Пока пустой список ролей
	}, nil
}

// ListStaffRequest представляет запрос на получение списка сотрудников
type ListStaffRequest struct {
	OnlyActive bool
	Page       int32
	PageSize   int32
}

// ListStaffResponse представляет ответ на получение списка сотрудников
type ListStaffResponse struct {
	Staff      []*CreateStaffResponse
	TotalCount int32
}

// ListStaff возвращает список сотрудников
func (s *StaffAppService) ListStaff(
	ctx context.Context,
	req *ListStaffRequest,
) (*ListStaffResponse, error) {
	offset := int(req.Page) * int(req.PageSize)
	limit := int(req.PageSize)

	staff, total, err := s.staffService.ListStaff(ctx, req.OnlyActive, offset, limit)
	if err != nil {
		return nil, err
	}

	staffResponses := make([]*CreateStaffResponse, len(staff))
	for i, st := range staff {
		staffResponses[i] = &CreateStaffResponse{
			ID:        int64(st.ID.ID()),
			WorkEmail: st.WorkEmail,
			WorkPhone: st.WorkPhone,
			FullName:  st.FullName,
			Position:  st.Position,
			IsActive:  st.IsActive,
			HireDate:  timestamppb.New(st.HireDate),
			Roles:     []string{}, // Пока пустой список ролей
		}
	}

	return &ListStaffResponse{
		Staff:      staffResponses,
		TotalCount: int32(total),
	}, nil
}
