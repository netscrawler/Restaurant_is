package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	"github.com/netscrawler/Restaurant_is/auth/internal/utils"
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
	log          *slog.Logger
	clientRepo   repository.ClientRepository
	staffRepo    repository.StaffRepository
	notify       NotifySender
	codeProvider CodeProvider
	jwtManager   *utils.JWTManager
}

func NewAuthService(
	log *slog.Logger,
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

	log := utils.LoggerWithTrace(ctx, a.log)

	user, err := a.clientRepo.GetClientByPhone(ctx, phone)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		log.Info(op+"not found client, creating new", slog.String("phone", phone))
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
		log := utils.LoggerWithTrace(ctx, a.log)
		code, err := utils.GenerateSecureCode()
		if err != nil {
			log.Error("error generate secure token", slog.Any("error", err))
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

	log := utils.LoggerWithTrace(ctx, a.log)

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
		log.Info(op+" Invalid code", slog.Any("user", user), slog.Int("code", codeInt))

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

func (a *AuthService) LoginStaff(
	ctx context.Context,
	email string,
	password string,
) (string, int64, string, int64, *models.Staff, error) {
	const op = "service.Auth.LoginStaff"
	log := utils.LoggerWithTrace(ctx, a.log).
		With(slog.String("fn", op), slog.String("email", email))

	user, err := a.staffRepo.GetStaffByEmail(ctx, email)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		log.Info("user not found")

		return "", 0, "", 0, nil, domain.ErrNotFound
	case err != nil:
		log.Info("internal error", slog.Any("error", err))

		return "", 0, "", 0, nil, domain.ErrInternal
	default:
	}

	if !user.IsActive {
		return "", 0, "", 0, nil, domain.ErrUserDeactivated
	}

	accessToken, aTokenExpire, refreshToken, rTokenExpire, err := a.jwtManager.GenerateTokenPair(
		user.ID.String(),
		string(models.UserTypeStaff),
		user.WorkEmail,
	)
	if err != nil {
		log.Info("error generating token pair", slog.Any("error", err))

		return "", 0, "", 0, nil, domain.ErrInternal
	}

	msg := "Your account has been logged in"

	a.notify.Send(ctx, email, fmt.Sprintf("%s %v", msg, time.Now().Format(time.DateTime)))

	return accessToken, aTokenExpire, refreshToken, rTokenExpire, user, nil
}
