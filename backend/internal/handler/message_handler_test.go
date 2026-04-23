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

func TestMsgList_Success_Returns200(t *testing.T) {
	msgs := []domain.Message{
		{ID: "msg-1", Body: "hello", SenderID: "uid-1", ConversationID: "conv-1", CreatedAt: time.Now()},
	}
	svc := &mockMsgService{
		listFn: func(_ context.Context, _, _ string, _ *time.Time) ([]domain.Message, error) { return msgs, nil },
	}
	h := NewMessageHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/conversations/conv-1/messages", nil)
	r = withUserID(r, "uid-1")
	r = withChiParams(r, map[string]string{"id": "conv-1"})
	w := httptest.NewRecorder()

	h.List(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp []map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 1 {
		t.Errorf("expected 1 message, got %d", len(resp))
	}
}

func TestMsgList_NoMessages_ReturnsEmptyArray(t *testing.T) {
	svc := &mockMsgService{
		listFn: func(_ context.Context, _, _ string, _ *time.Time) ([]domain.Message, error) { return nil, nil },
	}
	h := NewMessageHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/conversations/conv-1/messages", nil)
	r = withUserID(r, "uid-1")
	r = withChiParams(r, map[string]string{"id": "conv-1"})
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

func TestMsgList_WithAfterParam_PassesTimestamp(t *testing.T) {
	var gotAfter *time.Time
	svc := &mockMsgService{
		listFn: func(_ context.Context, _, _ string, after *time.Time) ([]domain.Message, error) {
			gotAfter = after
			return nil, nil
		},
	}
	h := NewMessageHandler(svc)

	after := time.Now().UTC().Format(time.RFC3339Nano)
	r := httptest.NewRequest(http.MethodGet, "/api/v1/conversations/conv-1/messages?after="+after, nil)
	r = withUserID(r, "uid-1")
	r = withChiParams(r, map[string]string{"id": "conv-1"})
	w := httptest.NewRecorder()

	h.List(w, r)

	if gotAfter == nil {
		t.Error("expected after timestamp to be passed to service")
	}
}

func TestMsgList_Forbidden_Returns403(t *testing.T) {
	svc := &mockMsgService{
		listFn: func(_ context.Context, _, _ string, _ *time.Time) ([]domain.Message, error) {
			return nil, errors.New("forbidden")
		},
	}
	h := NewMessageHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/conversations/conv-1/messages", nil)
	r = withUserID(r, "intruder")
	r = withChiParams(r, map[string]string{"id": "conv-1"})
	w := httptest.NewRecorder()

	h.List(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestMsgSend_ValidBody_Returns201(t *testing.T) {
	msg := &domain.Message{ID: "msg-1", Body: "hello", SenderID: "uid-1", ConversationID: "conv-1"}
	svc := &mockMsgService{
		sendFn: func(_ context.Context, _, _, body string) (*domain.Message, error) { return msg, nil },
	}
	h := NewMessageHandler(svc)

	body := bytes.NewBufferString(`{"body":"hello"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/conversations/conv-1/messages", body)
	r = withUserID(r, "uid-1")
	r = withChiParams(r, map[string]string{"id": "conv-1"})
	w := httptest.NewRecorder()

	h.Send(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMsgSend_EmptyBody_Returns400(t *testing.T) {
	h := NewMessageHandler(&mockMsgService{})

	body := bytes.NewBufferString(`{"body":""}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/conversations/conv-1/messages", body)
	r = withUserID(r, "uid-1")
	r = withChiParams(r, map[string]string{"id": "conv-1"})
	w := httptest.NewRecorder()

	h.Send(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestMsgSend_Forbidden_Returns403(t *testing.T) {
	svc := &mockMsgService{
		sendFn: func(_ context.Context, _, _, _ string) (*domain.Message, error) {
			return nil, errors.New("forbidden")
		},
	}
	h := NewMessageHandler(svc)

	body := bytes.NewBufferString(`{"body":"hello"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/conversations/conv-1/messages", body)
	r = withUserID(r, "intruder")
	r = withChiParams(r, map[string]string{"id": "conv-1"})
	w := httptest.NewRecorder()

	h.Send(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestMarkRead_Success_Returns204(t *testing.T) {
	svc := &mockMsgService{
		markReadFn: func(_ context.Context, _, _, _ string) error { return nil },
	}
	h := NewMessageHandler(svc)

	r := httptest.NewRequest(http.MethodPut, "/api/v1/conversations/conv-1/messages/msg-1/read", nil)
	r = withUserID(r, "uid-1")
	r = withChiParams(r, map[string]string{"id": "conv-1", "msgId": "msg-1"})
	w := httptest.NewRecorder()

	h.MarkRead(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestMarkRead_Forbidden_Returns403(t *testing.T) {
	svc := &mockMsgService{
		markReadFn: func(_ context.Context, _, _, _ string) error { return errors.New("forbidden") },
	}
	h := NewMessageHandler(svc)

	r := httptest.NewRequest(http.MethodPut, "/api/v1/conversations/conv-1/messages/msg-1/read", nil)
	r = withUserID(r, "intruder")
	r = withChiParams(r, map[string]string{"id": "conv-1", "msgId": "msg-1"})
	w := httptest.NewRecorder()

	h.MarkRead(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}
