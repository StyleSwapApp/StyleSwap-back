package usermanagement

import (
	"StyleSwap/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux {
	userConfig := New(configuration)
	router := chi.NewRouter()

	router.Get("/user/{id4Update}", userConfig.UpdateHandler)
	router.Delete("/{id4Delete}", userConfig.DeleteHandler)
	router.Get("/{id}", userConfig.GetUserHandler)

	return router
}
