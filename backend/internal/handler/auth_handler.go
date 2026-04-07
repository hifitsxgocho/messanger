package handler

import (
	"encoding/json"
	"errors"
	"messenger/backend/internal/service"
	"net/http"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" || req.Username == "" || len(req.Password) < 8 {
		Error(w, http.StatusBadRequest, "email, username and password (min 8 chars) required")
		return
	}

	result, err := h.auth.Register(r.Context(), service.RegisterInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrEmailTaken) {
		Error(w, http.StatusConflict, "email already taken")
		return
	}
	if err != nil {
		Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	JSON(w, http.StatusCreated, map[string]any{"token": result.Token, "user": result.User})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidCreds) {
		Error(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	JSON(w, http.StatusOK, map[string]any{"token": result.Token, "user": result.User})
}
