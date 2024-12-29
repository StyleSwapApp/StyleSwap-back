package usermanagement

import (
	"StyleSwap/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux{
	userConfig := New(configuration)
	router := chi.NewRouter()

	router.Post("/newUser", userConfig.UserHandler)
	router.Get("/login", userConfig.LoginHandler)

	return router
}