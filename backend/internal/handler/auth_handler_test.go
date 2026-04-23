package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"messenger/backend/internal/domain"
	"messenger/backend/internal/service"
)

func newAuthHandler(repo *mockUserRepoH) *AuthHandler {
	return NewAuthHandler(service.NewAuthService(repo, "test-secret"))
}

func TestAuthRegister_ValidInput_Returns201WithToken(t *testing.T) {
	repo := &mockUserRepoH{
		getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) { return nil, nil },
		createFn:     func(_ context.Context, _ *domain.User) error { return nil },
	}
	h := newAuthHandler(repo)

	body := bytes.NewBufferString(`{"email":"a@b.com","username":"alice","password":"password123"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if _, ok := resp["token"]; !ok {
		t.Error("expected token in response body")
	}
	if _, ok := resp["user"]; !ok {
		t.Error("expected user in response body")
	}
}

func TestAuthRegister_ShortPassword_Returns400(t *testing.T) {
	h := newAuthHandler(&mockUserRepoH{})

	body := bytes.NewBufferString(`{"email":"a@b.com","username":"alice","password":"short"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthRegister_MissingEmail_Returns400(t *testing.T) {
	h := newAuthHandler(&mockUserRepoH{})

	body := bytes.NewBufferString(`{"username":"alice","password":"password123"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthRegister_InvalidJSON_Returns400(t *testing.T) {
	h := newAuthHandler(&mockUserRepoH{})

	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthRegister_EmailTaken_Returns409(t *testing.T) {
	repo := &mockUserRepoH{
		getByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
			return &domain.User{Email: email}, nil
		},
	}
	h := newAuthHandler(repo)

	body := bytes.NewBufferString(`{"email":"taken@b.com","username":"bob","password":"password123"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	w := httptest.NewRecorder()

	h.Register(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestAuthLogin_ValidCredentials_Returns200WithToken(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &mockUserRepoH{
		getByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
			return &domain.User{ID: "uid-1", Email: email, PasswordHash: string(hash)}, nil
		},
	}
	h := newAuthHandler(repo)

	body := bytes.NewBufferString(`{"email":"a@b.com","password":"password123"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if _, ok := resp["token"]; !ok {
		t.Error("expected token in response body")
	}
}

func TestAuthLogin_WrongPassword_Returns401(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	repo := &mockUserRepoH{
		getByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
			return &domain.User{ID: "uid-1", Email: email, PasswordHash: string(hash)}, nil
		},
	}
	h := newAuthHandler(repo)

	body := bytes.NewBufferString(`{"email":"a@b.com","password":"wrong"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthLogin_InvalidJSON_Returns400(t *testing.T) {
	h := newAuthHandler(&mockUserRepoH{})

	r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
