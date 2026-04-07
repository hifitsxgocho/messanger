package repository

import (
	"context"
	"messenger/backend/internal/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Search(ctx context.Context, query string, limit int) ([]domain.UserPublic, error)
}

type ConversationRepository interface {
	Create(ctx context.Context, userAID, userBID string) (*domain.ConversationRow, error)
	GetByID(ctx context.Context, id string) (*domain.ConversationRow, error)
	FindBetween(ctx context.Context, userAID, userBID string) (*domain.ConversationRow, error)
	ListForUser(ctx context.Context, userID string) ([]domain.ConversationRow, error)
}

type MessageRepository interface {
	Create(ctx context.Context, msg *domain.Message) error
	ListByConversation(ctx context.Context, conversationID string, after *time.Time, limit int) ([]domain.Message, error)
	GetLastMessage(ctx context.Context, conversationID string) (*domain.Message, error)
	CountUnread(ctx context.Context, conversationID, readerID string) (int, error)
	MarkRead(ctx context.Context, messageID, readerID string) error
}
