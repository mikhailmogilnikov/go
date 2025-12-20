package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mikhailmogilnikov/go/final/auth/internal/pb/auth/v1"
	"github.com/mikhailmogilnikov/go/final/auth/internal/service"
)

// AuthServer реализует gRPC сервер авторизации
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

// NewAuthServer создаёт новый gRPC сервер
func NewAuthServer(authService *service.AuthService) *AuthServer {
	return &AuthServer{authService: authService}
}

// Register регистрирует нового пользователя
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, token, err := s.authService.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// Если пользователь уже существует
		if err.Error() == "user with this email already exists" {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		// Ошибки валидации
		if err.Error() == "invalid email format" || err.Error() == "email cannot be empty" ||
			err.Error() == "password must be at least 6 characters" {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}

	return &pb.RegisterResponse{
		UserId: userID,
		Token:  token,
	}, nil
}

// Login выполняет вход пользователя
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, token, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if err.Error() == "invalid email or password" {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return &pb.LoginResponse{
		UserId: userID,
		Token:  token,
	}, nil
}

// ValidateToken проверяет JWT токен
func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.GetToken() == "" {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	userID, email, valid := s.authService.ValidateToken(req.GetToken())
	return &pb.ValidateTokenResponse{
		Valid:  valid,
		UserId: userID,
		Email:  email,
	}, nil
}



