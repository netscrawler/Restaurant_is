package menugrpc

import (
	menuv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/menu"
	"google.golang.org/grpc"
)

type serverAPI struct {
	menuv1.UnimplementedMenuServiceServer
}

func Register(
	gRPCServer *grpc.Server,
) {
	menuv1.RegisterMenuServiceServer(
		gRPCServer,
		&serverAPI{},
	)
}
