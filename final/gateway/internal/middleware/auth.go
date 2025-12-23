package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authv1 "github.com/mikhailmogilnikov/go/final/gateway/internal/pb/auth/v1"
)

type AuthMiddleware struct {
	authClient authv1.AuthServiceClient
}

func NewAuthMiddleware(authClient authv1.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{authClient: authClient}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		token := parts[1]

		resp, err := m.authClient.ValidateToken(c.Request.Context(), &authv1.ValidateTokenRequest{
			Token: token,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to validate token"})
			return
		}

		if !resp.GetValid() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", resp.GetUserId())
		c.Set("email", resp.GetEmail())

		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get("user_id")
	if id, ok := userID.(int64); ok {
		return id
	}
	return 0
}



