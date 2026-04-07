package main

import (
	"context"
	"messenger/backend/internal/config"
	"messenger/backend/internal/handler"
	"messenger/backend/internal/middleware"
	"messenger/backend/internal/repository/postgres"
	"messenger/backend/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	cfg := config.Load()

	ctx := context.Background()

	// Connect to DB
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to postgres")
	}
	defer pool.Close()
	log.Info().Msg("connected to postgres")

	// Run migrations via goose
	db := stdlib.OpenDBFromPool(pool)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal().Err(err).Msg("goose set dialect")
	}
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("migrations failed")
	}
	log.Info().Msg("migrations applied")

	// Repositories
	userRepo := postgres.NewUserRepo(pool)
	convRepo := postgres.NewConversationRepo(pool)
	msgRepo := postgres.NewMessageRepo(pool)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	userSvc := service.NewUserService(userRepo, cfg.AvatarDir)
	convSvc := service.NewConversationService(convRepo, msgRepo, userRepo)
	msgSvc := service.NewMessageService(msgRepo, convRepo)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	convH := handler.NewConversationHandler(convSvc)
	msgH := handler.NewMessageHandler(msgSvc)

	router := handler.NewRouter(authH, userH, convH, msgH, middleware.TokenValidator(authSvc), cfg.AvatarDir)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.ServerPort).Msg("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down")
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutCtx)
}
