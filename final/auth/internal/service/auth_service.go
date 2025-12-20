package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/mikhailmogilnikov/go/final/auth/internal/domain"
)

// AuthService сервис авторизации
type AuthService struct {
	userRepo  domain.UserRepository
	jwtSecret []byte
	tokenTTL  time.Duration
}

// NewAuthService создаёт новый сервис авторизации
func NewAuthService(userRepo domain.UserRepository, jwtSecret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		tokenTTL:  tokenTTL,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, email, password string) (int64, string, error) {
	// Валидация
	if err := domain.ValidateEmail(email); err != nil {
		return 0, "", err
	}
	if err := domain.ValidatePassword(password); err != nil {
		return 0, "", err
	}

	// Проверяем, что пользователь не существует
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return 0, "", err
	}
	if existing != nil {
		return 0, "", errors.New("user with this email already exists")
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, "", err
	}

	// Создаём пользователя
	user := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return 0, "", err
	}

	// Генерируем токен
	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}

// Login выполняет вход пользователя
func (s *AuthService) Login(ctx context.Context, email, password string) (int64, string, error) {
	// Ищем пользователя
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return 0, "", err
	}
	if user == nil {
		return 0, "", errors.New("invalid email or password")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return 0, "", errors.New("invalid email or password")
	}

	// Генерируем токен
	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}

// ValidateToken проверяет JWT токен
func (s *AuthService) ValidateToken(tokenString string) (int64, string, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, "", false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", false
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", false
	}

	email, ok := claims["email"].(string)
	if !ok {
		return 0, "", false
	}

	return int64(userID), email, true
}

// generateToken генерирует JWT токен
func (s *AuthService) generateToken(userID int64, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(s.tokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}



