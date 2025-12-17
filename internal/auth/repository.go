package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type AuthRepository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Register(ctx context.Context, body RegisterRequest) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users(name, email, hashed_password) VALUES ($1, $2, $3)", body.Name, body.Email, body.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (UserWithPassword, error) {
	var user UserWithPassword
	err := r.db.QueryRow(ctx, "SELECT id, name, email, hashed_password, email_verified_at, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL", email).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.EmailVerifiedAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, userID string, token string) error {
	_, err := r.db.Exec(ctx, "INSERT INTO refresh_tokens(user_id, token, expires_at) VALUES($1, $2, $3)", userID, token, time.Now().Add(30*24*time.Hour))
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM refresh_tokens WHERE token = $1", token)
	if err != nil {
		return err
	}
	return nil
}
