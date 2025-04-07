package service

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	pb "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
)

type OAuthService struct {
	oauthRepo  repository.OAuthRepository
	clientRepo repository.ClientRepository
	// jwtManager *utils.JWTManager
}

func (s *OAuthService) HandleCallback(
	ctx context.Context,
	provider string,
	code string,
) (*pb.LoginResponse, error) {
	// Обработка callback от провайдера
	panic("implement me")
}

func (s *OAuthService) GetAuthURL(provider string) (string, error) {
	// Генерация URL для перенаправления
	panic("implement me")
}
