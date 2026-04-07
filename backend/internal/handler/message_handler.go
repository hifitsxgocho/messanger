package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"messenger/backend/internal/domain"
	"messenger/backend/internal/middleware"
)

type messageService interface {
	List(ctx context.Context, conversationID, callerID string, after *time.Time) ([]domain.Message, error)
	Send(ctx context.Context, conversationID, callerID, body string) (*domain.Message, error)
	MarkRead(ctx context.Context, conversationID, messageID, callerID string) error
}

type MessageHandler struct {
	svc messageService
}

func NewMessageHandler(svc messageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) List(w http.ResponseWriter, r *http.Request) {
	callerID := middleware.GetUserID(r.Context())
	conversationID := chi.URLParam(r, "id")

	var after *time.Time
	if v := r.URL.Query().Get("after"); v != "" {
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			after = &t
		}
	}

	msgs, err := h.svc.List(r.Context(), conversationID, callerID, after)
	if err != nil {
		if err.Error() == "forbidden" {
			Error(w, http.StatusForbidden, "forbidden")
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if msgs == nil {
		msgs = []domain.Message{}
	}
	JSON(w, http.StatusOK, msgs)
}

func (h *MessageHandler) Send(w http.ResponseWriter, r *http.Request) {
	callerID := middleware.GetUserID(r.Context())
	conversationID := chi.URLParam(r, "id")

	var body struct {
		Body string `json:"body"`
	}
	if err := decodeJSON(r, &body); err != nil || body.Body == "" {
		Error(w, http.StatusBadRequest, "body required")
		return
	}

	msg, err := h.svc.Send(r.Context(), conversationID, callerID, body.Body)
	if err != nil {
		if err.Error() == "forbidden" {
			Error(w, http.StatusForbidden, "forbidden")
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, msg)
}

func (h *MessageHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	callerID := middleware.GetUserID(r.Context())
	conversationID := chi.URLParam(r, "id")
	messageID := chi.URLParam(r, "msgId")

	if err := h.svc.MarkRead(r.Context(), conversationID, messageID, callerID); err != nil {
		if err.Error() == "forbidden" {
			Error(w, http.StatusForbidden, "forbidden")
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
