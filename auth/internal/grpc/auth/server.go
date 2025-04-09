package authgrpc

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/service"
	authv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
	"google.golang.org/grpc"
)

// type Auth interface {
// 	LoginClientInit(
// 		ctx context.Context,
// 		req *authv1.LoginClientRequest,
// 	) (*authv1.LoginInitResponse, error)
//
// 	LoginClientConfirm(
// 		context.Context,
// 		*authv1.LoginClientConfirmRequest,
// 	) (*authv1.LoginResponse, error)
// 	LoginStaff(ctx context.Context, req *authv1.LoginStaffRequest) (*authv1.LoginResponse, error)
// 	LoginYandex(ctx context.Context, req *authv1.OAuthYandexRequest) (*authv1.LoginResponse, error)
// 	Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error)
// 	Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.LoginResponse, error)
// }

type serverAPI struct {
	authv1.UnimplementedAuthServiceServer
	auth *service.AuthService
}

func Register(gRPCServer *grpc.Server, auth *service.AuthService) {
	authv1.RegisterAuthServiceServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) LoginClientInit(
	ctx context.Context,
	in *authv1.LoginClientRequest,
) (*authv1.LoginInitResponse, error) {
	err := s.auth.LoginClientInit(ctx, in.Phone)
	return &authv1.LoginInitResponse{
		Status:    "ok",
		ExpiresIn: 10000,
		Error:     "",
	}, err
}

// func (s *serverAPI) LoginStaff(
// 	ctx context.Context,
// 	in *authv1.LoginStaffRequest,
// ) (*authv1.LoginResponse, error) {
// 	token,err:= s.auth.LoginStaff(ctx, in)
// 	return &authv1.LoginResponse{
// 		AccessToken:  "",
// 		ExpiresIn:    0,
// 		RefreshToken: "",
// 		User:         &authv1.User{},
// 		Status:       "",
// 		Error:        "",
// 	}, err
// }

// func (s *serverAPI) LoginYandex(
// 	ctx context.Context,
// 	in *authv1.OAuthYandexRequest,
// ) (*authv1.LoginResponse, error) {
// 	return s.auth.LoginYandex(ctx, in)
// }
//
// func (s *serverAPI) Validate(
// 	ctx context.Context,
// 	in *authv1.ValidateRequest,
// ) (*authv1.ValidateResponse, error) {
// 	return s.auth.Validate(ctx, in)
// }
//
// func (s *serverAPI) Refresh(
// 	ctx context.Context,
// 	in *authv1.RefreshRequest,
// ) (*authv1.LoginResponse, error) {
// 	return s.auth.Refresh(ctx, in)
// }
