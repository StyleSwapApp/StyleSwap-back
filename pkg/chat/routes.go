package chat

import (
	"StyleSwap/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux{
	messageConfig := New(configuration)
	router := chi.NewRouter()
	
	router.Get("/ws", messageConfig.HandleWebSocket)

	return router
}