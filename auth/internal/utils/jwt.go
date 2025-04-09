package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
)

type JWTManager struct {
	secretKey            string
	refreshSecretKey     string
	tokenDuration        time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

func NewJWTManager(
	secret string,
	dur time.Duration,
	refreshSecret string,
	refreshDur time.Duration,
	issuer string,
) (*JWTManager, error) {
	// Проверка на пустые ключи
	if secret == "" || refreshSecret == "" {
		return nil, fmt.Errorf("%w (%w)", domain.ErrSecret, domain.ErrEmptySecret)
	}

	return &JWTManager{
		secretKey:            secret,
		refreshSecretKey:     refreshSecret,
		tokenDuration:        dur,
		refreshTokenDuration: refreshDur,
		issuer:               issuer,
	}, nil
}

// Claims структура расширена для соответствия proto-определению User.
type Claims struct {
	UserID   string   `json:"userId"`
	UserType string   `json:"userType"` // "client" или "staff"
	Roles    []string `json:"roles"`    // Роли пользователя из enum Role
	jwt.RegisteredClaims
}

// RefreshClaims содержит минимальную информацию для refresh токена.
type RefreshClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// Generate создает access token с полными данными пользователя.
func (m *JWTManager) Generate(userID string, userType string, roles []string) (string, error) {
	if userID == "" {
		return "", domain.ErrInvalidUserUUID
	}

	now := time.Now()
	claims := Claims{
		UserID:   userID,
		UserType: userType,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", fmt.Errorf("%w %w", domain.ErrInternalCodeGen, err)
	}

	return signedToken, nil
}

// GetRefreshTokenDuration возвращает настроенную длительность refresh токена.
func (m *JWTManager) GetRefreshTokenDuration() time.Duration {
	return m.refreshTokenDuration
}

// GenerateRefreshToken создает refresh токен с увеличенным сроком действия.
func (m *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	if userID == "" {
		return "", domain.ErrInvalidUserUUID
	}

	now := time.Now()
	refreshClaims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedRefreshToken, err := refreshToken.SignedString([]byte(m.refreshSecretKey))
	if err != nil {
		return "", fmt.Errorf("%w %w", domain.ErrInternalCodeGen, err)
	}

	return signedRefreshToken, nil
}

// GenerateTokenPair создает пару токенов (access + refresh) за один вызов.
func (m *JWTManager) GenerateTokenPair(
	userID string,
	userType string,
	roles []string,
) (string, string, error) {
	accessToken, err := m.Generate(userID, userType, roles)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := m.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Verify проверяет access token и возвращает claims.
func (m *JWTManager) Verify(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, domain.ErrInvalidToken
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnexpectedSignMethod
		}

		return []byte(m.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternalCodeParse, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	// Проверка issuer
	if claims.Issuer != m.issuer {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

// VerifyRefreshToken проверяет refresh токен.
func (m *JWTManager) VerifyRefreshToken(refreshTokenStr string) (*RefreshClaims, error) {
	if refreshTokenStr == "" {
		return nil, domain.ErrInvalidToken
	}

	refreshToken, err := jwt.Parse(
		refreshTokenStr,
		func(token *jwt.Token) (any, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.ErrUnexpectedSignMethod
			}

			return []byte(m.refreshSecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w (%w)", domain.ErrInternalCodeParse, err)
	}

	refreshClaims, ok := refreshToken.Claims.(*RefreshClaims)
	if !ok || !refreshToken.Valid {
		return nil, domain.ErrInvalidToken
	}

	// Проверка issuer
	if refreshClaims.Issuer != m.issuer {
		return nil, domain.ErrInvalidToken
	}

	return refreshClaims, nil
}

// RefreshAccessToken создает новый access token на основе действительного refresh токена.
func (m *JWTManager) RefreshAccessToken(
	refreshTokenStr string,
	userType string,
	roles []string,
) (string, error) {
	// Проверяем refresh токен
	refreshClaims, err := m.VerifyRefreshToken(refreshTokenStr)
	if err != nil {
		return "", err
	}

	// Создаем новый access token
	return m.Generate(refreshClaims.UserID, userType, roles)
}
