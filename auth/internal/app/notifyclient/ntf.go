package notifyclient

import (
	"context"
	"time"

	notify "github.com/netscrawler/RispProtos/proto/gen/go/v1/notify"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Notify notify.NotifyClient
	Conn   *grpc.ClientConn
}

func New(ctx context.Context, address string) (*Client, error) {
	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.5,
				MaxDelay:   5 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
		grpc.WithBlock(), // блокирует до соединения
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		Notify: notify.NewNotifyClient(conn),
		Conn:   conn,
	}, nil
}
