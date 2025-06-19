package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/order_service/internal/domain"
	"github.com/netscrawler/Restaurant_is/order_service/internal/models/dto"
	"github.com/netscrawler/Restaurant_is/order_service/internal/models/repository"
)

const AnswerTimeout = 15 * time.Second

type OrderProvider interface {
	Save(ctx context.Context, order *repository.Order) (uint64, error)
	Get(ctx context.Context, orderID uuid.UUID) (*repository.Order, error)
}

type EventSaver interface {
	Save(ctx context.Context, event *repository.Event) error
}

type DishGetter interface {
	Get(ctx context.Context, dishes []uuid.UUID) ([]*dto.Dish, error)
}

type Order struct {
	repo  OrderProvider
	event EventSaver
	dish  DishGetter
}

func NewOrder(repo OrderProvider, event EventSaver, dish DishGetter) *Order {
	return &Order{
		repo:  repo,
		event: event,
		dish:  dish,
	}
}

func (o *Order) Create(ctx context.Context, order *dto.OrderToCreate) (*dto.OrderCreated, error) {
	dishCtx, cancel := context.WithTimeout(ctx, AnswerTimeout)
	defer cancel()

	var dishIDs []uuid.UUID
	for _, v := range order.Items {
		dishIDs = append(dishIDs, v.Item)
	}

	dishes, err := o.dish.Get(dishCtx, dishIDs)
	if err != nil {
		panic(err)
	}

	dishMap := make(domain.DishList)

	for _, d := range dishes {
		for _, item := range order.Items {
			if d.ID == item.Item {
				dish := domain.NewDish(*d)
				dishMap[*dish] = item.Quanity
			}
		}
	}

	domainOrder, err := domain.NewOrder(
		order.UserID,
		dishMap,
		domain.OrderType(order.OrderType),
		"any",
	)
	if err != nil {
		return nil, err
	}

	orderNum, err := o.repo.Save(ctx, repository.NewOrder(domainOrder))
	if err != nil {
		return nil, err
	}

	domainOrder.SetNUM(orderNum)

	orderCreated, _ := dto.NewOrderCreated(
		domainOrder.ID(),
		domainOrder.NUM(),
		domainOrder.UserID(),
		domainOrder.Price(),
		domainOrder.Status(),
	)

	return orderCreated, nil
}

func (o *Order) UpdateStatus(ctx context.Context, orderID uuid.UUID, newStatus string) error {
	// order, err := o.repo.Get(ctx, orderID)
	// if err != nil {
	// 	return err
	// }
	// _, err = o.repo.Save(ctx, repository.NewOrder(order))
	return nil
}

func (o *Order) GetOrder(ctx context.Context, orderID string) (*repository.Order, error) {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return nil, err
	}
	return o.repo.Get(ctx, id)
}

func (o *Order) ListOrders(ctx context.Context, filter dto.OrderFilter) ([]*repository.Order, error) {
	return o.repo.(interface {
		ListOrders(context.Context, dto.OrderFilter) ([]*repository.Order, error)
	}).ListOrders(ctx, filter)
}

func (o *Order) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return err
	}
	order, err := o.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	order.Status = status
	order.UpdatedAt = time.Now()
	type updater interface {
		Update(context.Context, *repository.Order) error
	}
	if u, ok := o.repo.(updater); ok {
		return u.Update(ctx, order)
	}
	return nil
}
