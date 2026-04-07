package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"messenger/backend/internal/domain"
	"messenger/backend/internal/middleware"
)

type conversationService interface {
	GetOrCreate(ctx context.Context, callerID, partnerID string) (*domain.Conversation, error)
	GetByID(ctx context.Context, id, callerID string) (*domain.Conversation, error)
	ListForUser(ctx context.Context, userID string) ([]domain.Conversation, error)
}

type ConversationHandler struct {
	svc conversationService
}

func NewConversationHandler(svc conversationService) *ConversationHandler {
	return &ConversationHandler{svc: svc}
}

func (h *ConversationHandler) List(w http.ResponseWriter, r *http.Request) {
	callerID := middleware.GetUserID(r.Context())

	convs, err := h.svc.ListForUser(r.Context(), callerID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if convs == nil {
		convs = []domain.Conversation{}
	}
	JSON(w, http.StatusOK, convs)
}

func (h *ConversationHandler) Create(w http.ResponseWriter, r *http.Request) {
	callerID := middleware.GetUserID(r.Context())

	var body struct {
		UserID string `json:"userId"`
	}
	if err := decodeJSON(r, &body); err != nil || body.UserID == "" {
		Error(w, http.StatusBadRequest, "userId required")
		return
	}

	conv, err := h.svc.GetOrCreate(r.Context(), callerID, body.UserID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, conv)
}

func (h *ConversationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	callerID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	conv, err := h.svc.GetByID(r.Context(), id, callerID)
	if err != nil {
		if err.Error() == "forbidden" {
			Error(w, http.StatusForbidden, "forbidden")
			return
		}
		Error(w, http.StatusNotFound, "not found")
		return
	}
	JSON(w, http.StatusOK, conv)
}
