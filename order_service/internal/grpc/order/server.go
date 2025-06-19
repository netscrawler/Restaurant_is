package ordergrpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/order_service/internal/models/dto"
	"github.com/netscrawler/Restaurant_is/order_service/internal/models/repository"
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
	Create(ctx context.Context, order *dto.OrderToCreate) (*dto.OrderCreated, error)
	GetOrder(ctx context.Context, orderID string) (*repository.Order, error)
	ListOrders(ctx context.Context, filter dto.OrderFilter) ([]*repository.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
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

// Основные операции с заказами.
func (s *serverAPI) CreateOrder(
	ctx context.Context,
	r *orderv1.CreateOrderRequest,
) (*orderv1.OrderResponse, error) {
	userId, err := uuid.Parse(r.GetUserId().GetValue())
	if err != nil {
	}

	items := make([]dto.OrderItem, len(r.GetItems()))

	for i, item := range r.GetItems() {
		itemID, err := uuid.Parse(item.GetDishId().GetValue())
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

	createdOrder, err := s.orderProvider.Create(ctx, orderToCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &orderv1.OrderResponse{
		Id:          &orderv1.UUID{Value: createdOrder.ID},
		Status:      createdOrder.Status,
		TotalAmount: int64(createdOrder.Total),
	}, nil
}

func statusToProto(s string) orderv1.OrderStatus {
	switch s {
	case "created":
		return orderv1.OrderStatus_ORDER_STATUS_CREATED
	case "process":
		return orderv1.OrderStatus_ORDER_STATUS_CONFIRMED
	case "on_kitchen":
		return orderv1.OrderStatus_ORDER_STATUS_COOKING
	case "delivery":
		return orderv1.OrderStatus_ORDER_STATUS_READY
	case "delivered":
		return orderv1.OrderStatus_ORDER_STATUS_DELIVERED
	case "declined":
		return orderv1.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return orderv1.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func statusFromProto(s orderv1.OrderStatus) string {
	switch s {
	case orderv1.OrderStatus_ORDER_STATUS_CREATED:
		return "created"
	case orderv1.OrderStatus_ORDER_STATUS_CONFIRMED:
		return "process"
	case orderv1.OrderStatus_ORDER_STATUS_COOKING:
		return "on_kitchen"
	case orderv1.OrderStatus_ORDER_STATUS_READY:
		return "delivery"
	case orderv1.OrderStatus_ORDER_STATUS_DELIVERED:
		return "delivered"
	case orderv1.OrderStatus_ORDER_STATUS_CANCELLED:
		return "declined"
	default:
		return "unspecified"
	}
}

func (s *serverAPI) GetOrder(
	ctx context.Context,
	r *orderv1.GetOrderRequest,
) (*orderv1.OrderResponse, error) {
	order, err := s.orderProvider.GetOrder(ctx, r.GetOrderId().GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}
	return &orderv1.OrderResponse{
		Id:          &orderv1.UUID{Value: order.ID.String()},
		Status:      statusToProto(order.Status).String(),
		TotalAmount: int64(order.Price),
	}, nil
}

func (s *serverAPI) ListOrders(
	ctx context.Context,
	r *orderv1.ListOrdersRequest,
) (*orderv1.ListOrdersResponse, error) {
	filter := dto.OrderFilter{}
	if r.GetUserId() != nil {
		filter.UserID = r.GetUserId().GetValue()
	}
	if r.Status != nil {
		filter.Status = statusFromProto(r.GetStatus())
	}
	filter.Limit = int(r.GetPageSize())
	filter.Offset = int(r.GetPage()) * filter.Limit

	orders, err := s.orderProvider.ListOrders(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}
	resp := &orderv1.ListOrdersResponse{}
	for _, order := range orders {
		resp.Orders = append(resp.Orders, &orderv1.Order{
			Id:          &orderv1.UUID{Value: order.ID.String()},
			Status:      statusToProto(order.Status),
			TotalAmount: int64(order.Price),
		})
	}
	return resp, nil
}

func (s *serverAPI) UpdateOrderStatus(
	ctx context.Context,
	r *orderv1.UpdateOrderStatusRequest,
) (*emptypb.Empty, error) {
	err := s.orderProvider.UpdateOrderStatus(ctx, r.GetOrderId().GetValue(), statusFromProto(r.GetStatus()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// Платежи.
func (s *serverAPI) InitiatePayment(
	ctx context.Context,
	r *orderv1.PaymentRequest,
) (*orderv1.PaymentResponse, error) {
	return &orderv1.PaymentResponse{Status: "success"}, nil
}

func (s *serverAPI) ProcessPaymentCallback(
	ctx context.Context,
	r *orderv1.PaymentCallbackRequest,
) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// История и отчетность.
func (s *serverAPI) GetOrderHistory(
	ctx context.Context,
	r *orderv1.GetOrderRequest,
) (*orderv1.OrderHistoryResponse, error) {
	return &orderv1.OrderHistoryResponse{}, nil
}
