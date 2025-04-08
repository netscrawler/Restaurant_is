package service

import (
	"context"
	"fmt"
	"time"

	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	"github.com/netscrawler/Restaurant_is/auth/internal/utils"
	pb "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
)

type TokenService struct {
	tokenRepo  repository.TokenRepository
	jwtManager *utils.JWTManager
}

func NewTokenService(
	tokenRepo repository.TokenRepository,
	jwtManager *utils.JWTManager,
) *TokenService {
	return &TokenService{
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}
}

// GenerateTokens генерирует пару токенов (access и refresh) для пользователя
func (s *TokenService) GenerateTokens(
	ctx context.Context,
	user *pb.User,
	ipAddress, userAgent string,
) (accessToken, refreshToken string, err error) {
	if user == nil || user.Id == "" {
		return "", "", fmt.Errorf("invalid user")
	}

	// Определяем тип пользователя и роли
	userType := ""
	var protoRoles []pb.Role

	switch u := user.UserType.(type) {
	case *pb.User_Client:
		userType = "client"
		protoRoles = []pb.Role{pb.Role_ROLE_CLIENT}
	case *pb.User_Staff:
		userType = "staff"
		if len(user.Roles) > 0 {
			protoRoles = user.Roles
		} else {
			protoRoles = []pb.Role{pb.Role_ROLE_ADMIN}
		}
	default:
		return "", "", fmt.Errorf("unknown user type: %T", u)
	}

	// Генерируем access token
	accessToken, err = s.generateAccessToken(user.Id, userType, protoRoles)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Генерируем refresh token
	refreshToken, err = s.jwtManager.GenerateRefreshToken(user.Id)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Сохраняем refresh token в репозиторий
	err = s.tokenRepo.StoreRefreshToken(ctx, &models.RefreshToken{
		UserID:       user.Id,
		UserType:     userType,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.jwtManager.GetRefreshTokenDuration()),
		Revoked:      false,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Логируем действие авторизации
	_ = s.tokenRepo.LogTokenAction(ctx, user.Id, userType, "login", ipAddress, userAgent)

	return accessToken, refreshToken, nil
}

// generateAccessToken создает JWT access token с информацией о пользователе
func (s *TokenService) generateAccessToken(
	userID, userType string,
	roles []pb.Role,
) (string, error) {
	// Преобразуем роли из proto enum в строки для JWT
	roleStrings := make([]string, len(roles))
	for i, role := range roles {
		roleStrings[i] = role.String()
	}

	// Генерируем токен через JWT Manager
	return s.jwtManager.Generate(userID, userType, roleStrings)
}

// ValidateToken проверяет и декодирует токен
func (s *TokenService) ValidateToken(tokenString string) (*utils.Claims, error) {
	// Проверка токена с помощью JWT Manager
	claims, err := s.jwtManager.Verify(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

// RefreshTokens обновляет пару токенов используя refresh token
func (s *TokenService) RefreshTokens(
	ctx context.Context,
	refreshToken, ipAddress, userAgent string,
) (string, string, error) {
	// Проверяем refresh token в JWT
	refreshClaims, err := s.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Проверяем наличие токена в базе данных
	storedToken, err := s.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("refresh token not found: %w", err)
	}

	// Проверяем, не истёк ли токен
	if time.Now().After(storedToken.ExpiresAt) || storedToken.Revoked {
		// Отзываем токен если он еще действителен
		_ = s.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
		return "", "", fmt.Errorf("refresh token expired or revoked")
	}

	// Получаем информацию о пользователе
	userID := refreshClaims.UserID
	userType := storedToken.UserType

	// Теоретически здесь нужно получить пользователя из базы данных,
	// чтобы иметь актуальную информацию о типе и ролях
	// Но для примера предположим, что у нас есть метод получения пользователя:
	user, err := s.getUserByID(ctx, userID, userType)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %w", err)
	}

	// Отзываем старый refresh token
	err = s.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	// Логируем действие обновления токена
	_ = s.tokenRepo.LogTokenAction(ctx, userID, userType, "token_refresh", ipAddress, userAgent)

	// Генерируем новую пару токенов
	accessToken, newRefreshToken, err := s.GenerateTokens(ctx, user, ipAddress, userAgent)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// RevokeToken отзывает refresh токен
func (s *TokenService) RevokeToken(
	ctx context.Context,
	refreshToken, ipAddress, userAgent string,
) error {
	// Проверяем существование токена
	storedToken, err := s.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("refresh token not found: %w", err)
	}

	// Отзываем токен
	err = s.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	// Логируем действие отзыва токена
	_ = s.tokenRepo.LogTokenAction(
		ctx,
		storedToken.UserID,
		storedToken.UserType,
		"token_revoke",
		ipAddress,
		userAgent,
	)

	return nil
}

// RevokeAllUserTokens отзывает все токены пользователя (выход со всех устройств)
func (s *TokenService) RevokeAllUserTokens(
	ctx context.Context,
	userID, userType, ipAddress, userAgent string,
) error {
	// Отзываем все токены пользователя
	err := s.tokenRepo.DeleteAllUserTokens(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}

	// Логируем действие
	_ = s.tokenRepo.LogTokenAction(ctx, userID, userType, "logout", ipAddress, userAgent)

	return nil
}

// CleanupExpiredTokens удаляет все истекшие токены
func (s *TokenService) CleanupExpiredTokens(ctx context.Context) (int64, error) {
	return s.tokenRepo.CleanupExpiredTokens(ctx)
}

// Вспомогательный метод для получения пользователя по ID
// В реальном приложении нужно реализовать этот метод или внедрить UserRepository
func (s *TokenService) getUserByID(ctx context.Context, userID, userType string) (*pb.User, error) {
	// Пример-заглушка
	// В реальном коде здесь будет запрос к базе данных или другому сервису
	user := &pb.User{
		Id: userID,
	}

	// Заполняем информацию в зависимости от типа пользователя
	switch userType {
	case "client":
		user.UserType = &pb.User_Client{
			Client: &pb.Client{},
		}
		user.Roles = []pb.Role{pb.Role_ROLE_CLIENT}
	case "staff":
		user.UserType = &pb.User_Staff{
			Staff: &pb.Staff{},
		}
		user.Roles = []pb.Role{pb.Role_ROLE_ADMIN}
	default:
		return nil, fmt.Errorf("unknown user type: %s", userType)
	}

	return user, nil
}
