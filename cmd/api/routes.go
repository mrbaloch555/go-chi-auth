package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Config) routes() http.Handler {

	mux := chi.NewRouter()

	mux.Route("/users", func(r chi.Router) {
		r.Post("/login", app.Login)
		r.Post("/register", app.Register)
	})

	mux.With(app.Middleware("user")).Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hey, this is auth route"))
	})
	return mux
}
