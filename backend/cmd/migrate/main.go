package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	direction := flag.String("direction", "up", "up | down | reset | status | version")
	steps := flag.Int("steps", 0, "steps for down (0 = all)")
	flag.Parse()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://messenger:messenger_secret@localhost:5432/messenger?sslmode=disable"
	}
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: connect db: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: set dialect: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Database : %s\n", dbURL)
	fmt.Printf("Direction: %s\n\n", *direction)

	var runErr error
	switch *direction {
	case "up":
		runErr = goose.Up(db, migrationsPath)
	case "down":
		if *steps > 0 {
			runErr = goose.DownTo(db, migrationsPath, 0)
			for i := 0; i < *steps && runErr == nil; i++ {
				runErr = goose.Down(db, migrationsPath)
			}
		} else {
			runErr = goose.Down(db, migrationsPath)
		}
	case "reset":
		runErr = goose.Reset(db, migrationsPath)
	case "status":
		runErr = goose.Status(db, migrationsPath)
	case "version":
		runErr = goose.Version(db, migrationsPath)
	default:
		fmt.Fprintf(os.Stderr, "Unknown direction: %s\n", *direction)
		os.Exit(1)
	}

	if runErr != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", runErr)
		os.Exit(1)
	}
}
