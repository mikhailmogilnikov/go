package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mikhailmogilnikov/go/final/auth/internal/domain"
)

// UserRepository реализует domain.UserRepository для PostgreSQL
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создаёт новый репозиторий
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create создаёт нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		// Проверяем на дубликат email
		if isDuplicateKeyError(err) {
			return errors.New("user with this email already exists")
		}
		return err
	}
	return nil
}

// GetByEmail ищет пользователя по email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetByID ищет пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// isDuplicateKeyError проверяет, является ли ошибка дубликатом ключа
func isDuplicateKeyError(err error) bool {
	// PostgreSQL код ошибки 23505 - unique_violation
	return err != nil && contains(err.Error(), "23505")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}



