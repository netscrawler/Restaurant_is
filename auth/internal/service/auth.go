package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	"github.com/netscrawler/Restaurant_is/auth/internal/utils"
	"go.uber.org/zap"
)

type NotifySender interface {
	Send(ctx context.Context, to string, msg string)
}
type CodeProvider interface {
	Set(user uuid.UUID, code int)
	Get(user uuid.UUID) (int, bool)
	Delete(user uuid.UUID)
}

type AuthService struct {
	log          *zap.Logger
	clientRepo   repository.ClientRepository
	staffRepo    repository.StaffRepository
	notify       NotifySender
	codeProvider CodeProvider
	jwtManager   *utils.JWTManager
}

func NewAuthService(
	log *zap.Logger,
	clientRepo repository.ClientRepository,
	staffRepo repository.StaffRepository,
	notifySender NotifySender,
	codeProvider CodeProvider,
	jwtManager *utils.JWTManager,
) *AuthService {
	return &AuthService{
		clientRepo:   clientRepo,
		staffRepo:    staffRepo,
		jwtManager:   jwtManager,
		notify:       notifySender,
		log:          log,
		codeProvider: codeProvider,
	}
}

func (a *AuthService) LoginClientInit(ctx context.Context, phone string) error {
	const op = "service.Auth.LoginInit"

	user, err := a.clientRepo.GetClientByPhone(ctx, phone)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		a.log.Info(op+"not found client, creating new", zap.String("phone", phone))
		user = models.NewClient(phone)

		err = a.clientRepo.CreateClient(ctx, user)
		if err != nil {
			return domain.ErrInternal
		}
	case err != nil:
		return domain.ErrInternal
	default:
	}

	a.genAndSend(ctx, user)

	return nil
}

func (a *AuthService) genAndSend(ctx context.Context, u *models.Client) {
	go func(ctx context.Context, u *models.Client) {
		code, err := utils.GenerateSecureCode()
		if err != nil {
			a.log.Error("error generate secure token", zap.Error(err))
		}

		a.codeProvider.Set(u.ID, code)
		a.notify.Send(ctx, u.Phone, models.NewCodeMsg(code).String())
	}(ctx, u)
}

func (a *AuthService) LoginClientConfirm(
	ctx context.Context,
	phone string,
	code string,
) (string, int64, string, int64, *models.Client, error) {
	const op = "service.Auth.LoginClientConfirm"

	codeInt, err := strconv.Atoi(code)
	if err != nil {
		return "", 0, "", 0, nil, domain.ErrInvalidCode
	}

	user, err := a.clientRepo.GetClientByPhone(ctx, phone)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return "", 0, "", 0, nil, domain.ErrNotFound
	case err != nil:
		return "", 0, "", 0, nil, domain.ErrInternal
	default:
	}

	storedCode, exists := a.codeProvider.Get(user.ID)

	if !exists || storedCode != codeInt {
		a.log.Info(op+" Invalid code", zap.Any("user", user), zap.Int("code", codeInt))

		return "", 0, "", 0, nil, domain.ErrInvalidCode
	}

	accessToken, aTokenExpire, refreshToken, rTokenExpire, err := a.jwtManager.GenerateTokenPair(
		user.ID.String(),
		string(models.UserTypeClient),
		user.Phone,
	)
	if err != nil {
		return "", 0, "", 0, nil, domain.ErrInternal
	}

	return accessToken, aTokenExpire, refreshToken, rTokenExpire, user, nil
}

func (a *AuthService) Verify(ctx context.Context, token string) (bool, string, string, error) {
	cl, err := a.jwtManager.VerifyAccessToken(token)
	if err != nil {
		return false, "", "", fmt.Errorf("%w (%w)", domain.ErrInvalidCode, err)
	}

	return true, cl.UserID, cl.UserPhone, nil
}
