package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"messenger/backend/internal/domain"
	"messenger/backend/internal/middleware"
	"messenger/backend/internal/repository"
)

// mockUserRepoH implements repository.UserRepository for handler-level tests.
type mockUserRepoH struct {
	createFn     func(ctx context.Context, user *domain.User) error
	getByIDFn    func(ctx context.Context, id string) (*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	updateFn     func(ctx context.Context, user *domain.User) error
	searchFn     func(ctx context.Context, query string, limit int) ([]domain.UserPublic, error)
}

var _ repository.UserRepository = (*mockUserRepoH)(nil)

func (m *mockUserRepoH) Create(ctx context.Context, user *domain.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepoH) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepoH) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *mockUserRepoH) Update(ctx context.Context, user *domain.User) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepoH) Search(ctx context.Context, query string, limit int) ([]domain.UserPublic, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, query, limit)
	}
	return nil, nil
}

// mockConvService implements the unexported conversationService interface.
type mockConvService struct {
	getOrCreateFn func(ctx context.Context, callerID, partnerID string) (*domain.Conversation, error)
	getByIDFn     func(ctx context.Context, id, callerID string) (*domain.Conversation, error)
	listForUserFn func(ctx context.Context, userID string) ([]domain.Conversation, error)
}

func (m *mockConvService) GetOrCreate(ctx context.Context, callerID, partnerID string) (*domain.Conversation, error) {
	if m.getOrCreateFn != nil {
		return m.getOrCreateFn(ctx, callerID, partnerID)
	}
	return nil, nil
}

func (m *mockConvService) GetByID(ctx context.Context, id, callerID string) (*domain.Conversation, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id, callerID)
	}
	return nil, nil
}

func (m *mockConvService) ListForUser(ctx context.Context, userID string) ([]domain.Conversation, error) {
	if m.listForUserFn != nil {
		return m.listForUserFn(ctx, userID)
	}
	return nil, nil
}

// mockMsgService implements the unexported messageService interface.
type mockMsgService struct {
	listFn     func(ctx context.Context, conversationID, callerID string, after *time.Time) ([]domain.Message, error)
	sendFn     func(ctx context.Context, conversationID, callerID, body string) (*domain.Message, error)
	markReadFn func(ctx context.Context, conversationID, messageID, callerID string) error
}

func (m *mockMsgService) List(ctx context.Context, conversationID, callerID string, after *time.Time) ([]domain.Message, error) {
	if m.listFn != nil {
		return m.listFn(ctx, conversationID, callerID, after)
	}
	return nil, nil
}

func (m *mockMsgService) Send(ctx context.Context, conversationID, callerID, body string) (*domain.Message, error) {
	if m.sendFn != nil {
		return m.sendFn(ctx, conversationID, callerID, body)
	}
	return nil, nil
}

func (m *mockMsgService) MarkRead(ctx context.Context, conversationID, messageID, callerID string) error {
	if m.markReadFn != nil {
		return m.markReadFn(ctx, conversationID, messageID, callerID)
	}
	return nil
}

// withUserID injects a userID into the request context (simulates auth middleware).
func withUserID(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	return r.WithContext(ctx)
}

// withChiParams adds chi URL parameters to the request context.
func withChiParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
