package service

import (
	"context"
	"errors"

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
	tokenRepo    repository.TokenRepository
	oauthRepo    repository.OAuthRepository
	notify       NotifySender
	codeProvider CodeProvider
	jwtManager   *utils.JWTManager
	// oauthYandex *utils.YandexOAuth
	// codeCache *codeCache
}

func NewAuthService(
	log *zap.Logger,
	clientRepo repository.ClientRepository,
	staffRepo repository.StaffRepository,
	tokenRepo repository.TokenRepository,
	oauthRepo repository.OAuthRepository,
	notifySender NotifySender,
	codeProvider CodeProvider,
	jwtManager *utils.JWTManager,
	// oauthYandex *utils.YandexOAuth,
) *AuthService {
	return &AuthService{
		clientRepo: clientRepo,
		staffRepo:  staffRepo,
		tokenRepo:  tokenRepo,
		oauthRepo:  oauthRepo,
		jwtManager: jwtManager,
		// oauthYandex: oauthYandex,
		notify:       notifySender,
		log:          log,
		codeProvider: codeProvider,
	}
}

func (a *AuthService) LoginClinetInit(ctx context.Context, phone string) error {
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

	code, err := utils.GenerateSecureCode()
	if err != nil {
		return domain.ErrFailedCreateCode
	}

	a.codeProvider.Set(user.ID, code)
	// Todo: add token cache.
	go a.notify.Send(ctx, phone, models.NewCodeMsg(code).String())

	return nil
}

func (a *AuthService) LoginClientConfirm(ctx context.Context, phone string) (string, error) {
	return nil, nil
}
