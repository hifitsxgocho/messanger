package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"messenger/backend/internal/domain"
)

type MessageRepo struct {
	db *pgxpool.Pool
}

func NewMessageRepo(db *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) Create(ctx context.Context, msg *domain.Message) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO messages (id, conversation_id, sender_id, body, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW())
		RETURNING id, created_at`,
		msg.ConversationID, msg.SenderID, msg.Body)
	if err != nil {
		return err
	}

	row := r.db.QueryRow(ctx, `
		SELECT id, conversation_id, sender_id, body, created_at, read_at
		FROM messages
		WHERE conversation_id = $1 AND sender_id = $2
		ORDER BY created_at DESC LIMIT 1`,
		msg.ConversationID, msg.SenderID)
	return row.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Body, &msg.CreatedAt, &msg.ReadAt)
}

func (r *MessageRepo) ListByConversation(ctx context.Context, conversationID string, after *time.Time, limit int) ([]domain.Message, error) {
	var rows interface{ Next() bool; Scan(...any) error; Close(); Err() error }
	var err error

	if after != nil {
		rows, err = r.db.Query(ctx, `
			SELECT id, conversation_id, sender_id, body, created_at, read_at
			FROM messages
			WHERE conversation_id = $1 AND created_at > $2
			ORDER BY created_at ASC
			LIMIT $3`, conversationID, after, limit)
	} else {
		rows, err = r.db.Query(ctx, `
			SELECT id, conversation_id, sender_id, body, created_at, read_at
			FROM messages
			WHERE conversation_id = $1
			ORDER BY created_at ASC
			LIMIT $2`, conversationID, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Body, &m.CreatedAt, &m.ReadAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

func (r *MessageRepo) GetLastMessage(ctx context.Context, conversationID string) (*domain.Message, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, conversation_id, sender_id, body, created_at, read_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT 1`, conversationID)

	var m domain.Message
	if err := row.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Body, &m.CreatedAt, &m.ReadAt); err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MessageRepo) CountUnread(ctx context.Context, conversationID, readerID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM messages
		WHERE conversation_id = $1
		  AND sender_id != $2
		  AND read_at IS NULL`, conversationID, readerID).Scan(&count)
	return count, err
}

func (r *MessageRepo) MarkRead(ctx context.Context, messageID, readerID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE messages SET read_at = NOW()
		WHERE id = $1
		  AND sender_id != $2
		  AND read_at IS NULL`, messageID, readerID)
	return err
}
