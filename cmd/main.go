package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Prakash-Ravichandran/go-ecommerce-api/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback for local dev
	}

	cfg := config{
		addr: ":" + port,
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=ecom sslmode=disable"),
		},
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("connected to database", "dsn", cfg.db.dsn)

	api := application{
		config: cfg,
		db:     conn,
	}

	h := api.mount()
	err2 := api.run(h)

	if err2 != nil {
		slog.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}
}
