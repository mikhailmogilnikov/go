package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authv1 "github.com/mikhailmogilnikov/go/final/gateway/internal/pb/auth/v1"
)

// AuthHandler хендлер для авторизации
type AuthHandler struct {
	authClient authv1.AuthServiceClient
}

// NewAuthHandler создаёт новый хендлер
func NewAuthHandler(authClient authv1.AuthServiceClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

// RegisterRequest запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest запрос на вход
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse ответ авторизации
type AuthResponse struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// Register регистрация пользователя
// @Summary Регистрация нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authClient.Register(c.Request.Context(), &authv1.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		// Проверяем на конфликт (пользователь уже существует)
		if containsString(err.Error(), "AlreadyExists") {
			c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		UserID: resp.GetUserId(),
		Token:  resp.GetToken(),
	})
}

// Login вход пользователя
// @Summary Вход в систему
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authClient.Login(c.Request.Context(), &authv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if containsString(err.Error(), "Unauthenticated") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		UserID: resp.GetUserId(),
		Token:  resp.GetToken(),
	})
}

// Register регистрирует роуты
func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}



