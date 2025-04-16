package authgrpc

import (
	"context"
	"errors"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	authv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//	type Auth interface {
//		LoginClientInit(
//			ctx context.Context,
//			req *authv1.LoginClientRequest,
//		) (*authv1.LoginInitResponse, error)
//
//		LoginClientConfirm(
//			context.Context,
//			*authv1.LoginClientConfirmRequest,
//		) (*authv1.LoginResponse, error)
//		LoginStaff(ctx context.Context, req *authv1.LoginStaffRequest) (*authv1.LoginResponse, error)
//		LoginYandex(ctx context.Context, req *authv1.OAuthYandexRequest) (*authv1.LoginResponse, error)
//		Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error)
//		Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.LoginResponse, error)
//	}
type Auth interface {
	LoginClientInit(ctx context.Context, phone string) error
	LoginClientConfirm(
		ctx context.Context,
		phone string,
		code string,
	) (accessToken string, aTokenExpire int64, refreshToken string, rTokenExpire int64, client *models.Client, err error)
	Verify(ctx context.Context, token string) (bool, string, string, error)
}

type Audit interface {
	LogEvent(
		ctx context.Context,
		userID string,
		userType string,
		eventAction models.AuthEventAction,
		ipAddress, userAgent string,
	) error
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth  Auth
	audit Audit
}

func Register(
	gRPCServer *grpc.Server,
	auth Auth,
	audit Audit,
) {
	authv1.RegisterAuthServer(
		gRPCServer,
		&serverAPI{auth: auth, audit: audit},
	)
}

func (s *serverAPI) LoginClientInit(
	ctx context.Context,
	in *authv1.LoginClientInitRequest,
) (*authv1.LoginInitResponse, error) {
	err := s.auth.LoginClientInit(ctx, in.GetPhone())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &authv1.LoginInitResponse{
		Status: "ok",
		Error:  "",
	}, nil
}

func (s *serverAPI) LoginClientConfirm(
	ctx context.Context,
	in *authv1.LoginClientConfirmRequest,
) (*authv1.LoginResponse, error) {
	tkn, tknExp, rtkn, rtknExp, user, err := s.auth.LoginClientConfirm(
		ctx,
		in.GetPhone(),
		in.GetCode(),
	)

	switch {
	case errors.Is(err, domain.ErrInternal):
		return nil, status.Error(codes.Internal, err.Error())
	case errors.Is(err, domain.ErrInvalidCode):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	default:
	}

	return &authv1.LoginResponse{
		AccessToken:      tkn,
		ExpiresIn:        tknExp,
		RefreshToken:     rtkn,
		RefreshExpiresIn: rtknExp,
		User: &authv1.User{
			Id: user.ID.String(),
			UserType: &authv1.User_Client{
				Client: &authv1.Client{
					Email:    "",
					Phone:    user.Phone,
					FullName: "",
				},
			},
		},
		Status: "Success",
		Error:  "",
	}, nil
}

func (s *serverAPI) Validate(
	ctx context.Context,
	in *authv1.ValidateRequest,
) (*authv1.ValidateResponse, error) {
	valid, id, uphone, err := s.auth.Verify(ctx, in.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &authv1.ValidateResponse{
		Valid: valid,
		User: &authv1.User{
			Id: id,
			UserType: &authv1.User_Client{
				Client: &authv1.Client{
					Email:    "",
					Phone:    uphone,
					FullName: "",
				},
			},
		},
	}, nil
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
//
// func (s *serverAPI) Refresh(
// 	ctx context.Context,
// 	in *authv1.RefreshRequest,
// ) (*authv1.LoginResponse, error) {
// 	return s.auth.Refresh(ctx, in)
// }
