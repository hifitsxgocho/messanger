package postgres

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"messenger/backend/internal/domain"
)

type ConversationRepo struct {
	db *pgxpool.Pool
}

func NewConversationRepo(db *pgxpool.Pool) *ConversationRepo {
	return &ConversationRepo{db: db}
}

// normalizePair ensures the smaller UUID is always user_a_id to prevent duplicates.
func normalizePair(a, b string) (string, string) {
	if strings.Compare(a, b) > 0 {
		return b, a
	}
	return a, b
}

func (r *ConversationRepo) Create(ctx context.Context, userAID, userBID string) (*domain.ConversationRow, error) {
	a, b := normalizePair(userAID, userBID)

	row := r.db.QueryRow(ctx, `
		INSERT INTO conversations (user_a_id, user_b_id)
		VALUES ($1, $2)
		ON CONFLICT (user_a_id, user_b_id) DO UPDATE
			SET user_a_id = EXCLUDED.user_a_id
		RETURNING id, user_a_id, user_b_id, created_at`, a, b)

	var c domain.ConversationRow
	if err := row.Scan(&c.ID, &c.UserAID, &c.UserBID, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConversationRepo) GetByID(ctx context.Context, id string) (*domain.ConversationRow, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_a_id, user_b_id, created_at
		FROM conversations WHERE id = $1`, id)

	var c domain.ConversationRow
	if err := row.Scan(&c.ID, &c.UserAID, &c.UserBID, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConversationRepo) FindBetween(ctx context.Context, userAID, userBID string) (*domain.ConversationRow, error) {
	a, b := normalizePair(userAID, userBID)

	row := r.db.QueryRow(ctx, `
		SELECT id, user_a_id, user_b_id, created_at
		FROM conversations WHERE user_a_id = $1 AND user_b_id = $2`, a, b)

	var c domain.ConversationRow
	if err := row.Scan(&c.ID, &c.UserAID, &c.UserBID, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConversationRepo) ListForUser(ctx context.Context, userID string) ([]domain.ConversationRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_a_id, user_b_id, created_at
		FROM conversations
		WHERE user_a_id = $1 OR user_b_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []domain.ConversationRow
	for rows.Next() {
		var c domain.ConversationRow
		if err := rows.Scan(&c.ID, &c.UserAID, &c.UserBID, &c.CreatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}
	return convs, rows.Err()
}
