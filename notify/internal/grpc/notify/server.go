package notifygrpc

import (
	"context"

	notifyv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/notify"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotifySender interface {
	Send(ctx context.Context, recipient, message string) error
}

type serverAPI struct {
	notifyv1.UnimplementedNotifyServer
	notifyer NotifySender
}

func (s *serverAPI) Send(
	ctx context.Context,
	in *notifyv1.SendRequest,
) (*notifyv1.SendResponse, error) {
	err := s.notifyer.Send(ctx, in.GetPhone(), in.GetData())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &notifyv1.SendResponse{}, nil
}

func Register(
	gRPCServer *grpc.Server,
	notifyer NotifySender,
) {
	notifyv1.RegisterNotifyServer(
		gRPCServer,
		&serverAPI{notifyer: notifyer},
	)
}
