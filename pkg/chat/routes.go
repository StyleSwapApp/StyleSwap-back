package chat

import (
	"StyleSwap/config"

	"github.com/go-chi/chi/v5"
)

// Routes initialise les routes pour le chat
func Routes(configuration *config.Config) *chi.Mux {
	messageConfig := New(configuration)
	router := chi.NewRouter()

	// Route WebSocket pour le chat en temps réel
	router.Get("/ws", messageConfig.HandleWebSocket)

	// Route pour envoyer un message via HTTP (si nécessaire)
	router.Post("/message", messageConfig.SendMessageHandler)

	return router
}
