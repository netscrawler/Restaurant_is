package authgrpc

import (
	authv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
	"google.golang.org/grpc"
)

// type AuthServiceServer interface {
// 	LoginClient(context.Context, *authv1.LoginClientRequest) (*authv1.LoginResponse, error)
// 	LoginStaff(context.Context, *authv1.LoginStaffRequest) (*authv1.LoginResponse, error)
// 	LoginYandex(context.Context, *authv1.OAuthYandexRequest) (*authv1.LoginResponse, error)
// 	Validate(context.Context, *authv1.ValidateRequest) (*authv1.ValidateResponse, error)
// 	Refresh(context.Context, *authv1.RefreshRequest) (*authv1.LoginResponse, error)
// 	// mustEmbedUnimplementedAuthServiceServer()
// }

type Auth interface{}

type serverAPI struct {
	authv1.UnimplementedAuthServiceServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	authv1.RegisterAuthServiceServer(gRPCServer, &serverAPI{auth: auth})
}

// func (s *serverAPI) LoginClient(
// 	ctx context.Context,
// 	in *authv1.LoginClientRequest,
// ) (*authv1.LoginResponse, error) {
// 	panic("not implemented") // TODO: Implement
// }
//
// func (s *serverAPI) LoginStaff(
// 	ctx context.Context,
// 	in *authv1.LoginStaffRequest,
// ) (*authv1.LoginResponse, error) {
// 	panic("not implemented") // TODO: Implement
// }
//
// func (s *serverAPI) LoginYandex(
// 	ctx context.Context,
// 	in *authv1.OAuthYandexRequest,
// ) (*authv1.LoginResponse, error) {
// 	panic("not implemented") // TODO: Implement
// }
//
// func (s *serverAPI) Validate(
// 	ctx context.Context,
// 	in *authv1.ValidateRequest,
// ) (*authv1.ValidateResponse, error) {
// 	panic("not implemented") // TODO: Implement
// }
//
// func (s *serverAPI) Refresh(
// 	ctx context.Context,
// 	in *authv1.RefreshRequest,
// ) (*authv1.LoginResponse, error) {
// 	panic("not implemented") // TODO: Implement
// }
