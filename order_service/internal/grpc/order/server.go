package ordergrpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/order_service/internal/models/dto"
	orderv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type serverAPI struct {
	orderv1.UnimplementedOrderServiceServer
	orderProvider OrderProvider
}

type OrderProvider interface {
	CreateOrder(ctx context.Context, order *dto.OrderToCreate) (*dto.OrderCreated, error)
	// GetOrder(ctx context.Context, orderID string) (dto.Order, error)
	// ListOrders(ctx context.Context, filter dto.OrderFilter) ([]dto.Order, error)
	// UpdateOrderStatus(ctx context.Context, orderID string, status dto.OrderStatus) error
}

func Register(
	gRPCServer *grpc.Server,
	orderProvider OrderProvider,
) {
	orderv1.RegisterOrderServiceServer(
		gRPCServer,
		&serverAPI{orderProvider: orderProvider},
	)
}

// Основные операции с заказами
func (s *serverAPI) CreateOrder(
	ctx context.Context,
	r *orderv1.CreateOrderRequest,
) (*orderv1.OrderResponse, error) {
	userId, err := uuid.ParseBytes(r.GetUserId().GetValue())
	if err != nil {
	}

	items := make([]dto.OrderItem, len(r.GetItems()))
	for i, item := range r.GetItems() {
		itemID, err := uuid.ParseBytes(item.GetDishId().GetValue())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid item ID: %v", err)
		}

		items[i] = dto.OrderItem{
			Item:    itemID,
			Quanity: uint8(item.GetQuantity()),
		}
	}

	orderToCreate := dto.NewOrder(
		userId,
		orderv1.OrderType_name[int32(r.GetOrderType())],
		r.GetDeliveryAddress(),
		items,
	)

	createdOrder, err := s.orderProvider.CreateOrder(ctx, orderToCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &orderv1.OrderResponse{
		Id:          &orderv1.UUID{Value: createdOrder.ID},
		Status:      createdOrder.Status,
		TotalAmount: int64(createdOrder.Total),
	}, nil
}

func (s *serverAPI) GetOrder(
	ctx context.Context,
	r *orderv1.GetOrderRequest,
) (*orderv1.OrderResponse, error) {
	// order, err := s.orderProvider.GetOrder(ctx, r.GetOrderId())
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	// }
	// return &orderv1.OrderResponse{
	// 	Order: &orderv1.Order{
	// 		Id:        &orderv1.UUID{Value: order.ID},
	// 		UserId:    &orderv1.UUID{Value: order.UserID},
	// 		ProductId: &orderv1.UUID{Value: order.ProductID},
	// 		Quantity:  int32(order.Quantity),
	// 		Price:     uint64(order.Price),
	// 		Status:    string(order.Status),
	// 	},
	// }, nil
	panic("err")
}

func (s *serverAPI) ListOrders(
	ctx context.Context,
	r *orderv1.ListOrdersRequest,
) (*orderv1.ListOrdersResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) UpdateOrderStatus(
	ctx context.Context,
	r *orderv1.UpdateOrderStatusRequest,
) (*emptypb.Empty, error) {
	panic("not implemented") // TODO: Implement
}

// Платежи
func (s *serverAPI) InitiatePayment(
	ctx context.Context,
	r *orderv1.PaymentRequest,
) (*orderv1.PaymentResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *serverAPI) ProcessPaymentCallback(
	ctx context.Context,
	r *orderv1.PaymentCallbackRequest,
) (*emptypb.Empty, error) {
	panic("not implemented") // TODO: Implement
}

// История и отчетность
func (s *serverAPI) GetOrderHistory(
	ctx context.Context,
	r *orderv1.GetOrderRequest,
) (*orderv1.OrderHistoryResponse, error) {
	panic("not implemented") // TODO: Implement
}
