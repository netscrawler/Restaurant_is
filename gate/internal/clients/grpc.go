package clients

import (
	"fmt"
	"time"

	authv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
	menuv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/menu"
	orderv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/order"
	userv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	AuthClient  authv1.AuthClient
	UserClient  userv1.UserServiceClient
	MenuClient  menuv1.MenuServiceClient
	OrderClient orderv1.OrderServiceClient

	// Храним соединения для их закрытия
	connections []*grpc.ClientConn
}

func NewGRPCClients(config map[string]string) (*GRPCClients, error) {
	clients := &GRPCClients{
		connections: make([]*grpc.ClientConn, 0),
	}

	// Создаем соединения для каждого сервиса
	services := []string{"auth", "user", "menu", "order"}

	for _, service := range services {
		address := config[service]
		if address == "" {
			return nil, fmt.Errorf("no address configured for service: %s", service)
		}

		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithIdleTimeout(5*time.Second),

			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s service: %w", service, err)
		}

		// Сохраняем соединение для последующего закрытия
		clients.connections = append(clients.connections, conn)

		// Создаем конкретные клиенты
		switch service {
		case "auth":
			clients.AuthClient = authv1.NewAuthClient(conn)
		case "user":
			clients.UserClient = userv1.NewUserServiceClient(conn)
		case "menu":
			clients.MenuClient = menuv1.NewMenuServiceClient(conn)
		case "order":
			clients.OrderClient = orderv1.NewOrderServiceClient(conn)
		}
	}

	return clients, nil
}

func (c *GRPCClients) Close() {
	// Закрываем все соединения
	for _, conn := range c.connections {
		if conn != nil {
			conn.Close()
		}
	}
}
