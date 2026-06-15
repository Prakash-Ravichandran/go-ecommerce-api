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

	// run conventional REST server
	// h := api.mount()
	// err2 := api.run(h)

	// if err2 != nil {
	// 	slog.Error("Server has failed to start", "error", err)
	// 	os.Exit(1)
	// }

	// product service go routine
	go func() {
		grpcErr := api.runProductsGRPC()

		if grpcErr != nil {
			slog.Error("gRPC products server failed to start", "error", grpcErr)
			os.Exit(1)
		}
	}()

	// order service go routine
	go func() {
		grpcErr := api.runOrderGRPC()

		if grpcErr != nil {
			slog.Error("gRPC order server failed to start", "error", grpcErr)
		}
	}()

	// health service go routine
	go func() {
		grpcErr := api.runHealthGRPC()

		if grpcErr != nil {
			slog.Error("gRPC health server failed to start", "error", grpcErr)
		}
	}()

	// 3. Start the HTTP Gateway Router on the main thread (blocking)
	httpHandler := api.mount()
	if err := api.run(httpHandler); err != nil {
		slog.Error("HTTP Gateway failed: %v", "error", err)
	}
}
