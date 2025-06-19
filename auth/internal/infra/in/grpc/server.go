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
	) (accessToken string,
		aTokenExpire int64,
		refreshToken string,
		rTokenExpire int64,
		client *models.Client,
		err error)
	LoginStaff(
		ctx context.Context,
		email string,
		password string,
	) (accessToken string,
		aTokenExpire int64,
		refreshToken string,
		rTokenExpire int64,
		staff *models.Staff,
		err error)
}

type Token interface {
	Verify(ctx context.Context, token string) (bool, string, string, string, error)
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

type User interface {
	RegisterStaff(
		ctx context.Context,
		email, position string,
	) (staff *models.Staff, err error)
	DeactivateStaff(ctx context.Context, email string) (staff *models.Staff, err error)
	UpdateStaff(
		ctx context.Context,
		email string,
		newEmail, pswd, position *string,
	) (staff *models.Staff, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth  Auth
	audit Audit
	token Token
	user  User
}

//nolint:exhaustruct
func Register(
	gRPCServer *grpc.Server,
	auth Auth,
	audit Audit,
	token Token,
	user User,
) {
	authv1.RegisterAuthServer(
		gRPCServer,
		&serverAPI{auth: auth, audit: audit, token: token, user: user},
	)
}

func (s *serverAPI) RegisterStaff(
	ctx context.Context,
	in *authv1.RegisterStaffRequest,
) (*authv1.RegisterStaffResponse, error) {
	st, err := s.user.RegisterStaff(ctx, in.GetStaff().GetWorkEmail(), in.GetStaff().GetPosition())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "")
	}

	// Audit log for registration
	_ = s.audit.LogEvent(
		ctx,
		st.ID.String(),
		string(models.UserTypeStaff),
		models.ActionRegisterStaff,
		"", // TODO: extract ipAddress from ctx
		"", // TODO: extract userAgent from ctx
	)

	return &authv1.RegisterStaffResponse{
		Staff: &authv1.Staff{
			WorkEmail: st.WorkEmail,
			Position:  st.Position,
		},
	}, nil
}

//nolint:protogetter
func (s *serverAPI) UpdateStaff(
	ctx context.Context,
	in *authv1.UpdateStaffRequest,
) (*authv1.UpdateStaffResponse, error) {
	staff, err := s.user.UpdateStaff(
		ctx,
		in.CurrentEmail,
		in.NewWorkEmail,
		in.NewPassword,
		in.NewPosition,
	)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "")
	}

	return &authv1.UpdateStaffResponse{
		Staff: &authv1.Staff{
			WorkEmail: staff.WorkEmail,
			Position:  staff.Position,
		},
	}, nil
}

func (s *serverAPI) DeactivateStaff(
	ctx context.Context,
	in *authv1.DeactivateStaffRequest,
) (*authv1.DeactivateStaffResponse, error) {
	staff, err := s.user.DeactivateStaff(ctx, in.GetWorkEmail())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "")
	}

	return &authv1.DeactivateStaffResponse{
		WorkEmail: staff.WorkEmail,
		Position:  staff.Position,
		Status:    staff.IsActiveString(),
	}, nil
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

	// Audit log for client login
	_ = s.audit.LogEvent(
		ctx,
		user.ID.String(),
		string(models.UserTypeClient),
		models.ActionLoginClient,
		"", // TODO: extract ipAddress from ctx
		"", // TODO: extract userAgent from ctx
	)

	return &authv1.LoginResponse{
		AccessToken:           tkn,
		ExpiresIn:             tknExp,
		RefreshToken:          rtkn,
		RefreshTokenExpiresIn: rtknExp,
		User: &authv1.User{
			Id: user.ID.String(),
			UserType: &authv1.User_Client{
				Client: &authv1.Client{
					Email: "",
					Phone: user.Phone,
				},
			},
			Roles: []authv1.Role{authv1.Role_ROLE_CLIENT},
		},
	}, nil
}

func (s *serverAPI) Validate(
	ctx context.Context,
	in *authv1.ValidateRequest,
) (*authv1.ValidateResponse, error) {
	valid, id, uphone, utype, err := s.token.Verify(ctx, in.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	var resp *authv1.ValidateResponse

	switch utype {
	case string(models.UserTypeClient):
		resp = &authv1.ValidateResponse{
			Valid: valid,
			User: &authv1.User{
				Id: id,
				UserType: &authv1.User_Client{
					Client: &authv1.Client{
						Email: "",
						Phone: uphone,
					},
				},
				Roles: []authv1.Role{authv1.Role_ROLE_CLIENT},
			},
		}
	case string(models.UserTypeStaff):
		resp = &authv1.ValidateResponse{
			Valid: valid,
			User: &authv1.User{
				Id: id,
				UserType: &authv1.User_Staff{
					Staff: &authv1.Staff{
						WorkEmail: uphone,
						// TODO: FIX
						Position: "",
					},
				},
				Roles: []authv1.Role{authv1.Role_ROLE_STAFF},
			},
		}
	}

	return resp, nil
}

func (s *serverAPI) LoginStaff(
	ctx context.Context,
	in *authv1.LoginStaffRequest,
) (*authv1.LoginResponse, error) {
	tkn, tknExp, rtkn, rtknExp, user, err := s.auth.LoginStaff(
		ctx,
		in.GetStaff().GetWorkEmail(),
		in.GetPassword(),
	)

	switch {
	case errors.Is(err, domain.ErrInternal):
		return nil, status.Error(codes.Internal, err.Error())
	case errors.Is(err, domain.ErrUserDeactivated):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	default:
	}

	// Audit log for staff login
	_ = s.audit.LogEvent(
		ctx,
		user.ID.String(),
		string(models.UserTypeStaff),
		models.ActionLoginStaff,
		"", // TODO: extract ipAddress from ctx
		"", // TODO: extract userAgent from ctx
	)

	return &authv1.LoginResponse{
		AccessToken:           tkn,
		ExpiresIn:             tknExp,
		RefreshToken:          rtkn,
		RefreshTokenExpiresIn: rtknExp,
		User: &authv1.User{
			Id: user.ID.String(),
			UserType: &authv1.User_Staff{
				Staff: &authv1.Staff{
					WorkEmail: user.WorkEmail,
					Position:  user.Position,
				},
			},
			Roles: []authv1.Role{authv1.Role_ROLE_STAFF},
		},
	}, err
}

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
