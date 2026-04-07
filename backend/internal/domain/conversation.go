package domain

import "time"

type LastMessage struct {
	Body      string    `json:"body"`
	SenderID  string    `json:"senderId"`
	CreatedAt time.Time `json:"createdAt"`
}

type Conversation struct {
	ID          string      `json:"id"`
	Partner     UserPublic  `json:"partner"`
	LastMessage *LastMessage `json:"lastMessage"`
	UnreadCount int         `json:"unreadCount"`
	CreatedAt   time.Time   `json:"createdAt"`
}

// ConversationRow is the raw DB row for a conversation.
type ConversationRow struct {
	ID        string
	UserAID   string
	UserBID   string
	CreatedAt time.Time
}
