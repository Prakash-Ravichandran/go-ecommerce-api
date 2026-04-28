package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	// mount
	// run
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // important for rate limiting - random ID for each request
	r.Use(middleware.RealIP)    // get the real IP address of the client
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) //  recover from crashes

	// set a timeout value on the request context(ctx), that will signal
	// through ctx.Done() that the request has timed out
	// and further processing should be stopped
	r.Use(middleware.Timeout(60 * time.Second)) // if request takes more than 60s then stop it

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})
	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}

type config struct {
	addr string
	dbs  dbConfig
}

type dbConfig struct {
	dsn string
}
