package main

import (
	"StyleSwap/config"
	"StyleSwap/pkg/article"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Mount("/api", article.Routes(configuration))
	return router
}