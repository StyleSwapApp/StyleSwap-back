package login

import (
	"StyleSwap/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux {
	userConfig := New(configuration)
	router := chi.NewRouter()

	router.Get("/login", userConfig.LoginHandler)

	return router
}
