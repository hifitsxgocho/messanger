package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"messenger/backend/internal/domain"
	"messenger/backend/internal/service"
)

func newUserHandler(t *testing.T, repo *mockUserRepoH) *UserHandler {
	t.Helper()
	return NewUserHandler(service.NewUserService(repo, t.TempDir()))
}

func TestGetMe_UserFound_Returns200(t *testing.T) {
	repo := &mockUserRepoH{
		getByIDFn: func(_ context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: id, Username: "alice", Email: "a@b.com"}, nil
		},
	}
	h := newUserHandler(t, repo)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	r = withUserID(r, "uid-1")
	w := httptest.NewRecorder()

	h.GetMe(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["username"] != "alice" {
		t.Errorf("expected username alice, got %v", resp["username"])
	}
}

func TestGetMe_UserNotFound_Returns404(t *testing.T) {
	repo := &mockUserRepoH{
		getByIDFn: func(_ context.Context, _ string) (*domain.User, error) { return nil, nil },
	}
	h := newUserHandler(t, repo)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	r = withUserID(r, "uid-ghost")
	w := httptest.NewRecorder()

	h.GetMe(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateMe_ValidInput_Returns200WithUpdatedUser(t *testing.T) {
	user := &domain.User{ID: "uid-1", Username: "old", Bio: "old bio"}
	repo := &mockUserRepoH{
		getByIDFn: func(_ context.Context, _ string) (*domain.User, error) { return user, nil },
		updateFn:  func(_ context.Context, _ *domain.User) error { return nil },
	}
	h := newUserHandler(t, repo)

	body := bytes.NewBufferString(`{"username":"new","bio":"new bio"}`)
	r := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", body)
	r = withUserID(r, "uid-1")
	w := httptest.NewRecorder()

	h.UpdateMe(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["username"] != "new" {
		t.Errorf("expected username new, got %v", resp["username"])
	}
}

func TestSearch_EmptyQuery_ReturnsEmptyArray(t *testing.T) {
	h := newUserHandler(t, &mockUserRepoH{})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/search?q=", nil)
	w := httptest.NewRecorder()

	h.Search(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp []any
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 0 {
		t.Errorf("expected empty array, got %v", resp)
	}
}

func TestSearch_WithQuery_ReturnsResults(t *testing.T) {
	repo := &mockUserRepoH{
		searchFn: func(_ context.Context, _ string, _ int) ([]domain.UserPublic, error) {
			return []domain.UserPublic{{ID: "uid-1", Username: "alice"}}, nil
		},
	}
	h := newUserHandler(t, repo)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/search?q=ali", nil)
	w := httptest.NewRecorder()

	h.Search(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp []map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp) != 1 || resp[0]["username"] != "alice" {
		t.Errorf("unexpected results: %v", resp)
	}
}

func TestGetByID_UserFound_Returns200(t *testing.T) {
	repo := &mockUserRepoH{
		getByIDFn: func(_ context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: id, Username: "bob"}, nil
		},
	}
	h := newUserHandler(t, repo)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/uid-2", nil)
	r = withChiParams(r, map[string]string{"id": "uid-2"})
	w := httptest.NewRecorder()

	h.GetByID(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["username"] != "bob" {
		t.Errorf("expected username bob, got %v", resp["username"])
	}
}

func TestGetByID_UserNotFound_Returns404(t *testing.T) {
	repo := &mockUserRepoH{
		getByIDFn: func(_ context.Context, _ string) (*domain.User, error) { return nil, nil },
	}
	h := newUserHandler(t, repo)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/uid-ghost", nil)
	r = withChiParams(r, map[string]string{"id": "uid-ghost"})
	w := httptest.NewRecorder()

	h.GetByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUploadAvatar_ValidPNG_Returns200WithURL(t *testing.T) {
	user := &domain.User{ID: "uid-1"}
	repo := &mockUserRepoH{
		getByIDFn: func(_ context.Context, _ string) (*domain.User, error) { return user, nil },
		updateFn:  func(_ context.Context, _ *domain.User) error { return nil },
	}
	h := newUserHandler(t, repo)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	part, _ := mw.CreateFormFile("avatar", "photo.png")
	part.Write([]byte("fake-png-bytes")) //nolint
	mw.Close()

	r := httptest.NewRequest(http.MethodPost, "/api/v1/users/me/avatar", &buf)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	r = withUserID(r, "uid-1")
	w := httptest.NewRecorder()

	h.UploadAvatar(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["avatarUrl"] == "" {
		t.Error("expected avatarUrl in response")
	}
}

func TestUploadAvatar_InvalidExtension_Returns400(t *testing.T) {
	h := newUserHandler(t, &mockUserRepoH{})

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	part, _ := mw.CreateFormFile("avatar", "file.gif")
	part.Write([]byte("fake-data")) //nolint
	mw.Close()

	r := httptest.NewRequest(http.MethodPost, "/api/v1/users/me/avatar", &buf)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	r = withUserID(r, "uid-1")
	w := httptest.NewRecorder()

	h.UploadAvatar(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
