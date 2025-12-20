package domain

import "context"

// UserRepository интерфейс для работы с пользователями в БД
type UserRepository interface {
	// Create создаёт нового пользователя
	Create(ctx context.Context, user *User) error
	// GetByEmail ищет пользователя по email
	GetByEmail(ctx context.Context, email string) (*User, error)
	// GetByID ищет пользователя по ID
	GetByID(ctx context.Context, id int64) (*User, error)
}



