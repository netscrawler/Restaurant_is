package menugrpc

import (
	"context"

	menuv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/menu"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryProvider interface{}

type MenuProvider interface{}

type PromotionProvider interface{}

type serverAPI struct {
	category CategoryProvider
	menu     MenuProvider
	dish     PromotionProvider
	menuv1.UnimplementedMenuServiceServer
}

// Категории.
func (s *serverAPI) CreateCategory(
	ctx context.Context,
	in *menuv1.CategoryRequest,
) (*menuv1.CategoryResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) UpdateCategory(
	ctx context.Context,
	in *menuv1.UpdateCategoryRequest,
) (*menuv1.CategoryResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) ListCategories(
	ctx context.Context,
	in *menuv1.ListCategoriesRequest,
) (*menuv1.ListCategoriesResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) DeleteCategory(
	ctx context.Context,
	in *menuv1.DeleteCategoryRequest,
) (*emptypb.Empty, error) {
	panic("not implemented") // TODO: Implement
}

// Блюда.
func (s *serverAPI) CreateDish(
	ctx context.Context,
	in *menuv1.DishRequest,
) (*menuv1.DishResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) UpdateDish(
	ctx context.Context,
	in *menuv1.UpdateDishRequest,
) (*menuv1.DishResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) GetDish(
	ctx context.Context,
	in *menuv1.GetDishRequest,
) (*menuv1.DishResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) ListDishes(
	ctx context.Context,
	in *menuv1.ListDishesRequest,
) (*menuv1.ListDishesResponse, error) {
	panic("not implemented") // TODO: Implement
}

// Акции.
func (s *serverAPI) CreatePromotion(
	ctx context.Context,
	in *menuv1.PromotionRequest,
) (*menuv1.PromotionResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) UpdatePromotion(
	ctx context.Context,
	in *menuv1.UpdatePromotionRequest,
) (*menuv1.PromotionResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) ListActivePromotions(
	ctx context.Context,
	in *emptypb.Empty,
) (*menuv1.ListPromotionsResponse, error) {
	panic("not implemented") // TODO: Implement
}

// Изображения.
func (s *serverAPI) GenerateUploadURL(
	ctx context.Context,
	in *menuv1.ImageRequest,
) (*menuv1.ImageResponse, error) {
	panic("not implemented") // TODO: Implement
}

func Register(
	gRPCServer *grpc.Server,
	menu MenuProvider,
	dish PromotionProvider,
	category CategoryProvider,
) {
	menuv1.RegisterMenuServiceServer(
		gRPCServer,
		&serverAPI{
			category: category,
			menu:     menu,
			dish:     dish,
		},
	)
}
