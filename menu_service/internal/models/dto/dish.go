package dto

import (
	"time"

	"github.com/google/uuid"
	menuv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/menu"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateDishReq struct {
	Name           string
	Description    string
	Price          uint64
	CategoryID     uuid.UUID
	CookingTimeMin int32
	ImageURL       *string
	IsAvailable    bool
	Calories       *int32
}

type UpdateDishReq struct {
	ID             uuid.UUID
	Name           *string
	Description    *string
	Price          *uint64
	CategoryID     *int32
	CookingTimeMin *int32
	ImageURL       *string
	IsAvailable    *bool
	Calories       *int32
}

func NewUpdateDishReq(
	id uuid.UUID,
	name, desc *string,
	price *uint64,
	category, cookingTimeMin *int32,
	imageURL *string,
	isAvailable *bool,
	call *int32,
) *UpdateDishReq {
	return &UpdateDishReq{
		ID:             id,
		Name:           name,
		Description:    desc,
		Price:          price,
		CategoryID:     category,
		CookingTimeMin: cookingTimeMin,
		ImageURL:       imageURL,
		IsAvailable:    isAvailable,
		Calories:       call,
	}
}

type ListDishFilter struct {
	CategoryID    *int32
	OnlyAvailable bool
	Page          int32
	PageSize      int32
}

func NewListDishReq(categoryID *int32, available bool, page, pageSize int32) *ListDishFilter {
	return &ListDishFilter{
		CategoryID:    categoryID,
		OnlyAvailable: available,
		Page:          page,
		PageSize:      pageSize,
	}
}

type Dish struct {
	ID             uuid.UUID
	Name           string
	Description    string
	Price          uint64
	CategoryID     int32
	CookingTimeMin int32
	ImageURL       string
	IsAvailable    bool
	Calories       int32
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func (d *Dish) ToGRPCDto() *menuv1.Dish {
	return &menuv1.Dish{
		Id:             &menuv1.UUID{Value: d.ID.String()},
		Name:           d.Name,
		Description:    d.Name,
		Price:          d.Price,
		CategoryId:     d.CategoryID,
		CookingTimeMin: d.CookingTimeMin,
		ImageUrl:       d.ImageURL,
		IsAvailable:    d.IsAvailable,
		Calories:       d.Calories,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: d.CreatedAt.Unix(),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: d.UpdatedAt.Unix(),
		},
	}
}

func NewDish(
	name, desc string,
	price uint64,
	category int32,
	imageURL string,
	cookingTimeMin int32,
	calories int32,
) *Dish {
	created := time.Now()

	return &Dish{
		ID:             uuid.New(),
		Name:           name,
		Description:    desc,
		Price:          price,
		CategoryID:     category,
		CookingTimeMin: cookingTimeMin,
		ImageURL:       imageURL,
		IsAvailable:    false,
		Calories:       calories,
		CreatedAt:      &created,
		UpdatedAt:      &created,
	}
}
