package menugrpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/domain"
	"github.com/netscrawler/Restaurant_is/menu_service/internal/models/dto"
	menuv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/menu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DishProvider interface {
	// Create and save new dish.
	Create(ctx context.Context, dish *dto.Dish) (*dto.Dish, error)
	// Get dish by uuid.
	Get(ctx context.Context, dishID uuid.UUID) (*dto.Dish, error)
	// Update dish in storage.
	Update(ctx context.Context, dish *dto.UpdateDishReq) (*dto.Dish, error)
	// List returns dish by filter, offset and limit.
	List(ctx context.Context, filter *dto.ListDishFilter) ([]dto.Dish, error)
}

type ImageUrlCreator interface {
	// CreateURL generate s3 pre-signed url to save image.
	CreateURL(
		ctx context.Context,
		filename, contentType string,
	) (url string, objKey string, err error)
}

type serverAPI struct {
	dish  DishProvider
	image ImageUrlCreator
	menuv1.UnimplementedMenuServiceServer
}

// Блюда.
func (s *serverAPI) CreateDish(
	ctx context.Context,
	in *menuv1.DishRequest,
) (*menuv1.DishResponse, error) {
	dish := dto.NewDish(
		in.GetName(),
		in.GetDescription(),
		in.GetPrice(),
		in.GetCategoryId(),
		in.GetImageUrl(),
		in.GetCookingTimeMin(),
		in.GetCalories(),
	)

	createdDish, err := s.dish.Create(ctx, dish)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &menuv1.DishResponse{
		Dish: createdDish.ToGRPCDto(),
	}, nil
}

func (s *serverAPI) UpdateDish(
	ctx context.Context,
	in *menuv1.UpdateDishRequest,
) (*menuv1.DishResponse, error) {
	dishID, err := uuid.ParseBytes(in.GetId().GetValue())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidUUID.Error())
	}

	dishUpdate := dto.NewUpdateDishReq(
		dishID,
		in.Name,
		in.Description,
		in.Price,
		in.CategoryId,
		in.CookingTimeMin,
		in.ImageUrl,
		in.IsAvailable,
		in.Calories,
	)

	dish, err := s.dish.Update(ctx, dishUpdate)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &menuv1.DishResponse{
		Dish: dish.ToGRPCDto(),
	}, nil
}

func (s *serverAPI) GetDish(
	ctx context.Context,
	in *menuv1.GetDishRequest,
) (*menuv1.DishResponse, error) {
	dishID, err := uuid.ParseBytes(in.GetDishId().GetValue())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidUUID.Error())
	}

	dish, err := s.dish.Get(ctx, dishID)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &menuv1.DishResponse{
		Dish: dish.ToGRPCDto(),
	}, nil
}

func (s *serverAPI) ListDishes(
	ctx context.Context,
	in *menuv1.ListDishesRequest,
) (*menuv1.ListDishesResponse, error) {
	req := dto.NewListDishReq(
		in.CategoryId,
		in.OnlyAvailable,
		in.Page,
		in.PageSize,
	)

	dishes, err := s.dish.List(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dishList := make([]*menuv1.Dish, 0, len(dishes))

	for _, v := range dishes {
		dishList = append(dishList, v.ToGRPCDto())
	}

	return &menuv1.ListDishesResponse{
		Dishes:     dishList,
		TotalCount: int32(len(dishes)),
	}, nil
}

// Изображения.
func (s *serverAPI) GenerateUploadURL(
	ctx context.Context,
	in *menuv1.ImageRequest,
) (*menuv1.ImageResponse, error) {
	url, objKey, err := s.image.CreateURL(ctx, in.GetFilename(), in.GetContentType())
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &menuv1.ImageResponse{
		Url:       url,
		ObjectKey: objKey,
	}, nil
}

func Register(
	gRPCServer *grpc.Server,
	dish DishProvider,
	image ImageUrlCreator,
) {
	menuv1.RegisterMenuServiceServer(
		gRPCServer,
		&serverAPI{
			dish:  dish,
			image: image,
		},
	)
}
