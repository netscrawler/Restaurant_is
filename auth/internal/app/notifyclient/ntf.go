package notifyclient

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/config"
	notify "github.com/netscrawler/RispProtos/proto/gen/go/v1/notify"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Notify notify.NotifyClient
	Conn   *grpc.ClientConn
}

func New(ctx context.Context, cfg config.NotifyClient) (*Client, error) {
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
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		Notify: notify.NewNotifyClient(conn),
		Conn:   conn,
	}, nil
}
