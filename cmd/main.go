package main

import (
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
	}

	api := application{
		config: cfg,
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	h := api.mount()
	err := api.run(h)

	if err != nil {
		slog.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}
}
