package grpcinfra

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/order_service/internal/config"
	"github.com/netscrawler/Restaurant_is/order_service/internal/models/dto"
	menuclient "github.com/netscrawler/RispProtos/proto/gen/go/v1/menu"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type MenuClient struct {
	menu menuclient.MenuServiceClient
	Conn *grpc.ClientConn
	log  *slog.Logger
}

func New(ctx context.Context, cfg config.MenuClient) (*MenuClient, error) {
	conn, err := grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  cfg.BaseDelay,
				Multiplier: cfg.Multiplier,
				MaxDelay:   cfg.MaxDelay,
			},
			MinConnectTimeout: cfg.MinConnectTimeout,
		}),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	return &MenuClient{
		Conn: conn,
		menu: menuclient.NewMenuServiceClient(conn),
	}, nil
}

func (m *MenuClient) Get(ctx context.Context, dishes []uuid.UUID) ([]*dto.Dish, error) {
	type result struct {
		dish *dto.Dish
		err  error
	}

	out := make(chan result, len(dishes))

	var wg sync.WaitGroup

	wg.Add(len(dishes))

	for _, d := range dishes {
		go func(d uuid.UUID) {
			defer wg.Done()

			dish, err := m.get(ctx, d)
			out <- result{dish: dish, err: err}
		}(d)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	dishList := make([]*dto.Dish, 0, len(dishes))

	for r := range out {
		if r.err != nil {
			return nil, r.err
		}

		dishList = append(dishList, r.dish)
	}

	return dishList, nil
}

func (m *MenuClient) get(ctx context.Context, dish uuid.UUID) (*dto.Dish, error) {
	resp, err := m.menu.GetDish(ctx, &menuclient.GetDishRequest{
		DishId: &menuclient.UUID{
			Value: dish.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	d, err := dto.NewDish(
		resp.GetDish().GetId().GetValue(),
		resp.GetDish().GetName(),
		resp.GetDish().GetPrice(),
	)
	if err != nil {
		return nil, err
	}

	return d, nil
}
