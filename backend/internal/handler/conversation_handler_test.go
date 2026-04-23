package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"messenger/backend/internal/domain"
)

func TestConvList_NoConversations_ReturnsEmptyArray(t *testing.T) {
	svc := &mockConvService{
		listForUserFn: func(_ context.Context, _ string) ([]domain.Conversation, error) { return nil, nil },
	}
	h := NewConversationHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/conversations", nil)
	r = withUserID(r, "uid-1")
	w := httptest.NewRecorder()

	h.List(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp []any
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 0 {
		t.Errorf("expected empty array, got %v", resp)
	}
}

func TestConvList_WithConversations_ReturnsList(t *testing.T) {
	convs := []domain.Conversation{
		{ID: "conv-1", Partner: domain.UserPublic{ID: "uid-2", Username: "bob"}, CreatedAt: time.Now()},
	}
	svc := &mockConvService{
		listForUserFn: func(_ context.Context, _ string) ([]domain.Conversation, error) { return convs, nil },
	}
	h := NewConversationHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/conversations", nil)
	r = withUserID(r, "uid-1")
	w := httptest.NewRecorder()

	h.List(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp []map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 1 {
		t.Errorf("expected 1 conversation, got %d", len(resp))
	}
}

func TestConvCreate_ValidUserID_Returns200WithConversation(t *testing.T) {
	conv := &domain.Conversation{ID: "conv-1", Partner: domain.UserPublic{ID: "uid-2"}}
	svc := &mockConvService{
		getOrCreateFn: func(_ context.Context, _, _ string) (*domain.Conversation, error) { return conv, nil },
	}
	h := NewConversationHandler(svc)

	body := bytes.NewBufferString(`{"userId":"uid-2"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/conversations", body)
	r = withUserID(r, "uid-1")
	w := httptest.NewRecorder()

	h.Create(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["id"] != "conv-1" {
		t.Errorf("expected conv id conv-1, got %v", resp["id"])
	}
}

func TestConvCreate_MissingUserID_Returns400(t *testing.T) {
	h := NewConversationHandler(&mockConvService{})

	body := bytes.NewBufferString(`{}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/conversations", body)
	r = withUserID(r, "uid-1")
	w := httptest.NewRecorder()

	h.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestConvCreate_InvalidJSON_Returns400(t *testing.T) {
	h := NewConversationHandler(&mockConvService{})

	r := httptest.NewRequest(http.MethodPost, "/api/v1/conversations", bytes.NewBufferString("bad"))
	r = withUserID(r, "uid-1")
	w := httptest.NewRecorder()

	h.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestConvGetByID_Found_Returns200(t *testing.T) {
	conv := &domain.Conversation{ID: "conv-1", Partner: domain.UserPublic{ID: "uid-2"}}
	svc := &mockConvService{
		getByIDFn: func(_ context.Context, _, _ string) (*domain.Conversation, error) { return conv, nil },
	}
	h := NewConversationHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/conversations/conv-1", nil)
	r = withUserID(r, "uid-1")
	r = withChiParams(r, map[string]string{"id": "conv-1"})
	w := httptest.NewRecorder()

	h.GetByID(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestConvGetByID_Forbidden_Returns403(t *testing.T) {
	svc := &mockConvService{
		getByIDFn: func(_ context.Context, _, _ string) (*domain.Conversation, error) {
			return nil, errors.New("forbidden")
		},
	}
	h := NewConversationHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/conversations/conv-1", nil)
	r = withUserID(r, "intruder")
	r = withChiParams(r, map[string]string{"id": "conv-1"})
	w := httptest.NewRecorder()

	h.GetByID(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestConvGetByID_NotFound_Returns404(t *testing.T) {
	svc := &mockConvService{
		getByIDFn: func(_ context.Context, _, _ string) (*domain.Conversation, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewConversationHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/conversations/conv-999", nil)
	r = withUserID(r, "uid-1")
	r = withChiParams(r, map[string]string{"id": "conv-999"})
	w := httptest.NewRecorder()

	h.GetByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
