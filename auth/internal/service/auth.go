package service

import (
	"context"

	"github.com/netscrawler/Restaurant_is/auth/internal/repository"
	pb "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	clientRepo repository.ClientRepository
	staffRepo  repository.StaffRepository
	tokenRepo  repository.TokenRepository
	oauthRepo  repository.OAuthRepository
	// jwtManager  *utils.JWTManager
	// oauthYandex *utils.YandexOAuth
}

func NewAuthService(
	clientRepo repository.ClientRepository,
	staffRepo repository.StaffRepository,
	tokenRepo repository.TokenRepository,
	oauthRepo repository.OAuthRepository,
	// jwtManager *utils.JWTManager,
	// oauthYandex *utils.YandexOAuth,
) *AuthService {
	return &AuthService{
		clientRepo: clientRepo,
		staffRepo:  staffRepo,
		tokenRepo:  tokenRepo,
		oauthRepo:  oauthRepo,
		// jwtManager:  jwtManager,
		// oauthYandex: oauthYandex,
	}
}

func (s *AuthService) LoginClient(
	ctx context.Context,
	req *pb.LoginClientRequest,
) (*pb.LoginResponse, error) {
	// Реализация
	// 1. Поиск клиента
	// 2. Проверка пароля
	// 3. Генерация токенов
	// 4. Запись в аудит
	panic("implement me")
}

func (s *AuthService) LoginStaff(
	ctx context.Context,
	req *pb.LoginStaffRequest,
) (*pb.LoginResponse, error) {
	// Реализация
	// Аналогично клиенту, но для сотрудников

	panic("implement me")
}

func (s *AuthService) LoginYandex(
	ctx context.Context,
	req *pb.OAuthYandexRequest,
) (*pb.LoginResponse, error) {
	// Реализация
	// 1. Обмен кода на токен
	// 2. Получение данных пользователя
	// 3. Создание/обновление локальной записи
	// 4. Генерация JWT
	panic("implement me")
}

func (s *AuthService) Validate(
	ctx context.Context,
	req *pb.ValidateRequest,
) (*pb.ValidateResponse, error) {
	// Реализация
	panic("implement me")
}

func (s *AuthService) Refresh(
	ctx context.Context,
	req *pb.RefreshRequest,
) (*pb.LoginResponse, error) {
	// Реализация
	panic("implement me")
}
