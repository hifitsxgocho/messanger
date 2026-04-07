package service

import (
	"context"
	"errors"

	"messenger/backend/internal/domain"
	"messenger/backend/internal/repository"
)

type ConversationService struct {
	convRepo repository.ConversationRepository
	msgRepo  repository.MessageRepository
	userRepo repository.UserRepository
}

func NewConversationService(
	convRepo repository.ConversationRepository,
	msgRepo repository.MessageRepository,
	userRepo repository.UserRepository,
) *ConversationService {
	return &ConversationService{convRepo: convRepo, msgRepo: msgRepo, userRepo: userRepo}
}

func (s *ConversationService) GetOrCreate(ctx context.Context, callerID, partnerID string) (*domain.Conversation, error) {
	row, err := s.convRepo.Create(ctx, callerID, partnerID)
	if err != nil {
		return nil, err
	}
	return s.buildConversation(ctx, row, callerID)
}

func (s *ConversationService) GetByID(ctx context.Context, id, callerID string) (*domain.Conversation, error) {
	row, err := s.convRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row.UserAID != callerID && row.UserBID != callerID {
		return nil, errors.New("forbidden")
	}
	return s.buildConversation(ctx, row, callerID)
}

func (s *ConversationService) ListForUser(ctx context.Context, userID string) ([]domain.Conversation, error) {
	rows, err := s.convRepo.ListForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	convs := make([]domain.Conversation, 0, len(rows))
	for _, row := range rows {
		c, err := s.buildConversation(ctx, &row, userID)
		if err != nil {
			continue
		}
		convs = append(convs, *c)
	}
	return convs, nil
}

func (s *ConversationService) buildConversation(ctx context.Context, row *domain.ConversationRow, callerID string) (*domain.Conversation, error) {
	partnerID := row.UserAID
	if partnerID == callerID {
		partnerID = row.UserBID
	}

	partner, err := s.userRepo.GetByID(ctx, partnerID)
	if err != nil {
		return nil, err
	}

	conv := &domain.Conversation{
		ID:        row.ID,
		Partner:   partner.ToPublic(),
		CreatedAt: row.CreatedAt,
	}

	lastMsg, err := s.msgRepo.GetLastMessage(ctx, row.ID)
	if err == nil && lastMsg != nil {
		conv.LastMessage = &domain.LastMessage{
			Body:      lastMsg.Body,
			SenderID:  lastMsg.SenderID,
			CreatedAt: lastMsg.CreatedAt,
		}
	}

	unread, err := s.msgRepo.CountUnread(ctx, row.ID, callerID)
	if err == nil {
		conv.UnreadCount = unread
	}

	return conv, nil
}
