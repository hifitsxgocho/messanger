package postgres

import (
	"context"
	"errors"
	"fmt"
	"messenger/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	q := `INSERT INTO users (id, email, username, password_hash, bio, avatar_url, created_at, updated_at)
	      VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, q, u.ID, u.Email, u.Username, u.PasswordHash, u.Bio, u.AvatarURL, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("user create: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	q := `SELECT id, email, username, password_hash, bio, avatar_url, created_at, updated_at FROM users WHERE id = $1`
	u := &domain.User{}
	err := r.db.QueryRow(ctx, q, id).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("user get by id: %w", err)
	}
	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT id, email, username, password_hash, bio, avatar_url, created_at, updated_at FROM users WHERE email = $1`
	u := &domain.User{}
	err := r.db.QueryRow(ctx, q, email).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Bio, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("user get by email: %w", err)
	}
	return u, nil
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	q := `UPDATE users SET username=$1, bio=$2, avatar_url=$3, updated_at=$4 WHERE id=$5`
	_, err := r.db.Exec(ctx, q, u.Username, u.Bio, u.AvatarURL, u.UpdatedAt, u.ID)
	if err != nil {
		return fmt.Errorf("user update: %w", err)
	}
	return nil
}

func (r *UserRepo) Search(ctx context.Context, query string, limit int) ([]domain.UserPublic, error) {
	q := `SELECT id, username, bio, avatar_url FROM users WHERE username ILIKE $1 LIMIT $2`
	rows, err := r.db.Query(ctx, q, "%"+query+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("user search: %w", err)
	}
	defer rows.Close()

	var users []domain.UserPublic
	for rows.Next() {
		var u domain.UserPublic
		if err := rows.Scan(&u.ID, &u.Username, &u.Bio, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
