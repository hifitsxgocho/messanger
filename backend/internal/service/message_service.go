package service

import (
	"context"
	"errors"
	"time"

	"messenger/backend/internal/domain"
	"messenger/backend/internal/repository"
)

type MessageService struct {
	msgRepo  repository.MessageRepository
	convRepo repository.ConversationRepository
}

func NewMessageService(msgRepo repository.MessageRepository, convRepo repository.ConversationRepository) *MessageService {
	return &MessageService{msgRepo: msgRepo, convRepo: convRepo}
}

func (s *MessageService) List(ctx context.Context, conversationID, callerID string, after *time.Time) ([]domain.Message, error) {
	if err := s.checkAccess(ctx, conversationID, callerID); err != nil {
		return nil, err
	}
	return s.msgRepo.ListByConversation(ctx, conversationID, after, 100)
}

func (s *MessageService) Send(ctx context.Context, conversationID, callerID, body string) (*domain.Message, error) {
	if err := s.checkAccess(ctx, conversationID, callerID); err != nil {
		return nil, err
	}
	msg := &domain.Message{
		ConversationID: conversationID,
		SenderID:       callerID,
		Body:           body,
	}
	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *MessageService) MarkRead(ctx context.Context, conversationID, messageID, callerID string) error {
	if err := s.checkAccess(ctx, conversationID, callerID); err != nil {
		return err
	}
	return s.msgRepo.MarkRead(ctx, messageID, callerID)
}

func (s *MessageService) checkAccess(ctx context.Context, conversationID, callerID string) error {
	row, err := s.convRepo.GetByID(ctx, conversationID)
	if err != nil {
		return err
	}
	if row.UserAID != callerID && row.UserBID != callerID {
		return errors.New("forbidden")
	}
	return nil
}
