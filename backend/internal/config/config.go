package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	JWTSecret   string
	AvatarDir   string
}

func Load() *Config {
	// Load .env if it exists (ignored in production/Docker where env vars are set directly)
	godotenv.Load()

	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://messenger:messenger_secret@localhost:5432/messenger?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		AvatarDir:   getEnv("AVATAR_DIR", "./avatars"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
