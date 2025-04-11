package utils

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/netscrawler/Restaurant_is/auth/internal/config"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
)

// JWTManager реализует генерацию и валидацию Access и Refresh токенов на RS256.
type JWTManager struct {
	privateKey           *rsa.PrivateKey
	publicKey            *rsa.PublicKey
	refreshPrivateKey    *rsa.PrivateKey
	refreshPublicKey     *rsa.PublicKey
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

func NewJWTManager(cfg config.JWTConfig) (*JWTManager, error) {
	if cfg.PrivateKey == nil || cfg.PublicKey == nil || cfg.RefreshPrivateKey == nil ||
		cfg.RefreshPublicKey == nil {
		return nil, fmt.Errorf("%w: RSA ключи не заданы", domain.ErrSecret)
	}

	return &JWTManager{
		privateKey:           cfg.PrivateKey,
		publicKey:            cfg.PublicKey,
		refreshPrivateKey:    cfg.RefreshPrivateKey,
		refreshPublicKey:     cfg.RefreshPublicKey,
		accessTokenDuration:  cfg.AccessTTL,
		refreshTokenDuration: cfg.RefreshTTL,
		issuer:               cfg.Issuer,
	}, nil
}

// Claims представляет Access токен с полной информацией о пользователе.
type Claims struct {
	UserID    string `json:"userId"`
	UserType  string `json:"userType"`
	UserPhone string `json:"userPhone"`
	jwt.RegisteredClaims
}

// RefreshClaims представляет Refresh токен.
type RefreshClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func (m *JWTManager) generateToken(claims jwt.Claims, key *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "v1"

	signed, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("%w: %v", domain.ErrInternalCodeGen, err)
	}

	return signed, nil
}

func (m *JWTManager) GenerateAccessToken(userID, userType, userPhone string) (string, error) {
	if userID == "" {
		return "", domain.ErrInvalidUserUUID
	}

	now := time.Now()
	claims := Claims{
		UserID:    userID,
		UserType:  userType,
		UserPhone: userPhone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	return m.generateToken(claims, m.privateKey)
}

func (m *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	if userID == "" {
		return "", domain.ErrInvalidUserUUID
	}

	now := time.Now()
	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	return m.generateToken(claims, m.refreshPrivateKey)
}

func (m *JWTManager) GenerateTokenPair(userID, userType, userPhone string) (string, string, error) {
	accessToken, err := m.GenerateAccessToken(userID, userType, userPhone)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := m.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (m *JWTManager) VerifyAccessToken(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, domain.ErrInvalidToken
	}

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, domain.ErrUnexpectedSignMethod
		}
		return m.publicKey, nil
	}, jwt.WithLeeway(2*time.Second))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInternalCodeParse, err)
	}
	if claims.Issuer != m.issuer {
		return nil, domain.ErrInvalidToken
	}
	return claims, nil
}

func (m *JWTManager) VerifyRefreshToken(tokenStr string) (*RefreshClaims, error) {
	if tokenStr == "" {
		return nil, domain.ErrInvalidToken
	}

	claims := &RefreshClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, domain.ErrUnexpectedSignMethod
		}
		return m.refreshPublicKey, nil
	}, jwt.WithLeeway(2*time.Second))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInternalCodeParse, err)
	}
	if claims.Issuer != m.issuer {
		return nil, domain.ErrInvalidToken
	}
	return claims, nil
}

func (m *JWTManager) RefreshAccessToken(refreshToken, userType, userPhone string) (string, error) {
	claims, err := m.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}
	return m.GenerateAccessToken(claims.UserID, userType, userPhone)
}
