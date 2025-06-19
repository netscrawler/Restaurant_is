package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	authv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
)

type AuthMiddleware struct {
	authClient authv1.AuthClient
}

func NewAuthMiddleware(authClient authv1.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (am *AuthMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()

			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()

			return
		}

		// Валидируем токен через gRPC сервис auth
		validateReq := &authv1.ValidateRequest{
			Token: tokenString,
		}

		validateResp, err := am.authClient.Validate(context.Background(), validateReq)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()

			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("user_id", validateResp.GetUser().GetId())
		c.Set("user", validateResp.GetUser())
		c.Set("roles", validateResp.GetUser().GetRoles())
		c.Next()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No roles found"})
			c.Abort()

			return
		}

		userRoles, ok := roles.([]authv1.Role)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid roles format"})
			c.Abort()

			return
		}

		// Проверяем, есть ли у пользователя требуемая роль
		for _, role := range userRoles {
			if role.String() == requiredRole {
				c.Next()

				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}
