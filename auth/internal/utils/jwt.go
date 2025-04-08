package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
	issuer        string
}

func NewJWTManager(secret string, dur time.Duration, issuer string) *JWTManager {
	return &JWTManager{
		secretKey:     secret,
		tokenDuration: dur,
		issuer:        issuer,
	}
}

type Claims struct {
	UserID   string `json:"userId"`
	UserType string `json:"userType"` // "client" или "staff"
	jwt.RegisteredClaims
}

func (m *JWTManager) Generate(userID string, userType string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		UserType: userType,
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

func (m *JWTManager) Verify(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.ErrUnexpectedSignMethod
			}

			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w %w", domain.ErrInternalCodeParse, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}
