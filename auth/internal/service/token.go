package service

import (
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	pb "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
)

type TokenService struct {
	tokenRepo repository.TokenRepository
	// jwtManager *utils.JWTManager
}

func NewTokenService(
	tokenRepo repository.TokenRepository,
	// jwtManager *utils.JWTManager,
) *TokenService {
	return &TokenService{
		tokenRepo: tokenRepo,
		// jwtManager: jwtManager,
	}
}

func (s *TokenService) GenerateTokens(user *pb.User) (accessToken, refreshToken string, err error) {
	panic("implement me")
}

func (s *TokenService) generateAccessToken(
	userID, userType string,
	roles []pb.Role,
) (string, error) {
	panic("implement me")
}

func (s *TokenService) ValidateToken(tokenString string) (*models.Claims, error) {
	panic("implement me")
}
