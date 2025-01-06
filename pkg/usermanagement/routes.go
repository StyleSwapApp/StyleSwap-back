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
	router.Put("/updateUser", userConfig.UpdateHandler)
	router.Delete("/deleteUser", userConfig.DeleteHandler)

	return router
}