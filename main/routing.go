package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (config *Config) routing() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(config.SessionSaving)
	mux.Get("/", config.homepage)
	mux.Get("/login", config.login)
	mux.Get("/signup", config.signup)
	mux.Post("/login", config.postLogin)
	mux.Post("/signup", config.postSignup)
	mux.Get("/activate/", config.activate)
	mux.Get("/logout", config.logout)
	mux.Get("/userspace", config.userOnly)
	return mux
}

func (config *Config) SessionSaving(next http.Handler) http.Handler {
	return config.session.LoadAndSave(next)
}
