package dto

import "github.com/google/uuid"

type CreateCategoryReq struct {
	Name         string
	Description  string
	DisplayOrder int32
}

type UpdateCategoryReq struct {
	ID           uuid.UUID
	Name         *string
	Description  *string
	DisplayOrder *int32
	IsActive     *bool
}

type DeleteCategoryReq struct {
	ID    uuid.UUID
	Force bool
}

type ListCategoriesReq struct {
	OnlyActive bool
	Page       int32
	PageSize   int32
}

type Category struct {
	ID           uuid.UUID
	Name         string
	Description  string
	DisplayOrder int32
	IsActive     bool
}
