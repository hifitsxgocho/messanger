package handler

import (
	"encoding/json"
	"messenger/backend/internal/middleware"
	"messenger/backend/internal/service"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil || user == nil {
		Error(w, http.StatusNotFound, "user not found")
		return
	}
	JSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req struct {
		Username string `json:"username"`
		Bio      string `json:"bio"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	user, err := h.users.UpdateMe(r.Context(), userID, service.UpdateUserInput{
		Username: req.Username,
		Bio:      req.Bio,
	})
	if err != nil {
		Error(w, http.StatusInternalServerError, "update failed")
		return
	}
	JSON(w, http.StatusOK, user)
}

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	r.ParseMultipartForm(5 << 20) // 5MB
	file, header, err := r.FormFile("avatar")
	if err != nil {
		Error(w, http.StatusBadRequest, "avatar file required")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		Error(w, http.StatusBadRequest, "only jpg, png, webp allowed")
		return
	}

	avatarURL, err := h.users.UploadAvatar(r.Context(), userID, file, ext)
	if err != nil {
		Error(w, http.StatusInternalServerError, "upload failed")
		return
	}
	JSON(w, http.StatusOK, map[string]string{"avatarUrl": avatarURL})
}

func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	results, err := h.users.Search(r.Context(), q)
	if err != nil {
		Error(w, http.StatusInternalServerError, "search failed")
		return
	}
	if results == nil {
		JSON(w, http.StatusOK, []any{})
		return
	}
	JSON(w, http.StatusOK, results)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.users.GetPublicByID(r.Context(), id)
	if err != nil || user == nil {
		Error(w, http.StatusNotFound, "user not found")
		return
	}
	JSON(w, http.StatusOK, user)
}
