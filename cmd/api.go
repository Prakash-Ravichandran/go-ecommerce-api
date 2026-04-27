package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type application struct {
	config config
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	return r
}

type config struct {
	addr string
	dbs  dbConfig
}

type dbConfig struct {
	dsn string
}