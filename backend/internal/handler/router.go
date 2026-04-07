package handler

import (
	"messenger/backend/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	authH *AuthHandler,
	userH *UserHandler,
	convH *ConversationHandler,
	msgH *MessageHandler,
	authSvc middleware.TokenValidator,
	avatarDir string,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Serve uploaded avatars
	r.Handle("/avatars/*", http.StripPrefix("/avatars/", http.FileServer(http.Dir(avatarDir))))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public
		r.Post("/auth/register", authH.Register)
		r.Post("/auth/login", authH.Login)

		// Protected
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))

			r.Get("/users/me", userH.GetMe)
			r.Put("/users/me", userH.UpdateMe)
			r.Post("/users/me/avatar", userH.UploadAvatar)
			r.Get("/users/search", userH.Search)
			r.Get("/users/{id}", userH.GetByID)

			r.Get("/conversations", convH.List)
			r.Post("/conversations", convH.Create)
			r.Get("/conversations/{id}", convH.GetByID)

			r.Get("/conversations/{id}/messages", msgH.List)
			r.Post("/conversations/{id}/messages", msgH.Send)
			r.Put("/conversations/{id}/messages/{msgId}/read", msgH.MarkRead)
		})
	})

	return r
}
