package authgrpc

import (
	"context"

	authv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
	"google.golang.org/grpc"
)

type Auth interface {
	LoginClient(ctx context.Context, req *authv1.LoginClientRequest) (*authv1.LoginResponse, error)
	LoginStaff(ctx context.Context, req *authv1.LoginStaffRequest) (*authv1.LoginResponse, error)
	LoginYandex(ctx context.Context, req *authv1.OAuthYandexRequest) (*authv1.LoginResponse, error)
	Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error)
	Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.LoginResponse, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServiceServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	authv1.RegisterAuthServiceServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) LoginClient(
	ctx context.Context,
	in *authv1.LoginClientRequest,
) (*authv1.LoginResponse, error) {
	return s.auth.LoginClient(ctx, in)
}

func (s *serverAPI) LoginStaff(
	ctx context.Context,
	in *authv1.LoginStaffRequest,
) (*authv1.LoginResponse, error) {
	return s.auth.LoginStaff(ctx, in)
}

func (s *serverAPI) LoginYandex(
	ctx context.Context,
	in *authv1.OAuthYandexRequest,
) (*authv1.LoginResponse, error) {
	return s.auth.LoginYandex(ctx, in)
}

func (s *serverAPI) Validate(
	ctx context.Context,
	in *authv1.ValidateRequest,
) (*authv1.ValidateResponse, error) {
	return s.auth.Validate(ctx, in)
}

func (s *serverAPI) Refresh(
	ctx context.Context,
	in *authv1.RefreshRequest,
) (*authv1.LoginResponse, error) {
	return s.auth.Refresh(ctx, in)
}
